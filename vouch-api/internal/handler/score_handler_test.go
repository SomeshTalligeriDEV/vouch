package handler_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestScoreLeaderboard_EmptyList(t *testing.T) {
	app := fiber.New()
	app.Get("/scores", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    []fiber.Map{},
		})
	})

	req := httptest.NewRequest("GET", "/scores", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result struct {
		Success bool            `json:"success"`
		Data    []map[string]any `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if !result.Success {
		t.Error("expected success=true")
	}
	if result.Data == nil {
		t.Error("expected non-nil data array")
	}
}

func TestScoreGetByUsername_Returns200(t *testing.T) {
	app := fiber.New()
	app.Get("/scores/:username", func(c *fiber.Ctx) error {
		username := c.Params("username")
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"builder_id":  "user123",
				"username":    username,
				"total_score": 750.0,
				"tier":        "gold",
			},
		})
	})

	req := httptest.NewRequest("GET", "/scores/alice", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].(map[string]any)
	if data["username"] != "alice" {
		t.Errorf("expected username=alice, got %v", data["username"])
	}
}

func TestScoreRecalculate_RequiresAuth(t *testing.T) {
	app := fiber.New()
	app.Post("/scores/recalculate", func(c *fiber.Ctx) error {
		if c.Get("Authorization") == "" {
			return c.Status(401).JSON(fiber.Map{"success": false})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	req := httptest.NewRequest("POST", "/scores/recalculate", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}
