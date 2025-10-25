package usecase

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/takumi3488/cookiejar-server/internal/domain/entity"
	"github.com/takumi3488/cookiejar-server/internal/domain/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	tracer := otel.Tracer("cookiejar-server/usecase")
	ctx, span := tracer.Start(ctx, "StoreCookies", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	span.SetAttributes(attribute.Int("cookie.count", len(cookies)))

	// ホストごとにCookieをグループ化
	hostCookies := make(map[string][]*entity.Cookie)
	for _, cookie := range cookies {
		c := entity.NewCookie(cookie)
		host := c.Domain
		hostCookies[host] = append(hostCookies[host], c)
	}

	// 各ホストごとに一括保存
	now := time.Now()
	for host, cookieList := range hostCookies {
		if err := u.cookieRepo.UpsertMany(ctx, host, cookieList, now); err != nil {
			log.Printf("Failed to upsert cookies for host=%s: %v", host, err)
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to upsert cookies")
			return err
		}
	}

	span.SetStatus(codes.Ok, "Successfully stored all cookies")
	return nil
}

func (u *cookieUsecase) GetAllCookies(ctx context.Context) ([]*entity.Cookie, error) {
	return u.cookieRepo.FindAll(ctx)
}

func (u *cookieUsecase) GetCookiesByHost(ctx context.Context, host string) ([]*entity.Cookie, error) {
	return u.cookieRepo.FindByHost(ctx, host)
}
