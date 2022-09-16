package api

import (
	"dipl/db/sqlc"
	"dipl/token"
	"dipl/util"
	"fmt"
	"github.com/gin-gonic/gin"
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
	router.GET("")
	router.POST("/signup", server.createUser)
	router.POST("/authenticate", server.authUser)

	authGroup := router.Group("/").Use(authMiddleware(*server.tokenMaker))

	authGroup.GET("/users/:id", server.getUser)
	authGroup.DELETE("/users/:id", server.deleteUser)

	// authGroup.GET("/friends", server.getFriends)
	authGroup.GET("/friends/:id")
	authGroup.POST("/friends/:id")
	authGroup.DELETE("/friends/:id")

	authGroup.GET("/groups")
	authGroup.GET("/groups/:id")
	authGroup.POST("/groups/:id")
	authGroup.DELETE("/groups/:id")

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
