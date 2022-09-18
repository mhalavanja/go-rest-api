package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type authUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type authUserResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	Username             string    `json:"username"`
}

const wrongUsernameOrPassword = "Wrong username or password"

func (server *Server) authUser(ctx *gin.Context) {
	var req authUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("authUserRequest - ctx.ShouldBindJSON req =", req, "err =", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		log.Println("authUserRequest - server.store.GetUserByUsername req =", req, "err =", err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, wrongUsernameOrPassword)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		log.Println("authUserRequest - bcrypt.CompareHashAndPassword - wrong password")
		ctx.JSON(http.StatusUnauthorized, wrongUsernameOrPassword)
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		log.Println("authUserRequest - server.tokenMaker.CreateToken - err =", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := authUserResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
		Username:             user.Username,
	}

	ctx.JSON(http.StatusOK, rsp)
}
