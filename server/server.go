package server

import (
  "context"
	"fmt"
	"github.com/Rizkyyullah/pay-simple/auth"
	"github.com/Rizkyyullah/pay-simple/configs"
	"github.com/Rizkyyullah/pay-simple/users"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	authUC     auth.UseCase
	jwtService auth.JwtService
	engine     *gin.Engine
	address    string
}

func (s *Server) initRoute() {
	v1 := s.engine.Group(configs.APIGroup)

  auth.NewController(v1, s.authUC).Route()
}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.address); err != nil {
		log.Fatalf("server not running on address %s, because error %v", s.address, err.Error())
	}
}

func NewServer() *Server {
	tokenConfig := configs.LoadConfig()
	var dsn string
	if os.Getenv("APPMODE") == "DEPLOY" {
	  dsn = fmt.Sprintf(os.Getenv("DATABASE_URL"))
	} else {
	  dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta", configs.ENV.DB_Host, configs.ENV.DB_Port, configs.ENV.DB_User, configs.ENV.DB_Password, configs.ENV.DB_Name)
	}
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal("server.Connect Err :", err)
	}
	log.Printf("You are now connected to database '%s' as user '%s'", configs.ENV.DB_Name, configs.ENV.DB_User)

  // service
  jwtService := auth.NewJwtService(tokenConfig)

	// Repo
	usersRepo := users.NewRepository(conn)
	
	// UseCase
	authUC := auth.NewUseCase(usersRepo, jwtService)
	
	engine := gin.Default()
	address := fmt.Sprintf("%s:%d", configs.ENV.API_Host, configs.ENV.PORT)

	return &Server{
		authUC:     authUC,
		engine:     engine,
		address:    address,
	}
}
