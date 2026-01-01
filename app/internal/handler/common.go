package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	mysqlDriver "github.com/go-sql-driver/mysql"

	"go-gin-webapi/internal/auth"
	"go-gin-webapi/schemas"
)

func strPtr(s string) *string { return &s }

func badRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, schemas.BadRequestJSONResponse{Error: &msg})
}

func unauthorized(c *gin.Context) {
	msg := "unauthorized"
	c.JSON(http.StatusUnauthorized, schemas.UnauthorizedJSONResponse{Error: &msg})
}

func forbidden(c *gin.Context) {
	msg := "forbidden"
	c.JSON(http.StatusForbidden, schemas.ForbiddenJSONResponse{Error: &msg})
}

func notFound(c *gin.Context) {
	msg := "not found"
	c.JSON(http.StatusNotFound, schemas.NotFoundJSONResponse{Error: &msg})
}

func internalErr(c *gin.Context, err error) {
	msg := err.Error()
	c.JSON(http.StatusInternalServerError, schemas.InternalServerErrorJSONResponse{Error: &msg})
}

func (a *API) requireSelf(c *gin.Context, userID string) bool {
	uid, err := a.verifier.RequireUID(c)
	if err != nil {
		// Token missing/invalid -> 401.
		// Verifier misconfigured (e.g. Firebase Admin not configured) -> 500 to aid debugging.
		if errors.Is(err, auth.ErrUnauthorized) {
			unauthorized(c)
		} else {
			internalErr(c, err)
		}
		return false
	}
	if uid != userID {
		forbidden(c)
		return false
	}
	return true
}

func isMySQLDuplicate(err error) bool {
	var me *mysqlDriver.MySQLError
	return errors.As(err, &me) && me.Number == 1062
}

func isMySQLFKViolation(err error) bool {
	var me *mysqlDriver.MySQLError
	return errors.As(err, &me) && me.Number == 1452
}


