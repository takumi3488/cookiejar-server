package entity

import (
	"net/http"
	"testing"
	"time"
)

func TestNewCookie(t *testing.T) {
	tests := []struct {
		name       string
		httpCookie *http.Cookie
		want       *Cookie
	}{
		{
			name: "SameSite None",
			httpCookie: &http.Cookie{
				Name:     "test",
				Value:    "value",
				Domain:   ".example.com",
				Path:     "/",
				Expires:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			},
			want: &Cookie{
				Name:     "test",
				Value:    "value",
				Domain:   ".example.com",
				Path:     "/",
				Expires:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			},
		},
		{
			name: "SameSite Lax",
			httpCookie: &http.Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteLaxMode,
			},
			want: &Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteLaxMode,
			},
		},
		{
			name: "SameSite Strict",
			httpCookie: &http.Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteStrictMode,
			},
			want: &Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteStrictMode,
			},
		},
		{
			name: "SameSite Default",
			httpCookie: &http.Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteDefaultMode,
			},
			want: &Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteDefaultMode,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCookie(tt.httpCookie)
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
			if !got.Expires.Equal(tt.want.Expires) {
				t.Errorf("Expires = %v, want %v", got.Expires, tt.want.Expires)
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

func TestCookie_ToHTTPCookie(t *testing.T) {
	tests := []struct {
		name   string
		cookie *Cookie
		want   *http.Cookie
	}{
		{
			name: "SameSite none",
			cookie: &Cookie{
				Name:     "test",
				Value:    "value",
				Domain:   ".example.com",
				Path:     "/",
				Expires:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			},
			want: &http.Cookie{
				Name:     "test",
				Value:    "value",
				Domain:   ".example.com",
				Path:     "/",
				Expires:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			},
		},
		{
			name: "SameSite lax",
			cookie: &Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteLaxMode,
			},
			want: &http.Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteLaxMode,
			},
		},
		{
			name: "SameSite strict",
			cookie: &Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteStrictMode,
			},
			want: &http.Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteStrictMode,
			},
		},
		{
			name: "SameSite default",
			cookie: &Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteDefaultMode,
			},
			want: &http.Cookie{
				Name:     "test",
				Value:    "value",
				SameSite: http.SameSiteDefaultMode,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cookie.ToHTTPCookie()
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
			if !got.Expires.Equal(tt.want.Expires) {
				t.Errorf("Expires = %v, want %v", got.Expires, tt.want.Expires)
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
