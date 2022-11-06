package cookies

import "context"

type Repository interface {
	GetCookies(ctx context.Context, sellerID uint64) (Cookies, error)
}
