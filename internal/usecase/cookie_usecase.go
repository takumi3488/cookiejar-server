package usecase

import (
	"context"
	"net/http"

	"github.com/takumi3488/cookiejar-server/internal/domain/entity"
	"github.com/takumi3488/cookiejar-server/internal/domain/repository"
)

type CookieUsecase interface {
	StoreCookies(ctx context.Context, cookies []*http.Cookie) error
	GetAllCookies(ctx context.Context) ([]*entity.Cookie, error)

	GetCookiesByHost(ctx context.Context, host string) ([]*entity.Cookie, error)
}

type cookieUsecase struct {
	cookieRepo repository.CookieRepository
}

func NewCookieUsecase(cookieRepo repository.CookieRepository) CookieUsecase {
	return &cookieUsecase{
		cookieRepo: cookieRepo,
	}
}

func (u *cookieUsecase) StoreCookies(ctx context.Context, cookies []*http.Cookie) error {
	for _, cookie := range cookies {
		c := entity.NewCookie(cookie)
		if err := u.cookieRepo.Upsert(ctx, c); err != nil {
			return err
		}
	}

	return nil
}

func (u *cookieUsecase) GetAllCookies(ctx context.Context) ([]*entity.Cookie, error) {
	return u.cookieRepo.FindAll(ctx)
}

func (u *cookieUsecase) GetCookiesByHost(ctx context.Context, host string) ([]*entity.Cookie, error) {
	return u.cookieRepo.FindByHost(ctx, host)
}
