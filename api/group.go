package api

import (
	"database/sql"
	"github.com/mhalavanja/go-rest-api/db/sqlc"
	"github.com/mhalavanja/go-rest-api/token"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type groupResponse struct {
	GroupId   int64  `json:"groupId" binding:"required"`
	GroupName string `json:"groupName" binding:"required"`
	OwnerId   int64  `json:"ownerId" binding:"required"`
}

type createGroupRequest struct {
	GroupName string `json:"groupName" binding:"required"`
}

func (server *Server) createGroup(ctx *gin.Context) {
	var req createGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	arg := sqlc.CreateGroupParams{
		Name:        req.GroupName,
		UserIDOwner: userId,
	}

	group, err := server.store.CreateGroup(ctx, arg)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, group)
}

type groupIdRequest struct {
	GroupId int64 `json:"groupId" binding:"required"`
}

func (server *Server) deleteGroup(ctx *gin.Context) {
	var req groupIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

type groupIdUserIdRequest struct {
	GroupId int64 `json:"groupId" binding:"required"`
	UserId  int64 `json:"userId" binding:"required"`
}

func (server *Server) updateGroupOwner(ctx *gin.Context) {
	var req groupIdUserIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

type updateGroupNameRequest struct {
	GroupId   int64  `json:"groupId" binding:"required"`
	GroupName string `json:"groupName" binding:"required"`
}

func (server *Server) updateGroupName(ctx *gin.Context) {
	var req updateGroupNameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

func (server *Server) joinGroup(ctx *gin.Context) {
	var req groupIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

func (server *Server) leaveGroup(ctx *gin.Context) {
	var req groupIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

func (server *Server) addUserToGroup(ctx *gin.Context) {
	var req groupIdUserIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

func (server *Server) removeUserFromGroup(ctx *gin.Context) {
	var req groupIdUserIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

func (server *Server) addUserAsAdmin(ctx *gin.Context) {
	var req groupIdUserIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

func (server *Server) removeUserAsAdmin(ctx *gin.Context) {
	var req groupIdUserIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}
