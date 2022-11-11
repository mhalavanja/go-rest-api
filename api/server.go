package api

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mhalavanja/go-rest-api/db/sqlc"
	"github.com/mhalavanja/go-rest-api/token"
	"github.com/mhalavanja/go-rest-api/util"
)

type Server struct {
	config     util.Config
	store      *sqlc.Queries
	tokenMaker *token.JWTMaker
	router     *gin.Engine
}

func NewServer(config util.Config, store *sqlc.Queries) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create new JWTMaker: %w", err)
	}
	server := &Server{
		config:     config,
		tokenMaker: tokenMaker,
		store:      store,
	}
	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/register", server.createUser)
	router.POST("/authenticate", server.authUser)

	authGroup := router.Group("/").Use(authMiddleware(*server.tokenMaker))

	authGroup.GET("/user", server.getUser)
	authGroup.DELETE("/user", server.deleteUser)
	authGroup.PUT("/user", server.updateUser)

	authGroup.GET("/friends", server.getFriends)
	authGroup.GET("/friends/:id", server.getFriend)
	authGroup.POST("/friends", server.addFriend)
	authGroup.DELETE("/friends/:id", server.deleteFriend)

	authGroup.GET("/groups", server.getGroups)
	authGroup.GET("/groups/:id", server.getGroup)
	authGroup.POST("/groups", server.createGroup)
	authGroup.DELETE("/groups/leave/:id", server.leaveGroup)
	authGroup.DELETE("/groups/:id", server.deleteGroup)
	authGroup.POST("/groups/addUser", server.addFriendToGroup)
	authGroup.DELETE("/groups/removeUser", server.removeUserFromGroup)

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
