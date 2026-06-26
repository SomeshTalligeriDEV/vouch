package response_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
)

func TestSuccess_SetsSuccessTrue(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/", func(c *fiber.Ctx) error {
		return response.OK(c, map[string]string{"msg": "ok"})
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var out struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatal(err)
	}
	if !out.Success {
		t.Fatal("expected success: true")
	}
}

func TestError_SetsSuccessFalse(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/err", func(c *fiber.Ctx) error {
		return response.Error(c, fiber.StatusBadRequest, "bad_request", "something went wrong")
	})

	req := httptest.NewRequest("GET", "/err", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var out struct {
		Success bool `json:"success"`
		Err     struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatal(err)
	}
	if out.Success {
		t.Fatal("expected success: false")
	}
	if out.Err.Code != "bad_request" {
		t.Fatalf("expected code 'bad_request', got %q", out.Err.Code)
	}
}
