package models

// ContentDetail has uploaded content information
type ContentDetail struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Size     string `json:"size,omitempty"`
	UserID   string `json:"userId,omitempty"`
	UploadID string `json:"uploadId,omitempty"`
}

// ErrorMessage is returned as response when there is POST failure
type ErrorMessage struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// UserAccess defines access for each content
type UserAccess struct {
	ContentID string   `json:"contentId"`
	UserIds   []string `json:"userIds"`
}

// DownloadContent defines content to download
type DownloadContent struct {
	UserID      string `json:"userId"`
	ContentID   string `json:"contentId"`
	ContentName string `json:"contentName,omitempty"`
	Content     string `json:"content,omitempty"`
}
