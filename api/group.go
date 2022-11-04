package api

import (
	"log"
	"net/http"

	"github.com/mhalavanja/go-rest-api/db/sqlc"
	"github.com/mhalavanja/go-rest-api/token"

	"github.com/gin-gonic/gin"
)

func (server *Server) getGroups(ctx *gin.Context) {
	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	groupNames, err := server.store.GetGroups(ctx, userId)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, groupNames)
}

func (server *Server) getGroup(ctx *gin.Context) {
	var id int64
	if err := ctx.ShouldBindUri(&id); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	arg := sqlc.GetGroupParams{
		ID:          id,
		UserIDOwner: userId,
	}

	group, err := server.store.GetGroup(ctx, arg)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, group)
}

func (server *Server) createGroup(ctx *gin.Context) {
	var name string
	if err := ctx.ShouldBindJSON(&name); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	arg := sqlc.CreateGroupParams{
		GroupName: name,
		UserID:    userId,
	}

	err := server.store.CreateGroup(ctx, arg)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusCreated)
}

func (server *Server) deleteGroup(ctx *gin.Context) {
	var id int64
	if err := ctx.ShouldBindUri(&id); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	arg := sqlc.TryDeleteGroupParams{
		GroupID: id,
		UserID:  userId,
	}

	err := server.store.TryDeleteGroup(ctx, arg)
	if err != nil {
		//TODO: Hvatati error ako korisnik pokusa brisat grupu koja nije njegova, baca se iz procedure
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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

// func (server *Server) joinGroup(ctx *gin.Context) {
// 	var req groupIdRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		log.Print(err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// }

func (server *Server) leaveGroup(ctx *gin.Context) {
	var groupId int64
	if err := ctx.ShouldBindUri(&groupId); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	arg := sqlc.LeaveGroupParams{
		GroupID: groupId,
		UserID:  userId,
	}

	err := server.store.LeaveGroup(ctx, arg)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) addFriendToGroup(ctx *gin.Context) {
	var req groupIdUserIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId
	if userId == req.UserId {
		ctx.JSON(http.StatusUnprocessableEntity, "Can not add yourself to group")
		return
	}
	arg := sqlc.AddFriendToGroupParams{
		GroupID: req.GroupId,
		UserID:  req.UserId,
	}

	err := server.store.AddFriendToGroup(ctx, arg)
	if err != nil {
		log.Print(err.Error())
		// TODO: Dodati error koji se baca iz funkcije kada nisu prijatelji i hvatati ga ovdje te vracat prikladnu poruku
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) removeUserFromGroup(ctx *gin.Context) {
	var req groupIdUserIdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId
	if userId == req.UserId {
		ctx.JSON(http.StatusUnprocessableEntity, "Can not remove yourself to group. To leave a group use the following endpoint: /groups/:id/leave")
		return
	}
	arg := sqlc.RemoveUserFromGroupParams{
		GroupID: req.GroupId,
		UserID:  req.UserId,
	}

	err := server.store.RemoveUserFromGroup(ctx, arg)
	if err != nil {
		log.Print(err.Error())
		// TODO: Dodati error koji se baca iz funkcije kada user nije u toj grupi
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
