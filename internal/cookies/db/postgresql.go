package cookiesdb

import (
	"context"
	"wb-report-downloader/internal/cookies"
	"wb-report-downloader/pkg/client/postgresql"
)

type repository struct {
	client postgresql.Client
}

func (r *repository) GetCookies (ctx context.Context, sellerID uint64) (cookies.Cookies, error) {
	q := `
		SELECT access_token
		FROM sellers
		WHERE id = $1
	`

	var c cookies.Cookies
	err := r.client.QueryRow(ctx, q, sellerID).Scan(&c.RawCookies)
	if err != nil {
		return cookies.Cookies{}, err
	}
	return c, nil
}

func NewRepository(client postgresql.Client) cookies.Repository {
	return &repository{
		client: client,
	}
}
