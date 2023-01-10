package api

import (
	"log"
	"net/http"

	"github.com/lib/pq"
	"github.com/mhalavanja/go-rest-api/consts"
	"github.com/mhalavanja/go-rest-api/db/sqlc"
	"github.com/mhalavanja/go-rest-api/token"

	"github.com/gin-gonic/gin"
)

func (server *Server) getGroups(ctx *gin.Context) {
	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	groupNames, err := server.store.GetGroups(ctx, userId)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.JSON(http.StatusOK, groupNames)
}

type ID struct {
	Id int64 `uri:"id" binding:"required"`
}

func (server *Server) getGroup(ctx *gin.Context) {
	var id ID
	if err := ctx.ShouldBindUri(&id); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.ProvideGroupId)
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId
	arg := sqlc.GetGroupParams{
		ID:          id.Id,
		UserIDOwner: userId,
	}

	group, err := server.store.GetGroup(ctx, arg)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.JSON(http.StatusOK, group)
}

func (server *Server) createGroup(ctx *gin.Context) {
	var name string
	if err := ctx.ShouldBindJSON(&name); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.Provide+"group name")
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	arg := sqlc.CreateGroupParams{
		GroupName: name,
		UserID:    userId,
	}

	id, err := server.store.CreateGroup(ctx, arg)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.JSON(http.StatusCreated, id)
}

func (server *Server) deleteGroup(ctx *gin.Context) {
	var groupId ID
	if err := ctx.ShouldBindUri(&groupId); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.ProvideGroupId)
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	arg := sqlc.TryDeleteGroupParams{
		GroupID: groupId.Id,
		UserID:  userId,
	}

	err := server.store.TryDeleteGroup(ctx, arg)
	if err != nil {
		pqErr := err.(*pq.Error)
		if string(pqErr.Code) == "NOOWN" {
			ctx.JSON(http.StatusUnauthorized, consts.NotOwner)
			return
		}

		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.Status(http.StatusOK)
}

type groupIdUserIdRequest struct {
	GroupId int64 `json:"groupId" binding:"required"`
	UserId  int64 `json:"userId" binding:"required"`
}

// func (server *Server) updateGroupOwner(ctx *gin.Context) {
// 	var req groupIdUserIdRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		log.Print(err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// }

// type updateGroupNameRequest struct {
// 	GroupId   int64  `json:"groupId" binding:"required"`
// 	GroupName string `json:"groupName" binding:"required"`
// }

// func (server *Server) updateGroupName(ctx *gin.Context) {
// 	var req updateGroupNameRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		log.Print(err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// }

func (server *Server) leaveGroup(ctx *gin.Context) {
	var groupId int64
	if err := ctx.ShouldBindUri(&groupId); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.ProvideGroupId)
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	arg := sqlc.LeaveGroupParams{
		GroupID: groupId,
		UserID:  userId,
	}

	err := server.store.LeaveGroup(ctx, arg)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) addFriendToGroup(ctx *gin.Context) {
	var req groupIdUserIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.ProvideGroupIdAndUserId)
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId
	if userId == req.UserId {
		ctx.JSON(http.StatusUnprocessableEntity, "You can not add yourself to group")
		return
	}
	arg := sqlc.AddFriendToGroupParams{
		GroupID: req.GroupId,
		UserID:  req.UserId,
	}

	err := server.store.AddFriendToGroup(ctx, arg)
	if err != nil {
		pqErr := err.(*pq.Error)
		if string(pqErr.Code) == "NOFRN" {
			ctx.JSON(http.StatusUnauthorized, consts.NotFriends)
			return
		}

		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) removeUserFromGroup(ctx *gin.Context) {
	var req groupIdUserIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.ProvideGroupIdAndUserId)
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId
	if userId == req.UserId {
		ctx.JSON(http.StatusUnprocessableEntity, "You can not remove yourself to group. To leave a group use the following endpoint: /groups/:id/leave")
		return
	}
	arg := sqlc.RemoveUserFromGroupParams{
		GroupID: req.GroupId,
		UserID:  req.UserId,
	}

	err := server.store.RemoveUserFromGroup(ctx, arg)
	if err != nil {
		pqErr := err.(*pq.Error)
		if string(pqErr.Code) == "NOOWN" {
			ctx.JSON(http.StatusUnauthorized, consts.NotOwner)
			return
		}

		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) getGroupUsers(ctx *gin.Context) {
	var groupId ID
	if err := ctx.ShouldBindUri(&groupId); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.ProvideGroupId)
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId
	arg := sqlc.GetGroupUsersParams{
		UserID:  userId,
		GroupID: groupId.Id,
	}

	users, err := server.store.GetGroupUsers(ctx, arg)
	if err != nil {
		pqErr := err.(*pq.Error)
		if string(pqErr.Code) == "NOTIN" {
			ctx.JSON(http.StatusUnauthorized, consts.NotInGroup)
			return
		}

		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (server *Server) addGroupUser(ctx *gin.Context) {
	var groupId ID
	if err := ctx.ShouldBindUri(&groupId); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.ProvideGroupId)
		return
	}

	var friendUsername string
	if err := ctx.ShouldBindJSON(&friendUsername); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.ProvideGroupId)
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId
	arg := sqlc.AddFriendToGroupParams{
		UserID:         userId,
		GroupID:        groupId.Id,
		FriendUsername: friendUsername,
	}

	err := server.store.AddFriendToGroup(ctx, arg)
	if err != nil {
		pqErr := err.(*pq.Error)
		if string(pqErr.Code) == "NOTIN" {
			ctx.JSON(http.StatusUnauthorized, consts.NotInGroup)
			return
		}

		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (server *Server) deleteGroupUser(ctx *gin.Context) {
	var groupId ID
	if err := ctx.ShouldBindUri(&groupId); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.ProvideGroupId)
		return
	}

	var friendId int64
	if err := ctx.ShouldBindJSON(&friendId); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.ProvideGroupId)
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId
	arg := sqlc.RemoveUserFromGroupParams{
		UserID:   userId,
		GroupID:  groupId.Id,
		FriendID: friendId,
	}

	err := server.store.RemoveUserFromGroup(ctx, arg)
	if err != nil {
		pqErr := err.(*pq.Error)
		if string(pqErr.Code) == "NOTIN" {
			ctx.JSON(http.StatusUnauthorized, consts.NotInGroup)
			return
		}

		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.Status(http.StatusOK)
}

// func (server *Server) addUserAsAdmin(ctx *gin.Context) {
// 	var req groupIdUserIdRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		log.Print(err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// }

// func (server *Server) removeUserAsAdmin(ctx *gin.Context) {
// 	var req groupIdUserIdRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		log.Print(err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// }
