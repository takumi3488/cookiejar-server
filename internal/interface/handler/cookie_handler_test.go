package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/takumi3488/cookiejar-server/internal/domain/entity"
)

// モックユースケース
type mockCookieUsecase struct {
	storeCookiesFunc  func(ctx context.Context, cookies []*http.Cookie) error
	getAllCookiesFunc func(ctx context.Context) ([]*entity.Cookie, error)
}

func (m *mockCookieUsecase) StoreCookies(ctx context.Context, cookies []*http.Cookie) error {
	return m.storeCookiesFunc(ctx, cookies)
}

func (m *mockCookieUsecase) GetAllCookies(ctx context.Context) ([]*entity.Cookie, error) {
	return m.getAllCookiesFunc(ctx)
}

func TestCookieHandler_StoreCookies(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		storeCookiesErr error
		wantStatus     int
		wantResponse   map[string]interface{}
	}{
		{
			name: "正常にCookieを保存できる",
			requestBody: []*CookieRequest{
				{
					Name:     "test_cookie",
					Value:    "test_value",
					Domain:   ".example.com",
					Path:     "/",
					Secure:   true,
					HttpOnly: true,
					SameSite: "None",
				},
			},
			storeCookiesErr: nil,
			wantStatus:      200,
			wantResponse: map[string]interface{}{
				"status": "success",
				"count":  float64(1),
			},
		},
		{
			name: "複数のCookieを保存できる",
			requestBody: []*CookieRequest{
				{Name: "cookie1", Value: "value1"},
				{Name: "cookie2", Value: "value2"},
				{Name: "cookie3", Value: "value3"},
			},
			storeCookiesErr: nil,
			wantStatus:      200,
			wantResponse: map[string]interface{}{
				"status": "success",
				"count":  float64(3),
			},
		},
		{
			name:           "不正なJSON形式",
			requestBody:    "invalid json",
			storeCookiesErr: nil,
			wantStatus:      400,
			wantResponse: map[string]interface{}{
				"error": "Invalid JSON format or cookie structure",
			},
		},
		{
			name: "StoreCookiesでエラーが発生",
			requestBody: []*CookieRequest{
				{Name: "test_cookie", Value: "test_value"},
			},
			storeCookiesErr: errors.New("store error"),
			wantStatus:      500,
			wantResponse: map[string]interface{}{
				"error": "Failed to store cookies",
			},
		},
		{
			name:           "空のCookieリスト",
			requestBody:    []*CookieRequest{},
			storeCookiesErr: nil,
			wantStatus:      200,
			wantResponse: map[string]interface{}{
				"status": "success",
				"count":  float64(0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックユースケース作成
			mockUsecase := &mockCookieUsecase{
				storeCookiesFunc: func(ctx context.Context, cookies []*http.Cookie) error {
					return tt.storeCookiesErr
				},
			}

			// ハンドラー作成
			handler := NewCookieHandler(mockUsecase)

			// Fiberアプリケーション作成
			app := fiber.New()
			app.Post("/", handler.StoreCookies)

			// リクエストボディ作成
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			// HTTPリクエスト作成
			req, _ := http.NewRequest("POST", "/", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// リクエスト実行
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to execute request: %v", err)
			}

			// ステータスコード確認
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Status code = %v, want %v", resp.StatusCode, tt.wantStatus)
			}

			// レスポンスボディ確認
			respBody, _ := io.ReadAll(resp.Body)
			var gotResponse map[string]interface{}
			if err := json.Unmarshal(respBody, &gotResponse); err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			for key, wantValue := range tt.wantResponse {
				if gotValue, ok := gotResponse[key]; !ok {
					t.Errorf("Response missing key %v", key)
				} else if gotValue != wantValue {
					t.Errorf("Response[%v] = %v, want %v", key, gotValue, wantValue)
				}
			}
		})
	}
}

func TestCookieRequest_ToCookie(t *testing.T) {
	tests := []struct {
		name    string
		request *CookieRequest
		want    *http.Cookie
	}{
		{
			name: "SameSite None",
			request: &CookieRequest{
				Name:     "test",
				Value:    "value",
				Domain:   ".example.com",
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: "None",
			},
			want: &http.Cookie{
				Name:     "test",
				Value:    "value",
				Domain:   ".example.com",
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			},
		},
		{
			name: "SameSite Lax",
			request: &CookieRequest{
				Name:     "test",
				Value:    "value",
				SameSite: "Lax",
			},
			want: &http.Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteLaxMode,
			},
		},
		{
			name: "SameSite Strict",
			request: &CookieRequest{
				Name:     "test",
				Value:    "value",
				SameSite: "Strict",
			},
			want: &http.Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteStrictMode,
			},
		},
		{
			name: "SameSite未指定",
			request: &CookieRequest{
				Name:  "test",
				Value: "value",
			},
			want: &http.Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.ToCookie()

			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Value != tt.want.Value {
				t.Errorf("Value = %v, want %v", got.Value, tt.want.Value)
			}
			if got.Domain != tt.want.Domain {
				t.Errorf("Domain = %v, want %v", got.Domain, tt.want.Domain)
			}
			if got.Path != tt.want.Path {
				t.Errorf("Path = %v, want %v", got.Path, tt.want.Path)
			}
			if got.Secure != tt.want.Secure {
				t.Errorf("Secure = %v, want %v", got.Secure, tt.want.Secure)
			}
			if got.HttpOnly != tt.want.HttpOnly {
				t.Errorf("HttpOnly = %v, want %v", got.HttpOnly, tt.want.HttpOnly)
			}
			if got.SameSite != tt.want.SameSite {
				t.Errorf("SameSite = %v, want %v", got.SameSite, tt.want.SameSite)
			}
		})
	}
}
