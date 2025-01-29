package dto

type ProxyClientResponse struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	OS           string `json:"os"`
	DownloadLink string `json:"download_link"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
