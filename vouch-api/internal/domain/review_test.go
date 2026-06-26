package domain_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

func TestReview_IsValid_AcceptsRatings1To5(t *testing.T) {
	for rating := 1; rating <= 5; rating++ {
		r := domain.Review{Rating: rating}
		if !r.IsValid() {
			t.Fatalf("expected IsValid() true for rating %d", rating)
		}
	}
}

func TestReview_IsValid_RejectsOutOfRange(t *testing.T) {
	for _, rating := range []int{0, 6, -1} {
		r := domain.Review{Rating: rating}
		if r.IsValid() {
			t.Fatalf("expected IsValid() false for rating %d", rating)
		}
	}
}

func TestReview_IsVerifiedPurchase_TrueWhenSet(t *testing.T) {
	r := domain.Review{VerifiedPurchase: true}
	if !r.IsVerifiedPurchase() {
		t.Fatal("expected IsVerifiedPurchase() true")
	}
}
