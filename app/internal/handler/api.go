package handler

import (
	"go-gin-webapi/internal/auth"
	"go-gin-webapi/internal/config"
	"go-gin-webapi/internal/repo"
	"go-gin-webapi/schemas"
)

type API struct {
	repos    *repo.Repos
	idtk     *auth.IdentityToolkitClient
	fbAdmin  *auth.FirebaseAdmin
	verifier *auth.Verifier
}

func NewAPI(repos *repo.Repos, idtk *auth.IdentityToolkitClient, fbAdmin *auth.FirebaseAdmin, fbAdminErr error, authCfg config.AuthConfig) *API {
	return &API{
		repos:    repos,
		idtk:     idtk,
		fbAdmin:  fbAdmin,
		verifier: auth.NewVerifier(fbAdmin, fbAdminErr, authCfg),
	}
}

var _ schemas.ServerInterface = (*API)(nil)

