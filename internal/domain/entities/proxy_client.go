package entities

import "time"

type ProxyClient struct {
	ID           int64     `db:"id"`
	Title        string    `db:"title"`
	OS           string    `db:"os"`
	DownloadLink string    `db:"download_link"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
