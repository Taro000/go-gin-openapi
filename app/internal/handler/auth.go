package handler

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"

	"go-gin-webapi/internal/repo"
	"go-gin-webapi/schemas"
)

func (a *API) PostRegister(c *gin.Context) {
	var req schemas.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid json")
		return
	}
	if req.Email == nil || strings.TrimSpace(string(*req.Email)) == "" || req.Password == nil || strings.TrimSpace(*req.Password) == "" || req.Nickname == nil || strings.TrimSpace(*req.Nickname) == "" {
		badRequest(c, "email/password/nickname are required")
		return
	}
	if runeLen(strings.TrimSpace(*req.Nickname)) > 20 {
		badRequest(c, "nickname must be <= 20 chars")
		return
	}

	uid, idToken, refreshToken, err := a.idtk.SignUp(c.Request.Context(), string(*req.Email), *req.Password)
	if err != nil {
		// spec: 400/500 only
		badRequest(c, err.Error())
		return
	}

	if err := a.repos.Users.Create(c.Request.Context(), repo.User{
		UID:      uid,
		Nickname: *req.Nickname,
		Email:    string(*req.Email),
	}); err != nil {
		if isMySQLDuplicate(err) {
			badRequest(c, "email already exists")
			return
		}
		internalErr(c, err)
		return
	}

	c.JSON(201, schemas.RegisterUserResponse{
		Uid:          strPtr(uid),
		AccessToken:  strPtr(idToken),
		RefreshToken: strPtr(refreshToken),
	})
}

func (a *API) PostLogin(c *gin.Context) {
	var req schemas.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid json")
		return
	}
	if req.Email == nil || strings.TrimSpace(string(*req.Email)) == "" || req.Password == nil || strings.TrimSpace(*req.Password) == "" {
		badRequest(c, "email/password are required")
		return
	}

	uid, idToken, refreshToken, err := a.idtk.SignInWithPassword(c.Request.Context(), string(*req.Email), *req.Password)
	if err != nil {
		// spec: 400/500 only
		badRequest(c, err.Error())
		return
	}

	if _, err := a.repos.Users.GetByUID(c.Request.Context(), uid); err != nil {
		if err == sql.ErrNoRows {
			badRequest(c, "user not found in db (register first)")
			return
		}
		internalErr(c, err)
		return
	}

	c.JSON(201, schemas.LoginUserResponse{
		Uid:          strPtr(uid),
		AccessToken:  strPtr(idToken),
		RefreshToken: strPtr(refreshToken),
	})
}

func (a *API) PostLogout(c *gin.Context) {
	var req schemas.LogoutUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid json")
		return
	}
	if req.UserId == nil || strings.TrimSpace(*req.UserId) == "" {
		badRequest(c, "user_id is required")
		return
	}
	if !a.requireSelf(c, *req.UserId) {
		return
	}

	if a.fbAdmin != nil && a.fbAdmin.Auth != nil {
		_ = a.fbAdmin.Auth.RevokeRefreshTokens(c.Request.Context(), *req.UserId)
	}

	msg := "logged out"
	c.JSON(201, schemas.LogoutUserResponse{Message: &msg})
}


