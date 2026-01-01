package handler

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"

	"go-gin-webapi/schemas"
)

func (a *API) PostUsersUserIdTodosTodoIdGoodlucks(c *gin.Context, userId schemas.UserId, todoId schemas.TodoId) {
	if !a.requireSelf(c, string(userId)) {
		return
	}
	var req schemas.CreateGoodluckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid json")
		return
	}
	if req.UserId != nil && strings.TrimSpace(*req.UserId) != "" && *req.UserId != string(userId) {
		badRequest(c, "user_id mismatch")
		return
	}
	if req.TodoId != nil && strings.TrimSpace(*req.TodoId) != "" && *req.TodoId != string(todoId) {
		badRequest(c, "todo_id mismatch")
		return
	}
	if err := a.repos.Goodlucks.Create(c.Request.Context(), string(userId), string(todoId)); err != nil {
		if isMySQLFKViolation(err) {
			notFound(c)
			return
		}
		internalErr(c, err)
		return
	}
	msg := "goodluck created"
	c.JSON(200, schemas.CreateGoodluckResponse{Message: &msg})
}

func (a *API) DeleteUsersUserIdTodosTodoIdGoodlucks(c *gin.Context, userId schemas.UserId, todoId schemas.TodoId) {
	if !a.requireSelf(c, string(userId)) {
		return
	}
	if err := a.repos.Goodlucks.Delete(c.Request.Context(), string(userId), string(todoId)); err != nil {
		if err == sql.ErrNoRows {
			notFound(c)
			return
		}
		internalErr(c, err)
		return
	}
	msg := "goodluck deleted"
	c.JSON(200, schemas.DeleteGoodluckResponse{Message: &msg})
}


