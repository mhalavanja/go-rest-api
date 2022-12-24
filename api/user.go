package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/lib/pq"
	"github.com/mhalavanja/go-rest-api/consts"
	"github.com/mhalavanja/go-rest-api/db/sqlc"
	"github.com/mhalavanja/go-rest-api/token"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) getUser(ctx *gin.Context) {
	userId := ctx.MustGet(authPayload).(*token.Payload).UserId
	usr, err := server.store.GetUser(ctx, userId)

	if err != nil {
		log.Println("ERROR: getUser userId=", userId, " err=", err.Error())
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, consts.UserDoesNotExist)
			return
		}

		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
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
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.Provide+"email, username and password")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	arg := sqlc.CreateUserParams{
		Username:       req.Username,
		HashedPassword: string(hashedPassword),
		Email:          req.Email,
	}

	usr, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		pqErr := err.(*pq.Error)
		if string(pqErr.Code) == "23505" {
			ctx.JSON(http.StatusConflict, consts.UserAlreadyExists)
			return
		}

		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.JSON(http.StatusCreated, usr)
}

func (server *Server) deleteUser(ctx *gin.Context) {
	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	err := server.store.DeleteUser(ctx, userId)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.Status(http.StatusOK)
}

type updateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.Provide+"username, email and password")
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId
	var bytes []byte = nil
	if req.Password != "" {
		var err error
		bytes, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
			return
		}
	}

	var hashedPassword string = ""
	if bytes != nil {
		hashedPassword = string(bytes)
	}

	arg := sqlc.UpdateUserParams{
		ID:             userId,
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.Status(http.StatusOK)
}
