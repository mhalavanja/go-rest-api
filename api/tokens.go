package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mhalavanja/go-rest-api/consts"
	"github.com/mhalavanja/go-rest-api/db/sqlc"
	"github.com/mhalavanja/go-rest-api/token"
	"golang.org/x/crypto/bcrypt"
)

type authUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type authUserResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	UserID                int64     `json:"user_id"`
}

func (server *Server) authUser(ctx *gin.Context) {
	var req authUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("ERROR: ", req, "err =", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.Provide+"username and password")
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		log.Println("ERROR: ", req, "err =", err.Error())
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, consts.WrongUsernameOrPassword)
			return
		}
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, consts.WrongUsernameOrPassword)
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	session, err := server.store.GetSessionByUserId(ctx, user.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}
	if err == sql.ErrNoRows {
		refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.RefreshTokenDuration)
		if err != nil {
			log.Println("ERROR: ", err.Error())
			ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
			return
		}

		session, err = server.store.CreateSession(ctx, sqlc.CreateSessionParams{
			ID:           refreshPayload.ID,
			UserID:       user.ID,
			RefreshToken: refreshToken,
			UserAgent:    ctx.Request.UserAgent(),
			ClientIp:     ctx.ClientIP(),
			IsBlocked:    false,
			ExpiresAt:    refreshPayload.ExpiredAt,
		})
		if err != nil {
			log.Println("ERROR: ", err.Error())
			ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
			return
		}
	}

	rsp := authUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          session.RefreshToken,
		RefreshTokenExpiresAt: session.ExpiresAt,
		UserID:                user.ID,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type renewAccessTokenResponse struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	refreshToken, refreshPayload, shouldReturn := server.getRefreshToken(ctx)
	if shouldReturn {
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: ", err)
			ctx.JSON(http.StatusUnauthorized, consts.WrongUsernameOrPassword)
			return
		}
		log.Println("ERROR: ", err)
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	if session.IsBlocked {
		log.Println("Session is blocked")
		ctx.JSON(http.StatusUnauthorized, consts.BlockedSession)
		return
	}

	if session.UserID != refreshPayload.UserId {
		log.Println("session.UserID != refreshPayload.UserId")
		ctx.JSON(http.StatusUnauthorized, consts.IncorrectSession)
		return
	}

	if session.RefreshToken != refreshToken {
		log.Println("session.RefreshToken != refreshToken")
		ctx.JSON(http.StatusUnauthorized, consts.MismatchedSessionToken)
		return
	}

	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, consts.ExpiredSession)
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.UserId,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		log.Println("ERROR: ", err)
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getRefreshToken(ctx *gin.Context) (string, *token.Payload, bool) {
	refreshTokenHeader := ctx.GetHeader("refresh_token")

	if len(refreshTokenHeader) == 0 {
		err := errors.New("refresh_token header is not provided")
		log.Println("ERROR ", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return "", nil, true
	}
	fields := strings.Fields(refreshTokenHeader)

	refreshToken := fields[0]
	refreshPayload, err := server.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		log.Println("ERROR ", err)
		ctx.JSON(http.StatusUnauthorized, consts.WrongUsernameOrPassword)
		return "", nil, true
	}
	return refreshToken, refreshPayload, false
}

func (server *Server) deleteRefreshToken(ctx *gin.Context) {
	_ = ctx.MustGet(authPayload).(*token.Payload).UserId
	_, refreshPayload, shouldReturn := server.getRefreshToken(ctx)
	if shouldReturn {
		return
	}

	err := server.store.DeleteSession(ctx, refreshPayload.ID)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.Status(http.StatusOK)
}
