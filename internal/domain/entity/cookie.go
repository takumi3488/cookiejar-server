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
	SameSite http.SameSite
}

func NewCookie(httpCookie *http.Cookie) *Cookie {
	return &Cookie{
		Name:     httpCookie.Name,
		Value:    httpCookie.Value,
		Domain:   httpCookie.Domain,
		Path:     httpCookie.Path,
		Expires:  httpCookie.Expires,
		Secure:   httpCookie.Secure,
		HttpOnly: httpCookie.HttpOnly,
		SameSite: httpCookie.SameSite,
	}
}

func (c *Cookie) ToHTTPCookie() *http.Cookie {
	return &http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Domain:   c.Domain,
		Path:     c.Path,
		Expires:  c.Expires,
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
		SameSite: c.SameSite,
	}
}
