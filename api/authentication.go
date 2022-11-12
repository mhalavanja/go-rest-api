package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mhalavanja/go-rest-api/consts"
	"golang.org/x/crypto/bcrypt"
)

type authUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type authUserResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	UserID               int64     `json:"user_id"`
}

func (server *Server) authUser(ctx *gin.Context) {
	var req authUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("ERROR: authUserRequest - ctx.ShouldBindJSON req =", req, "err =", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.Provide+"username and password")
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		log.Println("ERROR: authUserRequest - server.store.GetUserByUsername req =", req, "err =", err.Error())
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, consts.WrongUsernameOrPassword)
			return
		}
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		log.Println("ERROR: authUserRequest - bcrypt.CompareHashAndPassword - wrong password")
		ctx.JSON(http.StatusUnauthorized, consts.WrongUsernameOrPassword)
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		log.Println("ERROR: authUserRequest - server.tokenMaker.CreateToken - err =", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	rsp := authUserResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
		UserID:               user.ID,
	}

	ctx.JSON(http.StatusOK, rsp)
}
