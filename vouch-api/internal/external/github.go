package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
)

// GitHubClient implements service.GitHubGateway against the real GitHub API.
type GitHubClient struct {
	clientID     string
	clientSecret string
	redirectURL  string
	http         *http.Client
}

// NewGitHubClient constructs a GitHubClient.
func NewGitHubClient(clientID, clientSecret, redirectURL string) *GitHubClient {
	return &GitHubClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		http:         &http.Client{Timeout: 10 * time.Second},
	}
}

// ExchangeCode swaps an OAuth code for an access token and fetches the profile.
func (g *GitHubClient) ExchangeCode(ctx context.Context, code string) (service.GitHubProfile, error) {
	token, err := g.exchange(ctx, code)
	if err != nil {
		return service.GitHubProfile{}, fmt.Errorf("GitHubClient.ExchangeCode: %w", err)
	}
	profile, err := g.fetchUser(ctx, token)
	if err != nil {
		return service.GitHubProfile{}, fmt.Errorf("GitHubClient.ExchangeCode: %w", err)
	}
	if profile.Email == "" {
		profile.Email, _ = g.fetchPrimaryEmail(ctx, token)
	}
	return profile, nil
}

func (g *GitHubClient) exchange(ctx context.Context, code string) (string, error) {
	form := url.Values{
		"client_id":     {g.clientID},
		"client_secret": {g.clientSecret},
		"code":          {code},
		"redirect_uri":  {g.redirectURL},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://github.com/login/oauth/access_token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := g.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error_description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.AccessToken == "" {
		return "", fmt.Errorf("github oauth failed: %s", body.Error)
	}
	return body.AccessToken, nil
}

func (g *GitHubClient) fetchUser(ctx context.Context, token string) (service.GitHubProfile, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := g.http.Do(req)
	if err != nil {
		return service.GitHubProfile{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return service.GitHubProfile{}, fmt.Errorf("github user api status %d", resp.StatusCode)
	}

	var u struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Bio       string `json:"bio"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return service.GitHubProfile{}, err
	}
	return service.GitHubProfile{
		ID: u.ID, Login: u.Login, Name: u.Name, Email: u.Email,
		AvatarURL: u.AvatarURL, Bio: u.Bio,
	}, nil
}

func (g *GitHubClient) fetchPrimaryEmail(ctx context.Context, token string) (string, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := g.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	return "", fmt.Errorf("no verified primary email")
}
