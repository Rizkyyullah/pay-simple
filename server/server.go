package server

import (
  "context"
	"fmt"
	"github.com/Rizkyyullah/pay-simple/configs"
	"github.com/Rizkyyullah/pay-simple/users"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	usersUC users.UseCase
	engine  *gin.Engine
	host    string
}

func (s *Server) initRoute() {
	v1 := s.engine.Group(configs.APIGroup)

  users.NewController(v1, s.usersUC).Route()
}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("server not running on host %s, because error %v", s.host, err.Error()))
	}
}

func NewServer() *Server {
	configs.LoadConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta", configs.ENV.Host, configs.ENV.Port, configs.ENV.User, configs.ENV.Password, configs.ENV.Name)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal("server.Connect Err :", err)
	}
	log.Printf("You are now connected to database '%s' as user '%s'", configs.ENV.Name, configs.ENV.User)

	// Repo
	usersRepo := users.NewRepository(conn)
	
	// Usecase
	usersUC := users.NewUseCase(usersRepo)
	
	engine := gin.Default()
	host := fmt.Sprintf("%s:%d", configs.ENV.Host, configs.ENV.APIPort)

	return &Server{
		usersUC:    usersUC,
		engine:     engine,
		host:       host,
	}
}
