package dto

type AssetUploadResponse struct {
	Path     string `json:"path"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
	Category string `json:"category"`
}
