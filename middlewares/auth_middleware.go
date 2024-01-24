package middlewares

import (
	"github.com/Rizkyyullah/pay-simple/shared/services"
	"github.com/Rizkyyullah/pay-simple/shared/common"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	RequireToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService services.JwtService
}

type AuthHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

func (m *authMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context){
		var authHeader AuthHeader
		if err := ctx.ShouldBindHeader(&authHeader); err != nil{
			log.Println("middlewares.RequireToken.authHeader Err :", err.Error())
			common.SendUnauthorizedResponse(ctx, "Unauthorized : " + err.Error())
			return
		}

		tokenHeader := strings.Replace(authHeader.AuthorizationHeader, "Bearer ", "", -1)
		if tokenHeader == ""{
			log.Printf("middlewares.RequireToken.tokenHeader \n")
			common.SendUnauthorizedResponse(ctx, "Unauthorized : Header tokens should not be empty")
			return
		}

		claims, err := m.jwtService.ParseToken(tokenHeader)
		if err != nil {
			log.Printf("middlewares.RequireToken.ParseToken: %v \n", err.Error())
			common.SendUnauthorizedResponse(ctx, "Unauthorized : " + err.Error())
			return
		}

		ctx.Set("userId", claims["userId"])

		validRole := false
		for _, role := range roles{
			if role == claims["role"]{
				validRole = true
				break
			}
		}

		if !validRole {
			log.Printf("RequireToken.validRole \n")
			common.SendForbiddenResponse(ctx, "You are prohibited from accessing these resources")
			return
		}

		ctx.Next()
	}
}

func NewAuthMiddleware(jwtService services.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService}
}
