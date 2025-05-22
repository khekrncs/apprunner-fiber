package models

import (
	"time"
)

type FileInfo struct {
	Key          string            `json:"key"`
	UserID       string            `json:"user_id"`
	Filename     string            `json:"filename"`
	ContentType  string            `json:"content_type"`
	Size         int64             `json:"size"`
	ETag         string            `json:"etag,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	UploadedAt   time.Time         `json:"uploaded_at"`
	LastModified time.Time         `json:"last_modified"`
}

type UploadFileRequest struct {
	Filename    string            `json:"filename" validate:"required"`
	ContentType string            `json:"content_type"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}
