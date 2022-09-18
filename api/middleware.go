package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"dipl/token"

	"github.com/gin-gonic/gin"
)

const authPayload = "authorization_payload"

func authMiddleware(tokenMaker token.JWTMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("authorization")

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			log.Println(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			log.Println(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != "bearer" {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			log.Println(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authPayload, payload)
		ctx.Next()
	}
}