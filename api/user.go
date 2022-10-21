package api

import (
	"database/sql"
	"github.com/mhalavanja/go-rest-api/db/sqlc"
	"github.com/mhalavanja/go-rest-api/token"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) getUser(ctx *gin.Context) {
	userId := ctx.MustGet(authPayload).(*token.Payload).UserId
	usr, err := server.store.GetUser(ctx, userId)

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

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := sqlc.CreateUserParams{
		Username:       req.Username,
		HashedPassword: string(hashedPassword),
		Email:          req.Email,
	}

	usr, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, usr)
}

func (server *Server) deleteUser(ctx *gin.Context) {
	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	err := server.store.DeleteUser(ctx, userId)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}

type updateUsernameRequest struct {
	Username string `json:"username" binding:"required"`
}

func (server *Server) updateUsername(ctx *gin.Context) {
	var req updateUsernameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	arg := sqlc.UpdateUsernameParams{
		Username: req.Username,
		ID:       userId,
	}

	err := server.store.UpdateUsername(ctx, arg)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}

type updateEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (server *Server) updateEmail(ctx *gin.Context) {
	var req updateEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	arg := sqlc.UpdateEmailParams{
		Email: req.Email,
		ID:    userId,
	}

	err := server.store.UpdateEmail(ctx, arg)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}

type updatePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func (server *Server) updatePassword(ctx *gin.Context) {
	var req updatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := sqlc.UpdatePasswordParams{
		HashedPassword: string(hashedPassword),
		ID:             userId,
	}

	err = server.store.UpdatePassword(ctx, arg)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}
