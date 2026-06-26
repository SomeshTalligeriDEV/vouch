package domain_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

func TestProblem_IsOpen(t *testing.T) {
	p := domain.Problem{Status: domain.ProblemStatusOpen}
	if !p.IsOpen() {
		t.Fatal("expected IsOpen() true for open problem")
	}
}

func TestProblem_IsClaimed(t *testing.T) {
	p := domain.Problem{Status: domain.ProblemStatusClaimed}
	if !p.IsClaimed() {
		t.Fatal("expected IsClaimed() true")
	}
}

func TestProblem_IsShipped(t *testing.T) {
	p := domain.Problem{Status: domain.ProblemStatusShipped}
	if !p.IsShipped() {
		t.Fatal("expected IsShipped() true")
	}
}

func TestProblem_IsPostedBy(t *testing.T) {
	p := domain.Problem{PosterID: "company1"}
	if !p.IsPostedBy("company1") {
		t.Fatal("expected IsPostedBy('company1') true")
	}
	if p.IsPostedBy("other") {
		t.Fatal("expected IsPostedBy('other') false")
	}
}
