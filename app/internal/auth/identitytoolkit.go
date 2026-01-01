package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type IdentityToolkitClient struct {
	apiKey string
	http   *http.Client
}

func NewIdentityToolkitClient(apiKey string) *IdentityToolkitClient {
	return &IdentityToolkitClient{
		apiKey: apiKey,
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type idtkAuthResponse struct {
	LocalID      string `json:"localId"`
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	Email        string `json:"email"`
	Error        *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *IdentityToolkitClient) SignUp(ctx context.Context, email, password string) (uid, idToken, refreshToken string, err error) {
	return c.call(ctx, "accounts:signUp", email, password)
}

func (c *IdentityToolkitClient) SignInWithPassword(ctx context.Context, email, password string) (uid, idToken, refreshToken string, err error) {
	return c.call(ctx, "accounts:signInWithPassword", email, password)
}

func (c *IdentityToolkitClient) call(ctx context.Context, method, email, password string) (uid, idToken, refreshToken string, err error) {
	if c.apiKey == "" {
		return "", "", "", errors.New("FIREBASE_API_KEY is required for register/login")
	}

	body := map[string]any{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}
	b, _ := json.Marshal(body)

	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/%s?key=%s", method, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer res.Body.Close()

	var out idtkAuthResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", "", "", err
	}

	if res.StatusCode >= 400 {
		if out.Error != nil && out.Error.Message != "" {
			return "", "", "", fmt.Errorf("identitytoolkit: %s", out.Error.Message)
		}
		return "", "", "", fmt.Errorf("identitytoolkit: http %d", res.StatusCode)
	}

	if out.LocalID == "" || out.IDToken == "" || out.RefreshToken == "" {
		return "", "", "", errors.New("identitytoolkit: unexpected empty response")
	}
	return out.LocalID, out.IDToken, out.RefreshToken, nil
}


