package work_log_summary_service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"worknote-api/config"
	"worknote-api/model"
	"worknote-api/repos/work_log_repo"
	"worknote-api/repos/work_log_summary_repo"
)

// OpenRouter API types
type openRouterRequest struct {
	Model    string          `json:"model"`
	Messages []openRouterMsg `json:"messages"`
}

type openRouterMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// GenerateSummary generates an AI-powered summary for a user's monthly work logs
func GenerateSummary(userID int64, month string) (*model.WorkLogSummary, error) {
	// Validate month format (YYYY-MM)
	if !isValidMonthFormat(month) {
		return nil, errors.New("invalid month format, expected YYYY-MM")
	}

	// Fetch all work logs for the user in the specified month
	workLogs, err := getWorkLogsForMonth(userID, month)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch work logs: %w", err)
	}

	if len(workLogs) == 0 {
		return nil, errors.New("no work logs found for the specified month")
	}

	// Build content string from work logs
	content := buildWorkLogContent(workLogs)

	// Call OpenRouter API for summarization
	summary, err := callOpenRouter(content)
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}

	// Upsert the summary to database
	return work_log_summary_repo.Upsert(userID, month, summary)
}

// GetSummary retrieves an existing summary for a user's month
func GetSummary(userID int64, month string) (*model.WorkLogSummary, error) {
	if !isValidMonthFormat(month) {
		return nil, errors.New("invalid month format, expected YYYY-MM")
	}
	return work_log_summary_repo.GetByMonth(userID, month)
}

// isValidMonthFormat validates the month format (YYYY-MM)
func isValidMonthFormat(month string) bool {
	match, _ := regexp.MatchString(`^\d{4}-(0[1-9]|1[0-2])$`, month)
	return match
}

// getWorkLogsForMonth retrieves work logs for a specific month
func getWorkLogsForMonth(userID int64, month string) ([]model.WorkLog, error) {
	allLogs, err := work_log_repo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	var monthLogs []model.WorkLog
	for _, log := range allLogs {
		// log.Date is in YYYY-MM-DD format; check if it starts with the month
		if strings.HasPrefix(log.Date, month) {
			monthLogs = append(monthLogs, log)
		}
	}

	return monthLogs, nil
}

// buildWorkLogContent formats work logs into a string for the AI prompt
func buildWorkLogContent(logs []model.WorkLog) string {
	var sb strings.Builder
	sb.WriteString("Here are my daily work logs for this month:\n\n")
	for _, log := range logs {
		sb.WriteString(fmt.Sprintf("## %s\n%s\n\n", log.Date, log.Content))
	}
	return sb.String()
}

// callOpenRouter calls the OpenRouter API to generate a summary
func callOpenRouter(content string) (string, error) {
	cfg := config.Get()

	if cfg.OpenRouterAPIKey == "" {
		return "", errors.New("OPENROUTER_API_KEY is not configured")
	}

	prompt := fmt.Sprintf(`You are a helpful assistant that summarizes work logs.

Please provide a concise but comprehensive summary of the following monthly work activities.
Highlight key accomplishments, recurring themes, and notable patterns.
Format the summary in a clear, professional manner.

%s`, content)

	reqBody := openRouterRequest{
		Model: cfg.OpenRouterModel,
		Messages: []openRouterMsg{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.OpenRouterAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenRouter API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var openRouterResp openRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if openRouterResp.Error != nil {
		return "", fmt.Errorf("OpenRouter error: %s", openRouterResp.Error.Message)
	}

	if len(openRouterResp.Choices) == 0 {
		return "", errors.New("no response from OpenRouter")
	}

	return openRouterResp.Choices[0].Message.Content, nil
}
