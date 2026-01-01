package auth

import (
	"context"
	"errors"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"go-gin-webapi/internal/config"
)

type FirebaseAdmin struct {
	Auth *auth.Client
}

func NewFirebaseAdmin(ctx context.Context, cfg config.FirebaseConfig) (*FirebaseAdmin, error) {
	// Allow emulator if configured.
	if cfg.AuthEmulatorHostport != "" {
		_ = os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", cfg.AuthEmulatorHostport)
	}

	opts := []option.ClientOption{}
	if cfg.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialsFile))
	} else if cfg.ServiceAccountJSON != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(cfg.ServiceAccountJSON)))
	} else {
		// If no explicit credentials are provided, firebase may still work via default credentials
		// (e.g., running on GCP). In local devcontainer this is usually missing.
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: cfg.ProjectID}, opts...)
	if err != nil {
		return nil, err
	}
	ac, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	if ac == nil {
		return nil, errors.New("firebase auth client is nil")
	}
	return &FirebaseAdmin{Auth: ac}, nil
}


