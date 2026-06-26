package service

import (
	"context"
	"errors"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/jwt"
)

func newTestJWT() *jwt.Manager {
	return jwt.NewManager(
		"test-secret-32-chars-long-xxxxxxxx",
		"refresh-secret-32-chars-long-xxxx",
	)
}

// fakeCompanyRepo is an in-memory CompanyRepository for unit tests.
type fakeCompanyRepo struct {
	byEmail map[string]*domain.Company
	bySlug  map[string]*domain.Company
	byID    map[string]*domain.Company
}

func newFakeCompanyRepo() *fakeCompanyRepo {
	return &fakeCompanyRepo{
		byEmail: make(map[string]*domain.Company),
		bySlug:  make(map[string]*domain.Company),
		byID:    make(map[string]*domain.Company),
	}
}

func (r *fakeCompanyRepo) Create(_ context.Context, c *domain.Company) error {
	if _, ok := r.byEmail[c.Email]; ok {
		return domain.ErrConflict
	}
	c.ID = "cid-" + c.Email
	r.byEmail[c.Email] = c
	r.bySlug[c.Slug] = c
	r.byID[c.ID] = c
	return nil
}
func (r *fakeCompanyRepo) GetByID(_ context.Context, id string) (*domain.Company, error) {
	if c, ok := r.byID[id]; ok {
		return c, nil
	}
	return nil, domain.ErrNotFound
}
func (r *fakeCompanyRepo) GetByEmail(_ context.Context, email string) (*domain.Company, error) {
	if c, ok := r.byEmail[email]; ok {
		return c, nil
	}
	return nil, domain.ErrNotFound
}
func (r *fakeCompanyRepo) GetBySlug(_ context.Context, slug string) (*domain.Company, error) {
	if c, ok := r.bySlug[slug]; ok {
		return c, nil
	}
	return nil, domain.ErrNotFound
}
func (r *fakeCompanyRepo) Update(_ context.Context, c *domain.Company) error {
	if _, ok := r.byID[c.ID]; !ok {
		return domain.ErrNotFound
	}
	r.byID[c.ID] = c
	return nil
}
func (r *fakeCompanyRepo) List(_ context.Context, _ int, _ int) ([]*domain.Company, int64, error) {
	var out []*domain.Company
	for _, c := range r.byID {
		out = append(out, c)
	}
	return out, int64(len(out)), nil
}

func TestCompanyRegister_Success(t *testing.T) {
	svc := NewCompanyService(newFakeCompanyRepo(), newTestJWT())
	res, err := svc.Register(ctx(), "Acme Inc", "cto@acme.com", "password123", "https://acme.com", "11-50")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if res.Company.Email != "cto@acme.com" {
		t.Fatalf("unexpected email: %s", res.Company.Email)
	}
	if res.Tokens.AccessToken == "" {
		t.Fatal("expected access token")
	}
}

func TestCompanyRegister_DuplicateEmail(t *testing.T) {
	repo := newFakeCompanyRepo()
	svc := NewCompanyService(repo, newTestJWT())
	if _, err := svc.Register(ctx(), "Acme", "cto@acme.com", "password123", "", ""); err != nil {
		t.Fatalf("first register: %v", err)
	}
	_, err := svc.Register(ctx(), "Acme2", "cto@acme.com", "password456", "", "")
	if err == nil {
		t.Fatal("expected duplicate email error")
	}
}

func TestCompanyLogin_WrongPassword(t *testing.T) {
	svc := NewCompanyService(newFakeCompanyRepo(), newTestJWT())
	if _, err := svc.Register(ctx(), "Acme", "cto@acme.com", "correctpassword", "", ""); err != nil {
		t.Fatalf("register: %v", err)
	}
	_, err := svc.Login(ctx(), "cto@acme.com", "wrongpassword")
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestCompanyLogin_UnknownEmail(t *testing.T) {
	svc := NewCompanyService(newFakeCompanyRepo(), newTestJWT())
	_, err := svc.Login(ctx(), "ghost@acme.com", "password")
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestCompanyRegister_WeakPassword(t *testing.T) {
	svc := NewCompanyService(newFakeCompanyRepo(), newTestJWT())
	_, err := svc.Register(ctx(), "Acme", "cto@acme.com", "short", "", "")
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for short password, got %v", err)
	}
}

func TestCompanyRegister_InvalidWebsite(t *testing.T) {
	svc := NewCompanyService(newFakeCompanyRepo(), newTestJWT())
	_, err := svc.Register(ctx(), "Acme", "cto@acme.com", "password123", "javascript:alert(1)", "")
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for javascript: URL, got %v", err)
	}
}

func TestCompanyLogin_Success(t *testing.T) {
	svc := NewCompanyService(newFakeCompanyRepo(), newTestJWT())
	if _, err := svc.Register(ctx(), "Acme", "cto@acme.com", "password123", "", ""); err != nil {
		t.Fatalf("register: %v", err)
	}
	res, err := svc.Login(ctx(), "cto@acme.com", "password123")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if res.Company.Email != "cto@acme.com" {
		t.Fatalf("unexpected email: %s", res.Company.Email)
	}
}
