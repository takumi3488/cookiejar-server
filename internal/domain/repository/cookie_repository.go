package repository

import (
	"context"

	"github.com/takumi3488/cookiejar-server/internal/domain/entity"
)

type CookieRepository interface {
	Upsert(ctx context.Context, cookie *entity.Cookie) error
	FindAll(ctx context.Context) ([]*entity.Cookie, error)
}
