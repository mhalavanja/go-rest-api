package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mhalavanja/go-rest-api/db/sqlc"
	"github.com/mhalavanja/go-rest-api/token"
	"github.com/mhalavanja/go-rest-api/util"
)

type Server struct {
	config     *util.Config
	store      *sqlc.Queries
	tokenMaker *token.JWTMaker
	router     *gin.Engine
	hub        *hub
}

func NewServer(config *util.Config, store *sqlc.Queries, hub *hub) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create new JWTMaker: %w", err)
	}

	server := &Server{
		config:     config,
		tokenMaker: tokenMaker,
		store:      store,
		hub:        hub,
	}

	router := gin.Default()
	router.SetTrustedProxies([]string{server.config.Client})

	router.POST("/register", server.createUser)
	router.POST("/tokens/authenticate", server.authUser)
	router.POST("/tokens/renewAccess", server.renewAccessToken)

	authGroup := router.Group("/").Use(authMiddleware(*server.tokenMaker))
	authGroup.DELETE("/tokens/refreshToken", server.deleteRefreshToken)

	authGroup.GET("/user", server.getUser)
	authGroup.DELETE("/user", server.deleteUser)
	authGroup.PUT("/user", server.updateUser)

	authGroup.GET("/friends", server.getFriends)
	authGroup.GET("/friends/:id", server.getFriend)
	authGroup.POST("/friends", server.addFriend)
	authGroup.DELETE("/friends/:id", server.deleteFriend)

	authGroup.GET("/groups", server.getGroups)
	authGroup.GET("/groups/:id", server.getGroup)
	authGroup.GET("/groups/:id/users", server.getGroupUsers)
	authGroup.POST("/groups/:id/users", server.addGroupUser)
	authGroup.DELETE("/groups/:id/users", server.deleteGroupUser)
	authGroup.POST("/groups", server.createGroup)
	authGroup.DELETE("/groups/:id/leave", server.leaveGroup)
	authGroup.DELETE("/groups/:id", server.deleteGroup)
	authGroup.POST("/groups/addUser", server.addFriendToGroup)
	authGroup.DELETE("/groups/removeUser", server.removeUserFromGroup)

	router.GET("/ws/groups/:id", hub.ServeWs)

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
