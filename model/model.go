package model

import "time"

// User represents a user in the database
type User struct {
	ID         int64     `db:"id"`
	Email      string    `db:"email"`
	GoogleID   string    `db:"google_id"`
	Username   string    `db:"username"`
	Name       string    `db:"name"`
	PictureURL string    `db:"picture_url"`
	Role       string    `db:"role"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// GoogleUserInfo represents the user info from Google OAuth
type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

// JobApplication represents a job application in the database
type JobApplication struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"user_id"`
	CompanyName string    `db:"company_name"`
	JobTitle    string    `db:"job_title"`
	JobURL      string    `db:"job_url"`
	SalaryRange string    `db:"salary_range"`
	Email       string    `db:"email"`
	Notes       string    `db:"notes"`
	State       string    `db:"state"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// JobApplicationLog represents a log entry for a job application
type JobApplicationLog struct {
	ID               int64     `db:"id"`
	JobApplicationID int64     `db:"job_application_id"`
	ProcessName      string    `db:"process_name"`
	Note             string    `db:"note"`
	AudioURL         string    `db:"audio_url"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

// WorkLog represents a daily work log entry
type WorkLog struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Date      string    `db:"date"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
