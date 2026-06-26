package pagination_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/pagination"
)

func TestClamp_ValidInput(t *testing.T) {
	page, limit := pagination.Clamp(2, 10)
	if page != 2 || limit != 10 {
		t.Errorf("expected (2,10), got (%d,%d)", page, limit)
	}
}

func TestClamp_ZeroPage(t *testing.T) {
	page, _ := pagination.Clamp(0, 20)
	if page != 1 {
		t.Errorf("expected page=1, got %d", page)
	}
}

func TestClamp_NegativePage(t *testing.T) {
	page, _ := pagination.Clamp(-5, 20)
	if page != 1 {
		t.Errorf("expected page=1, got %d", page)
	}
}

func TestClamp_LimitExceedsMax(t *testing.T) {
	_, limit := pagination.Clamp(1, 9999)
	if limit != pagination.MaxLimit {
		t.Errorf("expected limit=%d, got %d", pagination.MaxLimit, limit)
	}
}

func TestClamp_ZeroLimit(t *testing.T) {
	_, limit := pagination.Clamp(1, 0)
	if limit != pagination.DefaultLimit {
		t.Errorf("expected limit=%d, got %d", pagination.DefaultLimit, limit)
	}
}

func TestOffset_PageOne(t *testing.T) {
	off := pagination.Offset(1, 20)
	if off != 0 {
		t.Errorf("page 1 offset should be 0, got %d", off)
	}
}

func TestOffset_PageThree(t *testing.T) {
	off := pagination.Offset(3, 10)
	if off != 20 {
		t.Errorf("page 3 limit 10 offset should be 20, got %d", off)
	}
}

func TestHasNextPage_True(t *testing.T) {
	if !pagination.HasNextPage(1, 10, 25) {
		t.Error("expected HasNextPage=true when total > page*limit")
	}
}

func TestHasNextPage_False(t *testing.T) {
	if pagination.HasNextPage(3, 10, 25) {
		t.Error("expected HasNextPage=false on last page")
	}
}
