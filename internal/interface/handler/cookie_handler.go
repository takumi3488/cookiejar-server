package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/takumi3488/cookiejar-server/internal/usecase"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type CookieHandler struct {
	cookieUsecase usecase.CookieUsecase
}

func NewCookieHandler(cookieUsecase usecase.CookieUsecase) *CookieHandler {
	return &CookieHandler{
		cookieUsecase: cookieUsecase,
	}
}

type CookieRequest struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Path     string `json:"path,omitempty"`
	Domain   string `json:"domain,omitempty"`
	MaxAge   int    `json:"maxAge,omitempty"`
	Secure   bool   `json:"secure,omitempty"`
	HttpOnly bool   `json:"httpOnly,omitempty"`
	SameSite string `json:"sameSite,omitempty"`
}

func (c *CookieRequest) ToCookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Path:     c.Path,
		Domain:   c.Domain,
		MaxAge:   c.MaxAge,
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
	}
	switch c.SameSite {
	case "None":
		cookie.SameSite = http.SameSiteNoneMode
	case "Lax":
		cookie.SameSite = http.SameSiteLaxMode
	case "Strict":
		cookie.SameSite = http.SameSiteStrictMode
	}
	return cookie
}

func (h *CookieHandler) StoreCookies(c fiber.Ctx) error {
	ctx := c.Context()
	span := trace.SpanFromContext(ctx)

	var cookieReqs []*CookieRequest

	if err := json.Unmarshal(c.Body(), &cookieReqs); err != nil {
		log.Printf("Failed to parse JSON request body: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid JSON format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format or cookie structure",
		})
	}

	cookies := make([]*http.Cookie, len(cookieReqs))
	for i, req := range cookieReqs {
		cookies[i] = req.ToCookie()
	}

	if err := h.cookieUsecase.StoreCookies(ctx, cookies); err != nil {
		log.Printf("Failed to store cookies: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to store cookies")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store cookies",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"count":  len(cookies),
	})
}
