package jwt

import (
	"context"
	"testing"
	"time"
)

// inMemoryBlocklist is a test-only blocklist that doesn't need Redis.
type inMemoryBlocklist struct {
	blocked map[string]time.Time
}

func newInMemoryBlocklist() *inMemoryBlocklist {
	return &inMemoryBlocklist{blocked: make(map[string]time.Time)}
}

func (b *inMemoryBlocklist) Block(_ context.Context, key string, exp time.Duration) error {
	b.blocked[key] = time.Now().Add(exp)
	return nil
}

func (b *inMemoryBlocklist) IsBlocked(_ context.Context, key string) (bool, error) {
	exp, ok := b.blocked[key]
	if !ok {
		return false, nil
	}
	if time.Now().After(exp) {
		delete(b.blocked, key)
		return false, nil
	}
	return true, nil
}

func newTestManager() *Manager {
	return NewManager("test-secret-32-chars-long-xxxxxxxx", "refresh-secret-32-chars-long-xxxx")
}

func TestGenerate_RoundTrip(t *testing.T) {
	m := newTestManager()
	pair, err := m.Generate("uid1", "alice", "builder")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	claims, err := m.VerifyAccess(pair.AccessToken)
	if err != nil {
		t.Fatalf("VerifyAccess: %v", err)
	}
	if claims.UserID != "uid1" || claims.Username != "alice" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestVerifyRefresh_WrongSecret(t *testing.T) {
	m := newTestManager()
	other := NewManager("other-secret-32-chars-long-xxxxxx", "other-refresh-secret-32-chars-xxx")
	pair, _ := m.Generate("uid1", "alice", "builder")

	_, err := other.VerifyRefresh(pair.RefreshToken)
	if err == nil {
		t.Fatal("expected error verifying refresh token with wrong secret")
	}
}

func TestRevokeRefresh_BlocksReuse(t *testing.T) {
	bl := newInMemoryBlocklist()
	m := newTestManager().WithBlocklist(bl)

	pair, err := m.Generate("uid1", "alice", "builder")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	// Verify works before revocation.
	if _, err := m.VerifyRefresh(pair.RefreshToken); err != nil {
		t.Fatalf("VerifyRefresh before revocation: %v", err)
	}

	// Revoke the token.
	if err := m.RevokeRefresh(context.Background(), pair.RefreshToken); err != nil {
		t.Fatalf("RevokeRefresh: %v", err)
	}

	// Verify must fail after revocation.
	if _, err := m.VerifyRefresh(pair.RefreshToken); err == nil {
		t.Fatal("expected revoked token to be rejected")
	}
}

func TestSubjectType_Company(t *testing.T) {
	m := newTestManager()
	pair, err := m.GenerateTyped("cid1", "acme", "company", SubjectTypeCompany)
	if err != nil {
		t.Fatalf("GenerateTyped: %v", err)
	}
	claims, err := m.VerifyAccess(pair.AccessToken)
	if err != nil {
		t.Fatalf("VerifyAccess: %v", err)
	}
	if claims.SubjectType != SubjectTypeCompany {
		t.Fatalf("expected SubjectTypeCompany, got %q", claims.SubjectType)
	}
}

func TestAccessToken_NotUsableAsRefresh(t *testing.T) {
	m := newTestManager()
	pair, _ := m.Generate("uid1", "alice", "builder")

	_, err := m.VerifyRefresh(pair.AccessToken)
	if err == nil {
		t.Fatal("access token must not validate as refresh token (different secret)")
	}
}
