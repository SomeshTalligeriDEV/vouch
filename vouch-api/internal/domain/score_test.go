package domain

import "testing"

func TestComputeScore_StripeMultiplier(t *testing.T) {
	in := ScoreInputs{
		VerifiedUsers:   100,  // 100 * 10 = 1000
		MRR:             500,  // 500 * 2 = 1000
		AverageRating:   4.0,  // 4 * 50 * 5 = 1000
		ReviewCount:     50,
		NinetyDayGrowth: 1000, // 1000 * 0.1 = 100
		StripeVerified:  true,
	}
	got := ComputeScore(in)

	wantSubtotal := 1000.0 + 1000.0 + 1000.0 + 100.0
	if got.TotalScore != wantSubtotal {
		t.Fatalf("verified total = %v, want %v", got.TotalScore, wantSubtotal)
	}
	if got.StripeMultiplier != 1.0 {
		t.Fatalf("verified multiplier = %v, want 1.0", got.StripeMultiplier)
	}

	in.StripeVerified = false
	gotUnverified := ComputeScore(in)
	if gotUnverified.TotalScore != wantSubtotal*0.6 {
		t.Fatalf("unverified total = %v, want %v", gotUnverified.TotalScore, wantSubtotal*0.6)
	}
}

func TestComputeScore_Caps(t *testing.T) {
	in := ScoreInputs{
		VerifiedUsers:   1_000_000, // cap 30000
		MRR:             1_000_000, // cap 20000
		AverageRating:   5,
		ReviewCount:     1_000_000, // cap 15000
		NinetyDayGrowth: 1_000_000, // cap 5000
		StripeVerified:  true,
	}
	got := ComputeScore(in)
	if got.Breakdown.User != 30000 || got.Breakdown.Revenue != 20000 ||
		got.Breakdown.Impact != 15000 || got.Breakdown.Velocity != 5000 {
		t.Fatalf("caps not applied: %+v", got.Breakdown)
	}
	if got.Tier != Tier24Karat {
		t.Fatalf("tier = %v, want 24 Karat", got.Tier)
	}
}

func TestTierForScore(t *testing.T) {
	cases := []struct {
		score float64
		want  Tier
	}{
		{0, TierBronze}, {999, TierBronze},
		{1000, TierSilver}, {4999, TierSilver},
		{5000, TierGold}, {14999, TierGold},
		{15000, TierPlatinum}, {49999, TierPlatinum},
		{50000, Tier24Karat},
	}
	for _, c := range cases {
		if got := TierForScore(c.score); got != c.want {
			t.Errorf("TierForScore(%v) = %v, want %v", c.score, got, c.want)
		}
	}
}
