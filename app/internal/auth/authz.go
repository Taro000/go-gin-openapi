package auth

import (
	"context"
	"errors"
	"strings"

	fbauth "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"

	"go-gin-webapi/internal/config"
)

type Verifier struct {
	admin    *FirebaseAdmin
	adminErr error
	cfg      config.AuthConfig
}

func NewVerifier(admin *FirebaseAdmin, adminErr error, cfg config.AuthConfig) *Verifier {
	return &Verifier{admin: admin, adminErr: adminErr, cfg: cfg}
}

var ErrUnauthorized = errors.New("unauthorized")

func (v *Verifier) RequireUID(c *gin.Context) (string, error) {
	if v.cfg.Bypass {
		if uid := strings.TrimSpace(c.GetHeader("X-User-Id")); uid != "" {
			return uid, nil
		}
	}

	h := c.GetHeader("Authorization")
	if h == "" {
		return "", ErrUnauthorized
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrUnauthorized
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrUnauthorized
	}
	if v.admin == nil || v.admin.Auth == nil {
		if v.adminErr != nil {
			return "", v.adminErr
		}
		return "", errors.New("firebase admin not configured")
	}
	t, err := v.admin.Auth.VerifyIDToken(context.Background(), token)
	if err != nil {
		// Known "token is not acceptable" cases -> 401.
		// Anything else likely indicates server-side misconfiguration (project id/credentials/network),
		// so return the underlying error to help debugging (handler will map it to 500).
		if fbauth.IsIDTokenInvalid(err) ||
			fbauth.IsIDTokenExpired(err) ||
			fbauth.IsIDTokenRevoked(err) ||
			fbauth.IsUserDisabled(err) {
			return "", ErrUnauthorized
		}
		return "", err
	}
	if t == nil || t.UID == "" {
		return "", ErrUnauthorized
	}
	return t.UID, nil
}


