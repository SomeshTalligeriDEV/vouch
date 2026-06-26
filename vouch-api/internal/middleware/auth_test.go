package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/middleware"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/jwt"
)

func newTestApp(jwtMgr *jwt.Manager, handlers ...fiber.Handler) *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/test", append([]fiber.Handler{middleware.Auth(jwtMgr)}, handlers...)...)
	return app
}

func TestAuth_MissingToken_Returns401(t *testing.T) {
	jwtMgr := jwt.NewManager("secret-32-characters-long-padding", "refresh-secret-32chars-longpadding")
	app := newTestApp(jwtMgr, func(c *fiber.Ctx) error { return c.SendStatus(200) })

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestAuth_InvalidToken_Returns401(t *testing.T) {
	jwtMgr := jwt.NewManager("secret-32-characters-long-padding", "refresh-secret-32chars-longpadding")
	app := newTestApp(jwtMgr, func(c *fiber.Ctx) error { return c.SendStatus(200) })

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestAuth_ValidToken_Passes(t *testing.T) {
	jwtMgr := jwt.NewManager("secret-32-characters-long-padding", "refresh-secret-32chars-longpadding")
	pair, err := jwtMgr.GenerateTyped("user1", "alice", "user", "user")
	if err != nil {
		t.Fatal(err)
	}

	app := newTestApp(jwtMgr, func(c *fiber.Ctx) error { return c.SendStatus(200) })

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRequireSubjectType_WrongType_Returns403(t *testing.T) {
	jwtMgr := jwt.NewManager("secret-32-characters-long-padding", "refresh-secret-32chars-longpadding")
	pair, err := jwtMgr.GenerateTyped("u1", "alice", "user", "user")
	if err != nil {
		t.Fatal(err)
	}

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/company-only", middleware.Auth(jwtMgr), middleware.RequireSubjectType("company"), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/company-only", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 403 {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}

func TestRequireRole_WrongRole_Returns403(t *testing.T) {
	jwtMgr := jwt.NewManager("secret-32-characters-long-padding", "refresh-secret-32chars-longpadding")
	pair, err := jwtMgr.GenerateTyped("u1", "alice", "user", "user")
	if err != nil {
		t.Fatal(err)
	}

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/admin-only", middleware.Auth(jwtMgr), middleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 403 {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}
