package domain

import "time"

// Tier represents a builder's reputation band.
type Tier string

const (
	TierBronze   Tier = "Bronze"
	TierSilver   Tier = "Silver"
	TierGold     Tier = "Gold"
	TierPlatinum Tier = "Platinum"
	Tier24Karat  Tier = "24 Karat"
)

// Score component and multiplier constants. These define the Vouch scoring
// model and live in the domain so they are the single source of truth.
const (
	userScorePerUser    = 10.0
	userScoreCap        = 30000.0
	revenueScorePerMRR  = 2.0
	revenueScoreCap     = 20000.0
	impactScorePerPoint = 5.0
	impactScoreCap      = 15000.0
	velocityScorePerDel = 0.1
	velocityScoreCap    = 5000.0

	stripeMultiplierVerified   = 1.0
	stripeMultiplierUnverified = 0.6
)

// ScoreBreakdown holds the individual components that sum into the total.
type ScoreBreakdown struct {
	User     float64 `bson:"user" json:"user"`
	Revenue  float64 `bson:"revenue" json:"revenue"`
	Impact   float64 `bson:"impact" json:"impact"`
	Velocity float64 `bson:"velocity" json:"velocity"`
}

// BuilderScore is the computed reputation for a builder.
type BuilderScore struct {
	ID             string         `bson:"_id,omitempty" json:"id"`
	BuilderID      string         `bson:"builder_id" json:"builder_id"`
	TotalScore     float64        `bson:"total_score" json:"total_score"`
	Tier           Tier           `bson:"tier" json:"tier"`
	Breakdown      ScoreBreakdown `bson:"breakdown" json:"breakdown"`
	StripeVerified bool           `bson:"stripe_verified" json:"stripe_verified"`
	StripeMultiplier float64      `bson:"stripe_multiplier" json:"stripe_multiplier"`
	CalculatedAt   time.Time      `bson:"calculated_at" json:"calculated_at"`
	UpdatedAt      time.Time      `bson:"updated_at" json:"updated_at"`
}

// ScoreInputs are the aggregate metrics, across all of a builder's projects,
// required to compute a score.
type ScoreInputs struct {
	VerifiedUsers   int
	MRR             float64
	AverageRating   float64
	ReviewCount     int
	NinetyDayGrowth float64 // delta in verified users/revenue over the trailing 90 days
	StripeVerified  bool
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// ComputeScore applies the Vouch scoring model to the given inputs. This is a
// pure function — no IO, fully deterministic — so it can be unit tested in
// isolation and reused by both the service and the async worker.
func ComputeScore(in ScoreInputs) BuilderScore {
	breakdown := ScoreBreakdown{
		User:     min(float64(in.VerifiedUsers)*userScorePerUser, userScoreCap),
		Revenue:  min(in.MRR*revenueScorePerMRR, revenueScoreCap),
		Impact:   min(in.AverageRating*float64(in.ReviewCount)*impactScorePerPoint, impactScoreCap),
		Velocity: min(in.NinetyDayGrowth*velocityScorePerDel, velocityScoreCap),
	}

	multiplier := stripeMultiplierUnverified
	if in.StripeVerified {
		multiplier = stripeMultiplierVerified
	}

	subtotal := breakdown.User + breakdown.Revenue + breakdown.Impact + breakdown.Velocity
	total := subtotal * multiplier

	return BuilderScore{
		TotalScore:       total,
		Tier:             TierForScore(total),
		Breakdown:        breakdown,
		StripeVerified:   in.StripeVerified,
		StripeMultiplier: multiplier,
	}
}

// TierForScore maps a total score to its reputation tier.
func TierForScore(total float64) Tier {
	switch {
	case total >= 50000:
		return Tier24Karat
	case total >= 15000:
		return TierPlatinum
	case total >= 5000:
		return TierGold
	case total >= 1000:
		return TierSilver
	default:
		return TierBronze
	}
}
