package dto

type ProxyClientRequest struct {
	Title        string `json:"title" validate:"required"`
	OS           string `json:"os" validate:"required,oneof=windows macos linux android ios"`
	DownloadLink string `json:"download_link" validate:"required,url"`
}
