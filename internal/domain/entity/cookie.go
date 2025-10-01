package entity

import (
	"net/http"
	"time"
)

type Cookie struct {
	Name     string
	Value    string
	Domain   string
	Path     string
	Expires  time.Time
	Secure   bool
	HttpOnly bool
	SameSite string
}

func NewCookie(httpCookie *http.Cookie) *Cookie {
	sameSite := ""
	switch httpCookie.SameSite {
	case http.SameSiteDefaultMode:
		sameSite = "default"
	case http.SameSiteLaxMode:
		sameSite = "lax"
	case http.SameSiteStrictMode:
		sameSite = "strict"
	case http.SameSiteNoneMode:
		sameSite = "none"
	}

	return &Cookie{
		Name:     httpCookie.Name,
		Value:    httpCookie.Value,
		Domain:   httpCookie.Domain,
		Path:     httpCookie.Path,
		Expires:  httpCookie.Expires,
		Secure:   httpCookie.Secure,
		HttpOnly: httpCookie.HttpOnly,
		SameSite: sameSite,
	}
}

func (c *Cookie) ToHTTPCookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Domain:   c.Domain,
		Path:     c.Path,
		Expires:  c.Expires,
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
	}

	switch c.SameSite {
	case "lax":
		cookie.SameSite = http.SameSiteLaxMode
	case "strict":
		cookie.SameSite = http.SameSiteStrictMode
	case "none":
		cookie.SameSite = http.SameSiteNoneMode
	default:
		cookie.SameSite = http.SameSiteDefaultMode
	}

	return cookie
}
