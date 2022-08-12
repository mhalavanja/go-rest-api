package api

import (
	"dipl/db"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
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
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
