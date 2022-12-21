package api

import (
	"log"
	"net/http"

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
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getGroup(ctx *gin.Context) {
	var id ID
	if err := ctx.ShouldBindUri(&id); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.ProvideGroupId)
		return
	}

	userId := ctx.MustGet(authPayload).(*token.Payload).UserId

	log.Println(id)
	log.Println(userId)
	arg := sqlc.GetGroupParams{
		ID:          id.ID,
		UserIDOwner: userId,
	}

	group, err := server.store.GetGroup(ctx, arg)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, consts.InternalErrorMessage)
		return
	}

	ctx.JSON(http.StatusCreated, group)
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
	var id int64
	if err := ctx.ShouldBindUri(&id); err != nil {
		log.Println("ERROR: ", err.Error())
		ctx.JSON(http.StatusBadRequest, consts.ProvideGroupId)
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
		log.Println("ERROR: ", err.Error())
		// TODO: Dodati error koji se baca iz funkcije kada nisu prijatelji i hvatati ga ovdje te vracat prikladnu poruku
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
		log.Println("ERROR: ", err.Error())
		// TODO: Dodati error koji se baca iz funkcije kada user nije u toj grupi
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
