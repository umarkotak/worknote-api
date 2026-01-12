package contract

import (
	"worknote-api/model"

	"github.com/go-jose/go-jose/v3/jwt"
)

// GoogleAuthRequest is the request body for Google authentication
type GoogleAuthRequest struct {
	IDToken string `json:"id_token"`
}

// AuthResponse is the response from authentication
type AuthResponse struct {
	AccessToken string      `json:"access_token"`
	User        *model.User `json:"user"`
}

// ErrorResponse is a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// TokenClaims represents the claims in the JWE token
type TokenClaims struct {
	jwt.Claims
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// UserInfo represents the authenticated user info in request context
type UserInfo struct {
	UserID int64
	Email  string
	Role   string
}

// CreateJobApplicationRequest is the request body for creating a job application
type CreateJobApplicationRequest struct {
	CompanyName string `json:"company_name"`
	JobTitle    string `json:"job_title"`
	JobURL      string `json:"job_url,omitempty"`
	SalaryRange string `json:"salary_range,omitempty"`
	Email       string `json:"email,omitempty"`
	Notes       string `json:"notes,omitempty"`
	State       string `json:"state,omitempty"`
}

// UpdateJobApplicationRequest is the request body for updating a job application
type UpdateJobApplicationRequest struct {
	CompanyName string `json:"company_name,omitempty"`
	JobTitle    string `json:"job_title,omitempty"`
	JobURL      string `json:"job_url,omitempty"`
	SalaryRange string `json:"salary_range,omitempty"`
	Email       string `json:"email,omitempty"`
	Notes       string `json:"notes,omitempty"`
	State       string `json:"state,omitempty"`
}

// JobApplicationResponse is the response for a job application
type JobApplicationResponse struct {
	ID          int64  `json:"id"`
	CompanyName string `json:"company_name"`
	JobTitle    string `json:"job_title"`
	JobURL      string `json:"job_url,omitempty"`
	SalaryRange string `json:"salary_range,omitempty"`
	Email       string `json:"email,omitempty"`
	Notes       string `json:"notes,omitempty"`
	State       string `json:"state"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// JobApplicationListResponse is the response for listing job applications
type JobApplicationListResponse struct {
	Data  []JobApplicationResponse `json:"data"`
	Total int                      `json:"total"`
}

// CreateJobApplicationLogRequest is the request body for creating a job application log
type CreateJobApplicationLogRequest struct {
	ProcessName string `json:"process_name"`
	Note        string `json:"note,omitempty"`
	AudioURL    string `json:"audio_url,omitempty"`
}

// UpdateJobApplicationLogRequest is the request body for updating a job application log
type UpdateJobApplicationLogRequest struct {
	ProcessName string `json:"process_name,omitempty"`
	Note        string `json:"note,omitempty"`
	AudioURL    string `json:"audio_url,omitempty"`
}

// JobApplicationLogResponse is the response for a job application log
type JobApplicationLogResponse struct {
	ID               int64  `json:"id"`
	JobApplicationID int64  `json:"job_application_id"`
	ProcessName      string `json:"process_name"`
	Note             string `json:"note,omitempty"`
	AudioURL         string `json:"audio_url,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// JobApplicationLogListResponse is the response for listing job application logs
type JobApplicationLogListResponse struct {
	Data []JobApplicationLogResponse `json:"data"`
}

// UpsertWorkLogRequest is the request body for upserting a work log
type UpsertWorkLogRequest struct {
	Date    string `json:"date"`
	Content string `json:"content"`
}

// WorkLogResponse is the response for a work log
type WorkLogResponse struct {
	ID        int64  `json:"id"`
	Date      string `json:"date"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// WorkLogListResponse is the response for listing work logs
type WorkLogListResponse struct {
	Data []WorkLogResponse `json:"data"`
}
