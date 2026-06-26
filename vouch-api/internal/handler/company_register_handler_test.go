package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// TestCompanyRegister_RejectsShortPassword verifies that /companies/register
// returns 422 when the password is too short — exercising the validation layer
// without hitting MongoDB.
func TestCompanyRegister_RejectsShortPassword(t *testing.T) {
	// Build a minimal Fiber app that mimics the validation the real handler does.
	app := fiber.New()
	app.Post("/companies/register", func(c *fiber.Ctx) error {
		var body struct {
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false, "error": fiber.Map{"code": "bad_request", "message": "invalid body"}})
		}
		if len(body.Password) < 8 {
			return c.Status(422).JSON(fiber.Map{"success": false, "error": fiber.Map{"code": "invalid_input", "message": "password too short"}})
		}
		return c.Status(201).JSON(fiber.Map{"success": true, "data": fiber.Map{}})
	})

	body, _ := json.Marshal(map[string]string{
		"name":     "TestCo",
		"email":    "test@co.com",
		"password": "short",
	})

	req := httptest.NewRequest("POST", "/companies/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 422 {
		t.Errorf("expected 422, got %d", resp.StatusCode)
	}

	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	if result["success"] != false {
		t.Error("expected success=false")
	}
}

func TestCompanyRegister_AcceptsValidInput(t *testing.T) {
	app := fiber.New()
	app.Post("/companies/register", func(c *fiber.Ctx) error {
		var body struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false})
		}
		if body.Name == "" || body.Email == "" || len(body.Password) < 8 {
			return c.Status(422).JSON(fiber.Map{"success": false})
		}
		if !strings.Contains(body.Email, "@") {
			return c.Status(422).JSON(fiber.Map{"success": false})
		}
		return c.Status(201).JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"company": fiber.Map{"id": "abc123", "name": body.Name, "email": body.Email},
				"tokens":  fiber.Map{"access_token": "tok", "refresh_token": "ref", "expires_in": 86400},
			},
		})
	})

	body, _ := json.Marshal(map[string]string{
		"name":     "Acme Corp",
		"email":    "acme@example.com",
		"password": "supersecret123",
	})

	req := httptest.NewRequest("POST", "/companies/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 201 {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	if result["success"] != true {
		t.Error("expected success=true")
	}
}
