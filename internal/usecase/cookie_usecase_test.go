package usecase

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/takumi3488/cookiejar-server/internal/domain/entity"
)

// モックリポジトリ
type mockCookieRepository struct {
	upsertFunc     func(ctx context.Context, cookie *entity.Cookie) error
	findAllFunc    func(ctx context.Context) ([]*entity.Cookie, error)
	findByHostFunc func(ctx context.Context, host string) ([]*entity.Cookie, error)
}

func (m *mockCookieRepository) Upsert(ctx context.Context, cookie *entity.Cookie) error {
	return m.upsertFunc(ctx, cookie)
}

func (m *mockCookieRepository) FindAll(ctx context.Context) ([]*entity.Cookie, error) {
	return m.findAllFunc(ctx)
}

func (m *mockCookieRepository) FindByHost(ctx context.Context, host string) ([]*entity.Cookie, error) {
	if m.findByHostFunc != nil {
		return m.findByHostFunc(ctx, host)
	}
	return nil, nil
}

func TestCookieUsecase_StoreCookies(t *testing.T) {
	tests := []struct {
		name      string
		cookies   []*http.Cookie
		upsertErr error
		wantErr   bool
	}{
		{
			name: "正常に保存できる",
			cookies: []*http.Cookie{
				{Name: "cookie1", Value: "value1"},
				{Name: "cookie2", Value: "value2"},
			},
			upsertErr: nil,
			wantErr:   false,
		},
		{
			name: "Upsertでエラーが発生",
			cookies: []*http.Cookie{
				{Name: "cookie1", Value: "value1"},
			},
			upsertErr: errors.New("upsert error"),
			wantErr:   true,
		},
		{
			name:      "空のCookieリスト",
			cookies:   []*http.Cookie{},
			upsertErr: nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockCookieRepository{
				upsertFunc: func(ctx context.Context, cookie *entity.Cookie) error {
					return tt.upsertErr
				},
			}

			uc := NewCookieUsecase(mockRepo)
			err := uc.StoreCookies(context.Background(), tt.cookies)

			if (err != nil) != tt.wantErr {
				t.Errorf("StoreCookies() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCookieUsecase_GetAllCookies(t *testing.T) {
	tests := []struct {
		name          string
		findAllResult []*entity.Cookie
		findAllErr    error
		wantErr       bool
		wantLen       int
	}{
		{
			name: "正常にCookieを取得できる",
			findAllResult: []*entity.Cookie{
				{Name: "cookie1", Value: "value1"},
				{Name: "cookie2", Value: "value2"},
			},
			findAllErr: nil,
			wantErr:    false,
			wantLen:    2,
		},
		{
			name:          "FindAllでエラーが発生",
			findAllResult: nil,
			findAllErr:    errors.New("find all error"),
			wantErr:       true,
			wantLen:       0,
		},
		{
			name:          "空のCookieリストを取得",
			findAllResult: []*entity.Cookie{},
			findAllErr:    nil,
			wantErr:       false,
			wantLen:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockCookieRepository{
				findAllFunc: func(ctx context.Context) ([]*entity.Cookie, error) {
					return tt.findAllResult, tt.findAllErr
				},
			}

			uc := NewCookieUsecase(mockRepo)
			result, err := uc.GetAllCookies(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllCookies() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && len(result) != tt.wantLen {
				t.Errorf("GetAllCookies() len = %v, want %v", len(result), tt.wantLen)
			}
		})
	}
}
