package work_log_download_service

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"worknote-api/model"
	"worknote-api/repos/work_log_repo"
)

var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// DownloadRequest represents the request parameters for downloading worklogs
type DownloadRequest struct {
	StartDate string
	EndDate   string
}

// Validate validates the download request
func (r *DownloadRequest) Validate() error {
	if r.StartDate == "" {
		return errors.New("start_date is required")
	}
	if r.EndDate == "" {
		return errors.New("end_date is required")
	}
	if !dateRegex.MatchString(r.StartDate) {
		return errors.New("start_date must be in YYYY-MM-DD format")
	}
	if !dateRegex.MatchString(r.EndDate) {
		return errors.New("end_date must be in YYYY-MM-DD format")
	}
	if r.StartDate > r.EndDate {
		return errors.New("start_date must be before or equal to end_date")
	}
	return nil
}

// DownloadWorkLogs retrieves worklogs within a date range and generates markdown
func DownloadWorkLogs(userID int64, req *DownloadRequest) (string, string, error) {
	if err := req.Validate(); err != nil {
		return "", "", err
	}

	logs, err := work_log_repo.ListByUserIDAndDateRange(userID, req.StartDate, req.EndDate)
	if err != nil {
		return "", "", err
	}

	markdown := GenerateMarkdown(logs)
	filename := GenerateFilename(req.StartDate, req.EndDate)

	return markdown, filename, nil
}

// GenerateMarkdown converts worklogs to markdown format
func GenerateMarkdown(logs []model.WorkLog) string {
	if len(logs) == 0 {
		return ""
	}

	var sb strings.Builder

	// Try multiple date formats
	dateFormats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}

	for _, log := range logs {
		// Parse date and format as "Fri, 21 November 2025"
		var parsedDate time.Time
		var err error
		for _, format := range dateFormats {
			parsedDate, err = time.Parse(format, log.Date)
			if err == nil {
				break
			}
		}

		if err != nil {
			// Fallback to original date if parsing fails
			sb.WriteString(log.Date)
			sb.WriteString("\n")
		} else {
			sb.WriteString(parsedDate.Format("Mon, 2 January 2006"))
			sb.WriteString("\n")
		}

		// Split content by lines and add bullet points
		lines := strings.Split(log.Content, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// If line already starts with "-", don't add another
			if strings.HasPrefix(line, "-") {
				sb.WriteString(line)
			} else {
				sb.WriteString("- ")
				sb.WriteString(line)
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// GenerateFilename creates a filename for the download
func GenerateFilename(startDate, endDate string) string {
	// Format: worklog-2024-01-01-to-2024-01-31.md
	return "worklog-" + startDate + "-to-" + endDate + ".md"
}
