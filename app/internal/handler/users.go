package handler

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"go-gin-webapi/schemas"
)

func (a *API) GetUsersUserId(c *gin.Context, userId schemas.UserId) {
	if !a.requireSelf(c, string(userId)) {
		return
	}
	u, err := a.repos.Users.GetByUID(c.Request.Context(), string(userId))
	if err != nil {
		if err == sql.ErrNoRows {
			notFound(c)
			return
		}
		internalErr(c, err)
		return
	}

	email := openapi_types.Email(u.Email)
	c.JSON(200, schemas.GetUserDetailResponse{
		Nickname: &u.Nickname,
		Email:    &email,
	})
}

func (a *API) PutUsersUserId(c *gin.Context, userId schemas.UserId) {
	if !a.requireSelf(c, string(userId)) {
		return
	}
	var req schemas.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid json")
		return
	}

	var emailStr *string
	if req.Email != nil {
		s := string(*req.Email)
		emailStr = &s
	}

	u, err := a.repos.Users.Update(c.Request.Context(), string(userId), req.Nickname, emailStr)
	if err != nil {
		if err == sql.ErrNoRows {
			notFound(c)
			return
		}
		internalErr(c, err)
		return
	}

	email := openapi_types.Email(u.Email)
	c.JSON(200, schemas.UpdateUserResponse{
		Nickname: &u.Nickname,
		Email:    &email,
	})
}


