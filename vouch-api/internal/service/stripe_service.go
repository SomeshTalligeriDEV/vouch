package service

import (
	"context"
	"fmt"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// StripeGateway abstracts the external Stripe API. Vouch uses read-only OAuth:
// it exchanges an auth code for a connected account id and reads revenue. It
// never creates charges.
type StripeGateway interface {
	// ExchangeCode swaps an OAuth authorization code for a connected Stripe
	// account id.
	ExchangeCode(ctx context.Context, code string) (accountID string, err error)
	// FetchRevenue reads current MRR and customer counts for a connected
	// account.
	FetchRevenue(ctx context.Context, accountID string) (mrr float64, customers int, currency string, err error)
}

// StripeService holds business logic for Stripe connection and revenue sync.
type StripeService struct {
	users    domain.UserRepository
	stripe   domain.StripeRepository
	gateway  StripeGateway
	scoreEnq ScoreEnqueuer
}

// NewStripeService constructs a StripeService.
func NewStripeService(
	users domain.UserRepository,
	stripeRepo domain.StripeRepository,
	gateway StripeGateway,
	scoreEnq ScoreEnqueuer,
) *StripeService {
	return &StripeService{users: users, stripe: stripeRepo, gateway: gateway, scoreEnq: scoreEnq}
}

// Connect completes the Stripe OAuth flow for a user: it exchanges the code,
// stores the connected account id, takes an initial revenue snapshot, and
// schedules a score recalculation.
func (s *StripeService) Connect(ctx context.Context, userID, code string) (*domain.StripeSnapshot, error) {
	accountID, err := s.gateway.ExchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("StripeService.Connect: %w", err)
	}
	if err := s.users.SetStripeAccount(ctx, userID, accountID); err != nil {
		return nil, fmt.Errorf("StripeService.Connect: %w", err)
	}
	snap, err := s.Sync(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("StripeService.Connect: %w", err)
	}
	return snap, nil
}

// Sync reads current revenue for a user's connected account and persists a
// fresh snapshot. Invoked synchronously on connect and from the worker.
func (s *StripeService) Sync(ctx context.Context, userID string) (*domain.StripeSnapshot, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("StripeService.Sync: %w", err)
	}
	if !user.HasStripe() {
		return nil, fmt.Errorf("StripeService.Sync: %w", domain.ErrStripeNotVerified)
	}

	mrr, customers, currency, err := s.gateway.FetchRevenue(ctx, user.StripeAccountID)
	if err != nil {
		return nil, fmt.Errorf("StripeService.Sync: %w", err)
	}

	snap := &domain.StripeSnapshot{
		BuilderID:      userID,
		MRR:            mrr,
		TotalCustomers: customers,
		Currency:       currency,
	}
	if err := s.stripe.Save(ctx, snap); err != nil {
		return nil, fmt.Errorf("StripeService.Sync: %w", err)
	}

	if s.scoreEnq != nil {
		_ = s.scoreEnq.EnqueueScoreRecalc(ctx, userID)
	}
	return snap, nil
}
