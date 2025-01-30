package dto

import "time"

type ProxyClientDTO struct {
	ID           int64
	Title        string
	OS           string
	DownloadLink string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
