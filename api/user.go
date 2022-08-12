package api

import (
	"database/sql"
	"log"
	"net/http"

	"dipl/db/sqlc"
	"github.com/gin-gonic/gin"
)

type GetDeleteUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1" `
}

func (server *Server) getUser(ctx *gin.Context) {
	var req GetDeleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	usr, err := server.store.GetUser(ctx, req.ID)

	if err != nil {
		log.Print(err.Error())
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, usr)
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required" `
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	usr, err := server.store.CreateUser(ctx, sqlc.CreateUserParams(req))
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, usr)
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req GetDeleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteUser(ctx, req.ID)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}
