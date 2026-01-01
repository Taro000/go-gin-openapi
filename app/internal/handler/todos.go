package handler

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-gin-webapi/internal/repo"
	"go-gin-webapi/schemas"
)

func (a *API) GetUsersUserIdTodos(c *gin.Context, userId schemas.UserId) {
	if !a.requireSelf(c, string(userId)) {
		return
	}

	todos, err := a.repos.Todos.ListByOwner(c.Request.Context(), string(userId))
	if err != nil {
		internalErr(c, err)
		return
	}

	out := make(schemas.GetTodoListResponse, 0, len(todos))
	for _, t := range todos {
		id := t.ID
		title := t.Title
		status := t.Status
		var due *time.Time
		if t.DueDatetime != nil {
			due = t.DueDatetime
		}
		out = append(out, struct {
			DueDatetime *time.Time `json:"due_datetime,omitempty"`
			Id          *string    `json:"id,omitempty"`
			Status      *string    `json:"status,omitempty"`
			Title       *string    `json:"title,omitempty"`
		}{
			DueDatetime: due,
			Id:          &id,
			Status:      &status,
			Title:       &title,
		})
	}
	c.JSON(200, out)
}

func (a *API) PostUsersUserIdTodos(c *gin.Context, userId schemas.UserId) {
	if !a.requireSelf(c, string(userId)) {
		return
	}
	var req schemas.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid json")
		return
	}
	if req.Title == nil || strings.TrimSpace(*req.Title) == "" || req.Content == nil {
		badRequest(c, "title/content are required")
		return
	}
	status := "00"
	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		status = *req.Status
	}
	if err := a.repos.Statuses.Ensure(c.Request.Context(), status); err != nil {
		internalErr(c, err)
		return
	}

	id := uuid.NewString()
	t := repo.Todo{
		ID:          id,
		Owner:       string(userId),
		Status:      status,
		Title:       *req.Title,
		Content:     *req.Content,
		DueDatetime: req.DueDatetime,
	}
	if err := a.repos.Todos.Create(c.Request.Context(), t); err != nil {
		internalErr(c, err)
		return
	}
	c.JSON(200, schemas.CreateTodoResponse{Id: &id})
}

func (a *API) GetUsersUserIdTodosTodoId(c *gin.Context, userId schemas.UserId, todoId schemas.TodoId) {
	if !a.requireSelf(c, string(userId)) {
		return
	}
	t, err := a.repos.Todos.GetByIDOwner(c.Request.Context(), string(todoId), string(userId))
	if err != nil {
		if err == sql.ErrNoRows {
			notFound(c)
			return
		}
		internalErr(c, err)
		return
	}
	c.JSON(200, schemas.GetTodoDetailResponse{
		Title:       &t.Title,
		Content:     &t.Content,
		Status:      &t.Status,
		DueDatetime: t.DueDatetime,
	})
}

func (a *API) PutUsersUserIdTodosTodoId(c *gin.Context, userId schemas.UserId, todoId schemas.TodoId) {
	if !a.requireSelf(c, string(userId)) {
		return
	}
	var req schemas.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "invalid json")
		return
	}
	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		if err := a.repos.Statuses.Ensure(c.Request.Context(), *req.Status); err != nil {
			internalErr(c, err)
			return
		}
	}
	if err := a.repos.Todos.UpdateByIDOwner(c.Request.Context(), string(todoId), string(userId), req.Title, req.Content, req.Status, req.DueDatetime); err != nil {
		if err == sql.ErrNoRows {
			notFound(c)
			return
		}
		internalErr(c, err)
		return
	}
	id := string(todoId)
	c.JSON(200, schemas.UpdateTodoResponse{Id: &id})
}

func (a *API) DeleteUsersUserIdTodosTodoId(c *gin.Context, userId schemas.UserId, todoId schemas.TodoId) {
	if !a.requireSelf(c, string(userId)) {
		return
	}
	if err := a.repos.Todos.DeleteByIDOwner(c.Request.Context(), string(todoId), string(userId)); err != nil {
		if err == sql.ErrNoRows {
			notFound(c)
			return
		}
		internalErr(c, err)
		return
	}
	msg := "deleted"
	c.JSON(200, schemas.DeleteTodoResponse{Message: &msg})
}


