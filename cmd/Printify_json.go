package cmd

type PRINTIFY_ImageUploadRequest struct {
	FileName string `json:"file_name"`
	URL      string `json:"url"`
}

type PRINTIFY_ImageUploadResponse struct {
	ID         string `json:"id"`
	FileName   string `json:"file_name"`
	Height     int    `json:"height"`
	Width      int    `json:"width"`
	Size       int    `json:"size"`
	MimeType   string `json:"mime_type"`
	PreviewURL string `json:"preview_url"`
	UploadTime string `json:"upload_time"`
}
