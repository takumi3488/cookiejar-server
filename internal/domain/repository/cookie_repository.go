package repository

import (
	"context"
	"time"

	"github.com/takumi3488/cookiejar-server/internal/domain/entity"
)

type CookieRepository interface {
	Upsert(ctx context.Context, cookie *entity.Cookie, updatedAt time.Time) error
	UpsertMany(ctx context.Context, host string, cookies []*entity.Cookie, updatedAt time.Time) error
	FindAll(ctx context.Context) ([]*entity.Cookie, error)

	FindByHost(ctx context.Context, host string) ([]*entity.Cookie, error)
}
