package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhalavanja/go-rest-api/db/sqlc"
	"github.com/mhalavanja/go-rest-api/token"
)

func (server *Server) getFriends(ctx *gin.Context) {
	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	friendNames, err := server.store.GetFriends(ctx, userId)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, friendNames)
}

func (server *Server) getFriend(ctx *gin.Context) {
	var friendId int64
	if err := ctx.ShouldBindUri(&friendId); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	arg := sqlc.GetFriendParams{
		UserID:       userId,
		UserIDFriend: friendId,
	}

	friend, err := server.store.GetFriend(ctx, arg)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sqlc.User{
		ID:       friendId,
		Username: friend.Username,
		Email:    friend.Email,
	})
}

func (server *Server) addFriend(ctx *gin.Context) {
	var friendId int64
	if err := ctx.ShouldBindJSON(&friendId); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	arg := sqlc.AddFriendParams{
		UserIDFriend: friendId,
		UserID:       userId,
	}

	err := server.store.AddFriend(ctx, arg)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusCreated)
}

func (server *Server) deleteFriend(ctx *gin.Context) {
	var friendId int64
	if err := ctx.ShouldBindUri(&friendId); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	arg := sqlc.DeleteFriendParams{
		UserIDFriend: friendId,
		UserID:       userId,
	}

	err := server.store.DeleteFriend(ctx, arg)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
