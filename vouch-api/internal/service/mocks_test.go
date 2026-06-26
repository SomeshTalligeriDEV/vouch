package service

import (
	"context"
	"sync"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// In-memory fakes implementing the domain repository interfaces, used to test
// service business logic without a database.

type fakeUserRepo struct {
	mu    sync.Mutex
	byID  map[string]*domain.User
	byUN  map[string]*domain.User
	seq   int
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{byID: map[string]*domain.User{}, byUN: map[string]*domain.User{}}
}

func (f *fakeUserRepo) add(u *domain.User) *domain.User {
	f.seq++
	if u.ID == "" {
		u.ID = "u" + itoa(f.seq)
	}
	f.byID[u.ID] = u
	f.byUN[u.Username] = u
	return u
}

func (f *fakeUserRepo) Create(_ context.Context, u *domain.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.byUN[u.Username]; ok {
		return domain.ErrAlreadyExists
	}
	f.add(u)
	return nil
}
func (f *fakeUserRepo) GetByID(_ context.Context, id string) (*domain.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.byID[id]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}
func (f *fakeUserRepo) GetByUsername(_ context.Context, un string) (*domain.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.byUN[un]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}
func (f *fakeUserRepo) GetByEmail(context.Context, string) (*domain.User, error) {
	return nil, domain.ErrNotFound
}
func (f *fakeUserRepo) GetByGitHubID(context.Context, int64) (*domain.User, error) {
	return nil, domain.ErrNotFound
}
func (f *fakeUserRepo) Update(context.Context, *domain.User) error { return nil }
func (f *fakeUserRepo) SetStripeAccount(_ context.Context, id, acct string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.byID[id]; ok {
		u.StripeAccountID = acct
		return nil
	}
	return domain.ErrNotFound
}

type fakeProjectRepo struct {
	mu     sync.Mutex
	byID   map[string]*domain.Project
	seq    int
	ratings map[string]domain.ReviewStats
}

func newFakeProjectRepo() *fakeProjectRepo {
	return &fakeProjectRepo{byID: map[string]*domain.Project{}, ratings: map[string]domain.ReviewStats{}}
}

func (f *fakeProjectRepo) add(p *domain.Project) *domain.Project {
	f.seq++
	if p.ID == "" {
		p.ID = "p" + itoa(f.seq)
	}
	f.byID[p.ID] = p
	return p
}
func (f *fakeProjectRepo) Create(_ context.Context, p *domain.Project) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.add(p)
	return nil
}
func (f *fakeProjectRepo) GetByID(_ context.Context, id string) (*domain.Project, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if p, ok := f.byID[id]; ok {
		return p, nil
	}
	return nil, domain.ErrNotFound
}
func (f *fakeProjectRepo) GetBySlug(context.Context, string) (*domain.Project, error) {
	return nil, domain.ErrNotFound
}
func (f *fakeProjectRepo) List(context.Context, domain.ProjectFilter) ([]*domain.Project, int64, error) {
	return nil, 0, nil
}
func (f *fakeProjectRepo) ListByBuilder(_ context.Context, builderID string) ([]*domain.Project, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*domain.Project
	for _, p := range f.byID {
		if p.BuilderID == builderID {
			out = append(out, p)
		}
	}
	return out, nil
}
func (f *fakeProjectRepo) Update(_ context.Context, p *domain.Project) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.byID[p.ID] = p
	return nil
}
func (f *fakeProjectRepo) UpdateRatingStats(_ context.Context, id string, s domain.ReviewStats) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ratings[id] = s
	if p, ok := f.byID[id]; ok {
		p.ReviewCount = s.Count
		p.AverageRating = s.Average
	}
	return nil
}

type fakeReviewRepo struct {
	mu   sync.Mutex
	list []*domain.Review
	seq  int
}

func newFakeReviewRepo() *fakeReviewRepo { return &fakeReviewRepo{} }

func (f *fakeReviewRepo) Create(_ context.Context, r *domain.Review) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.seq++
	r.ID = "r" + itoa(f.seq)
	f.list = append(f.list, r)
	return nil
}
func (f *fakeReviewRepo) GetByProjectAndReviewer(_ context.Context, pid, rid string) (*domain.Review, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, r := range f.list {
		if r.ProjectID == pid && r.ReviewerID == rid {
			return r, nil
		}
	}
	return nil, domain.ErrNotFound
}
func (f *fakeReviewRepo) ListByProject(context.Context, string, int, int) ([]*domain.Review, int64, error) {
	return nil, 0, nil
}
func (f *fakeReviewRepo) StatsForProject(_ context.Context, pid string) (domain.ReviewStats, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var sum, n int
	for _, r := range f.list {
		if r.ProjectID == pid {
			sum += r.Rating
			n++
		}
	}
	if n == 0 {
		return domain.ReviewStats{}, nil
	}
	return domain.ReviewStats{Count: n, Average: float64(sum) / float64(n)}, nil
}

type fakeProblemRepo struct {
	mu   sync.Mutex
	byID map[string]*domain.Problem
}

func newFakeProblemRepo() *fakeProblemRepo {
	return &fakeProblemRepo{byID: map[string]*domain.Problem{}}
}

func (f *fakeProblemRepo) Create(_ context.Context, p *domain.Problem) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if p.ID == "" {
		p.ID = "prob" + itoa(len(f.byID)+1)
	}
	f.byID[p.ID] = p
	return nil
}
func (f *fakeProblemRepo) GetByID(_ context.Context, id string) (*domain.Problem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if p, ok := f.byID[id]; ok {
		cp := *p
		return &cp, nil
	}
	return nil, domain.ErrNotFound
}
func (f *fakeProblemRepo) GetBySlug(context.Context, string) (*domain.Problem, error) {
	return nil, domain.ErrNotFound
}
func (f *fakeProblemRepo) List(context.Context, domain.ProblemFilter) ([]*domain.Problem, int64, error) {
	return nil, 0, nil
}
func (f *fakeProblemRepo) Update(context.Context, *domain.Problem) error { return nil }
func (f *fakeProblemRepo) Claim(_ context.Context, id, builderID string) (*domain.Problem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.byID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if p.Status != domain.ProblemStatusOpen {
		return nil, domain.ErrProblemClaimed
	}
	p.Status = domain.ProblemStatusClaimed
	p.ClaimedBy = builderID
	cp := *p
	return &cp, nil
}
func (f *fakeProblemRepo) AddUpvote(_ context.Context, id, userID string) (*domain.Problem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.byID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if !p.HasUpvoted(userID) {
		p.UpvotedBy = append(p.UpvotedBy, userID)
		p.Upvotes++
	}
	cp := *p
	return &cp, nil
}

type fakeStripeRepo struct {
	snap *domain.StripeSnapshot
}

func (f *fakeStripeRepo) Save(_ context.Context, s *domain.StripeSnapshot) error {
	f.snap = s
	return nil
}
func (f *fakeStripeRepo) Latest(_ context.Context, _ string) (*domain.StripeSnapshot, error) {
	if f.snap == nil {
		return nil, domain.ErrNotFound
	}
	return f.snap, nil
}

type fakeScoreRepo struct {
	upserted *domain.BuilderScore
}

func (f *fakeScoreRepo) GetByBuilderID(context.Context, string) (*domain.BuilderScore, error) {
	if f.upserted == nil {
		return nil, domain.ErrNotFound
	}
	return f.upserted, nil
}
func (f *fakeScoreRepo) Upsert(_ context.Context, s *domain.BuilderScore) error {
	f.upserted = s
	return nil
}
func (f *fakeScoreRepo) TopBuilders(context.Context, int) ([]*domain.BuilderScore, error) {
	return nil, nil
}

// recordingEnqueuer captures enqueue calls so tests can assert on them.
type recordingEnqueuer struct {
	scoreCalls   []string
	emailCalls   []string
}

func (r *recordingEnqueuer) EnqueueScoreRecalc(_ context.Context, builderID string) error {
	r.scoreCalls = append(r.scoreCalls, builderID)
	return nil
}
func (r *recordingEnqueuer) EnqueueProblemClaimedEmail(_ context.Context, problemID string) error {
	r.emailCalls = append(r.emailCalls, problemID)
	return nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
