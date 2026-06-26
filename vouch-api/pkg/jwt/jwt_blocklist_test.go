package jwt_test

import (
	"context"
	"testing"
	"time"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/jwt"
)

type inMemBlocklist struct {
	blocked map[string]struct{}
}

func newInMemBlocklist() *inMemBlocklist {
	return &inMemBlocklist{blocked: make(map[string]struct{})}
}

func (b *inMemBlocklist) Block(_ context.Context, hash string, _ time.Duration) error {
	b.blocked[hash] = struct{}{}
	return nil
}

func (b *inMemBlocklist) IsBlocked(_ context.Context, hash string) (bool, error) {
	_, ok := b.blocked[hash]
	return ok, nil
}

func TestTokenRotation_OldRefreshBlockedAfterRotate(t *testing.T) {
	bl := newInMemBlocklist()
	mgr := jwt.NewManager(
		"access-secret-32-characters-pad1",
		"refresh-secret-32-characters-pad",
	).WithBlocklist(bl)

	pair, err := mgr.GenerateTyped("u1", "alice", "user", "user")
	if err != nil {
		t.Fatal(err)
	}

	// Revoke the refresh token
	if err := mgr.RevokeRefresh(context.Background(), pair.RefreshToken); err != nil {
		t.Fatal(err)
	}

	// Verification should now fail
	_, err = mgr.VerifyRefresh(pair.RefreshToken)
	if err == nil {
		t.Fatal("expected error for revoked refresh token, got nil")
	}
}

func TestAccessToken_NotUsableAsRefresh(t *testing.T) {
	mgr := jwt.NewManager(
		"access-secret-32-characters-pad1",
		"refresh-secret-32-characters-pad",
	)

	pair, err := mgr.GenerateTyped("u1", "alice", "user", "user")
	if err != nil {
		t.Fatal(err)
	}

	// Access token must not verify as a refresh token
	_, err = mgr.VerifyRefresh(pair.AccessToken)
	if err == nil {
		t.Fatal("expected error using access token as refresh token")
	}
}

func TestRefreshToken_NotUsableAsAccess(t *testing.T) {
	mgr := jwt.NewManager(
		"access-secret-32-characters-pad1",
		"refresh-secret-32-characters-pad",
	)

	pair, err := mgr.GenerateTyped("u1", "alice", "user", "user")
	if err != nil {
		t.Fatal(err)
	}

	// Refresh token must not verify as an access token (different secret)
	_, err = mgr.VerifyAccess(pair.RefreshToken)
	if err == nil {
		t.Fatal("expected error using refresh token as access token")
	}
}
