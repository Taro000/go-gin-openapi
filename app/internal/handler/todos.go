package handler

import (
	"database/sql"
	"errors"
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
		apiStatus, ok := todoCodeToStatus(t.Status)
		if !ok {
			internalErr(c, errors.New("invalid todo status code in db"))
			return
		}
		status := apiStatus

		var due *schemas.TodoDueDatetime
		if t.DueDatetime != nil {
			s := formatTodoDueDatetime(*t.DueDatetime)
			due = &s
		}
		out = append(out, struct {
			DueDatetime *schemas.TodoDueDatetime `json:"due_datetime,omitempty"`
			Id          *string                `json:"id,omitempty"`
			Status      *schemas.TodoStatus    `json:"status,omitempty"`
			Title       *string                `json:"title,omitempty"`
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

	title := strings.TrimSpace(*req.Title)
	if runeLen(title) > 30 {
		badRequest(c, "title must be <= 30 chars")
		return
	}
	if runeLen(*req.Content) > 1000 {
		badRequest(c, "content must be <= 1000 chars")
		return
	}

	apiStatus := schemas.TodoStatus("未着手")
	if req.Status != nil && strings.TrimSpace(string(*req.Status)) != "" {
		apiStatus = *req.Status
	}
	statusCode, ok := todoStatusToCode(apiStatus)
	if !ok {
		badRequest(c, "status must be one of: 未着手, 進行中, 完了, 保留")
		return
	}

	var due *time.Time
	if req.DueDatetime != nil && strings.TrimSpace(string(*req.DueDatetime)) != "" {
		t, err := parseTodoDueDatetime(*req.DueDatetime)
		if err != nil {
			badRequest(c, err.Error())
			return
		}
		due = &t
	}

	if err := a.repos.Statuses.Ensure(c.Request.Context(), statusCode); err != nil {
		internalErr(c, err)
		return
	}

	id := uuid.NewString()
	t := repo.Todo{
		ID:          id,
		Owner:       string(userId),
		Status:      statusCode,
		Title:       title,
		Content:     *req.Content,
		DueDatetime: due,
	}
	if err := a.repos.Todos.Create(c.Request.Context(), t); err != nil {
		internalErr(c, err)
		return
	}
	c.JSON(201, schemas.CreateTodoResponse{Id: &id})
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

	apiStatus, ok := todoCodeToStatus(t.Status)
	if !ok {
		internalErr(c, errors.New("invalid todo status code in db"))
		return
	}
	var due *schemas.TodoDueDatetime
	if t.DueDatetime != nil {
		s := formatTodoDueDatetime(*t.DueDatetime)
		due = &s
	}

	c.JSON(200, schemas.GetTodoDetailResponse{
		Title:       &t.Title,
		Content:     &t.Content,
		Status:      &apiStatus,
		DueDatetime: due,
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

	var title *string
	if req.Title != nil {
		s := strings.TrimSpace(*req.Title)
		if s == "" {
			badRequest(c, "title must not be empty")
			return
		}
		if runeLen(s) > 30 {
			badRequest(c, "title must be <= 30 chars")
			return
		}
		title = &s
	}

	if req.Content != nil && runeLen(*req.Content) > 1000 {
		badRequest(c, "content must be <= 1000 chars")
		return
	}

	var status *string
	if req.Status != nil && strings.TrimSpace(string(*req.Status)) != "" {
		code, ok := todoStatusToCode(*req.Status)
		if !ok {
			badRequest(c, "status must be one of: 未着手, 進行中, 完了, 保留")
			return
		}
		if err := a.repos.Statuses.Ensure(c.Request.Context(), code); err != nil {
			internalErr(c, err)
			return
		}
		status = &code
	}

	var dueDatetime *time.Time
	if req.DueDatetime != nil && strings.TrimSpace(string(*req.DueDatetime)) != "" {
		t, err := parseTodoDueDatetime(*req.DueDatetime)
		if err != nil {
			badRequest(c, err.Error())
			return
		}
		dueDatetime = &t
	}

	if err := a.repos.Todos.UpdateByIDOwner(c.Request.Context(), string(todoId), string(userId), title, req.Content, status, dueDatetime); err != nil {
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
	c.Status(204)
}


