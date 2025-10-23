package persistence

import (
	"context"
	"encoding/json"

	"github.com/takumi3488/cookiejar-server/db"
	"github.com/takumi3488/cookiejar-server/internal/domain/entity"
	"github.com/takumi3488/cookiejar-server/internal/domain/repository"
)

type cookieRepository struct {
	queries *db.Queries
}

func NewCookieRepository(queries *db.Queries) repository.CookieRepository {
	return &cookieRepository{
		queries: queries,
	}
}

func (r *cookieRepository) Upsert(ctx context.Context, cookie *entity.Cookie) error {
	// 今のところドメインをホストとして使用し、Cookie保存にJSONを使用
	cookieJSON, err := json.Marshal(cookie)
	if err != nil {
		return err
	}

	return r.queries.UpsertCookies(ctx, db.UpsertCookiesParams{
		Host:    cookie.Domain,
		Cookies: string(cookieJSON),
	})
}

func (r *cookieRepository) FindAll(ctx context.Context) ([]*entity.Cookie, error) {
	cookies, err := r.queries.ListCookies(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Cookie, 0)
	for _, c := range cookies {
		var cookie entity.Cookie
		if err := json.Unmarshal([]byte(c.Cookies), &cookie); err != nil {
			// アンマーシャルが失敗した場合、このCookieをスキップ
			continue
		}
		result = append(result, &cookie)
	}
	return result, nil
}
