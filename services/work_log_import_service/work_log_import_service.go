package work_log_import_service

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"worknote-api/repos/work_log_repo"
)

var (
	// Matches date formats like "Fri, 21 November 2025" or "Monday, 1 January 2025"
	dateLineRegex = regexp.MustCompile(`^[A-Za-z]+,\s*\d{1,2}\s+[A-Za-z]+\s+\d{4}$`)
)

// ImportedWorkLog represents a worklog parsed from markdown
type ImportedWorkLog struct {
	Date    string
	Content string
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	Imported int      `json:"imported"`
	Updated  int      `json:"updated"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors,omitempty"`
}

// ParseMarkdown parses markdown content and extracts worklogs
func ParseMarkdown(content string) ([]ImportedWorkLog, error) {
	lines := strings.Split(content, "\n")
	var worklogs []ImportedWorkLog
	var currentDate string
	var currentContent []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			// Empty line signifies end of current worklog entry
			if currentDate != "" && len(currentContent) > 0 {
				worklogs = append(worklogs, ImportedWorkLog{
					Date:    currentDate,
					Content: strings.Join(currentContent, "\n"),
				})
			}
			currentDate = ""
			currentContent = nil
			continue
		}

		// Check if this is a date line
		if dateLineRegex.MatchString(line) {
			// Save previous worklog if exists
			if currentDate != "" && len(currentContent) > 0 {
				worklogs = append(worklogs, ImportedWorkLog{
					Date:    currentDate,
					Content: strings.Join(currentContent, "\n"),
				})
			}
			// Parse new date
			parsedDate, err := parseDateLine(line)
			if err != nil {
				return nil, errors.New("failed to parse date: " + line)
			}
			currentDate = parsedDate
			currentContent = nil
			continue
		}

		// Check if this is a bullet point line
		if strings.HasPrefix(line, "-") {
			// Remove bullet point and leading whitespace
			content := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			if content != "" {
				currentContent = append(currentContent, content)
			}
		} else if currentDate != "" {
			// Line without bullet point but within a worklog entry
			currentContent = append(currentContent, line)
		}
	}

	// Don't forget the last entry
	if currentDate != "" && len(currentContent) > 0 {
		worklogs = append(worklogs, ImportedWorkLog{
			Date:    currentDate,
			Content: strings.Join(currentContent, "\n"),
		})
	}

	if len(worklogs) == 0 {
		return nil, errors.New("no valid worklogs found in markdown")
	}

	return worklogs, nil
}

// parseDateLine parses a date line like "Fri, 21 November 2025" and returns YYYY-MM-DD
func parseDateLine(dateStr string) (string, error) {
	// Try parsing with various layouts
	layouts := []string{
		"Mon, 2 January 2006",
		"Monday, 2 January 2006",
		"Mon, 2 Jan 2006",
		"Monday, 2 Jan 2006",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return t.Format("2006-01-02"), nil
		}
	}

	return "", errors.New("unable to parse date: " + dateStr)
}

// ImportWorkLogs imports parsed worklogs for a user
func ImportWorkLogs(userID int64, worklogs []ImportedWorkLog) (*ImportResult, error) {
	result := &ImportResult{
		Errors: make([]string, 0),
	}

	for _, wl := range worklogs {
		// Check if worklog already exists
		existing, err := work_log_repo.GetByDate(userID, wl.Date)
		if err != nil {
			result.Errors = append(result.Errors, "error checking date "+wl.Date+": "+err.Error())
			continue
		}

		if existing != nil {
			// Update existing
			_, err = work_log_repo.Upsert(userID, wl.Date, wl.Content)
			if err != nil {
				result.Errors = append(result.Errors, "error updating "+wl.Date+": "+err.Error())
			} else {
				result.Updated++
			}
		} else {
			// Create new
			_, err = work_log_repo.Upsert(userID, wl.Date, wl.Content)
			if err != nil {
				result.Errors = append(result.Errors, "error creating "+wl.Date+": "+err.Error())
			} else {
				result.Imported++
			}
		}
	}

	return result, nil
}

// ImportFromMarkdown imports worklogs directly from markdown content
func ImportFromMarkdown(userID int64, markdownContent string) (*ImportResult, error) {
	worklogs, err := ParseMarkdown(markdownContent)
	if err != nil {
		return nil, err
	}

	return ImportWorkLogs(userID, worklogs)
}
