package api

import (
	"dipl/db"
	"dipl/token"
	"dipl/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      *db.Store
	tokenMaker *token.JWTMaker
	router     *gin.Engine
}

func NewServer(config util.Config, store *db.Store) (*Server, error) {
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

	router.GET("/users/:id", server.getUser)
	router.POST("/users", server.createUser)
	router.DELETE("/users/:id", server.deleteUser)

	// router.GET("/friends", server.getFriends)
	router.GET("/friends/:id")
	router.POST("/friends/:id")
	router.DELETE("/friends/:id")

	router.GET("/groups")
	router.GET("/groups/:id")
	router.POST("/groups/:id")
	router.DELETE("/groups/:id")

	router.GET("")
	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
