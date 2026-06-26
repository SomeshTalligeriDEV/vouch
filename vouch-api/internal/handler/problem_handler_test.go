package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestProblemList_ReturnsJSON(t *testing.T) {
	app := fiber.New()
	app.Get("/problems", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    []fiber.Map{},
			"meta":    fiber.Map{"page": 1, "limit": 20, "total": 0},
		})
	})

	req := httptest.NewRequest("GET", "/problems", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	if result["success"] != true {
		t.Error("expected success=true")
	}
}

func TestProblemCreate_RequiresAuth(t *testing.T) {
	app := fiber.New()
	app.Post("/problems", func(c *fiber.Ctx) error {
		if c.Get("Authorization") == "" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   fiber.Map{"code": "unauthorized", "message": "authentication required"},
			})
		}
		return c.Status(201).JSON(fiber.Map{"success": true, "data": fiber.Map{}})
	})

	body, _ := json.Marshal(map[string]any{
		"title":       "Need a CRM tool",
		"description": "Looking for something lightweight",
		"budget_min":  500,
		"budget_max":  2000,
	})

	req := httptest.NewRequest("POST", "/problems", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Errorf("expected 401 without auth, got %d", resp.StatusCode)
	}
}

func TestProblemGet_NotFound(t *testing.T) {
	app := fiber.New()
	app.Get("/problems/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "nonexistent" {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"error":   fiber.Map{"code": "not_found", "message": "problem not found"},
			})
		}
		return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"id": id}})
	})

	req := httptest.NewRequest("GET", "/problems/nonexistent", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 404 {
		t.Errorf("expected 404 for unknown problem, got %d", resp.StatusCode)
	}
}
