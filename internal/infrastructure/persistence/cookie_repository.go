package persistence

import (
	"context"
	"encoding/json"
	"time"

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

func (r *cookieRepository) Upsert(ctx context.Context, cookie *entity.Cookie, updatedAt time.Time) error {
	// 既存のCookieを取得
	existingCookies, err := r.FindByHost(ctx, cookie.Domain)
	if err != nil {
		// レコードが存在しない場合は空配列として扱う
		existingCookies = []*entity.Cookie{}
	}

	// 同じ名前のCookieを探してマージ（更新）、なければ追加
	found := false
	for i, existing := range existingCookies {
		if existing.Name == cookie.Name {
			existingCookies[i] = cookie
			found = true
			break
		}
	}
	if !found {
		existingCookies = append(existingCookies, cookie)
	}

	// Cookie配列をJSON化
	cookiesJSON, err := json.Marshal(existingCookies)
	if err != nil {
		return err
	}

	return r.queries.UpsertCookies(ctx, db.UpsertCookiesParams{
		Host:      cookie.Domain,
		Cookies:   string(cookiesJSON),
		UpdatedAt: updatedAt,
	})
}

func (r *cookieRepository) UpsertMany(ctx context.Context, host string, cookies []*entity.Cookie, updatedAt time.Time) error {
	// 既存のCookieを取得
	existingCookies, err := r.FindByHost(ctx, host)
	if err != nil {
		// レコードが存在しない場合は空配列として扱う
		existingCookies = []*entity.Cookie{}
	}

	// 既存のCookieをマップ化（名前をキーに）
	existingMap := make(map[string]*entity.Cookie)
	for _, existing := range existingCookies {
		existingMap[existing.Name] = existing
	}

	// 新しいCookieでマップを更新
	for _, cookie := range cookies {
		existingMap[cookie.Name] = cookie
	}

	// マップから配列を再構築
	mergedCookies := make([]*entity.Cookie, 0, len(existingMap))
	for _, cookie := range existingMap {
		mergedCookies = append(mergedCookies, cookie)
	}

	// Cookie配列をJSON化
	cookiesJSON, err := json.Marshal(mergedCookies)
	if err != nil {
		return err
	}

	return r.queries.UpsertCookies(ctx, db.UpsertCookiesParams{
		Host:      host,
		Cookies:   string(cookiesJSON),
		UpdatedAt: updatedAt,
	})
}

func (r *cookieRepository) FindAll(ctx context.Context) ([]*entity.Cookie, error) {
	cookies, err := r.queries.ListCookies(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Cookie, 0)
	for _, c := range cookies {
		var cookieList []*entity.Cookie
		if err := json.Unmarshal([]byte(c.Cookies), &cookieList); err != nil {
			// 配列としてのアンマーシャルが失敗した場合、単一のCookieとして試す（後方互換性）
			var cookie entity.Cookie
			if err := json.Unmarshal([]byte(c.Cookies), &cookie); err != nil {
				// アンマーシャルが失敗した場合、このCookieをスキップ
				continue
			}
			result = append(result, &cookie)
		} else {
			result = append(result, cookieList...)
		}
	}
	return result, nil
}

func (r *cookieRepository) FindByHost(ctx context.Context, host string) ([]*entity.Cookie, error) {
	c, err := r.queries.GetCookiesByHost(ctx, host)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Cookie, 0)
	var cookieList []*entity.Cookie
	if err := json.Unmarshal([]byte(c.Cookies), &cookieList); err != nil {
		// 配列としてのアンマーシャルが失敗した場合、単一のCookieとして試す（後方互換性）
		var cookie entity.Cookie
		if err := json.Unmarshal([]byte(c.Cookies), &cookie); err != nil {
			return nil, err
		}
		result = append(result, &cookie)
	} else {
		result = append(result, cookieList...)
	}
	return result, nil
}
