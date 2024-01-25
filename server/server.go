package server

import (
  "context"
	"fmt"
	"github.com/Rizkyyullah/pay-simple/auth"
	"github.com/Rizkyyullah/pay-simple/configs"
	"github.com/Rizkyyullah/pay-simple/middlewares"
	"github.com/Rizkyyullah/pay-simple/users"
	"github.com/Rizkyyullah/pay-simple/products"
	"github.com/Rizkyyullah/pay-simple/shared/services"
	"github.com/Rizkyyullah/pay-simple/transactions"
	"github.com/Rizkyyullah/pay-simple/transactions_detail"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	authUC     auth.UseCase
	usersUC    users.UseCase
	productsUC    products.UseCase
	transactionsUC    transactions.UseCase
	transactionsDetailUC    transactions_detail.UseCase
	jwtService services.JwtService
	engine     *gin.Engine
	address    string
}

func (s *Server) initRoute() {
	v1 := s.engine.Group(configs.APIGroup)
  authMiddleware := middlewares.NewAuthMiddleware(s.jwtService)
  
  auth.NewController(v1, s.authUC).Route()
  users.NewController(v1, s.usersUC, authMiddleware).Route()
  products.NewController(v1, s.productsUC, authMiddleware).Route()
  transactions.NewController(v1, s.transactionsUC, authMiddleware).Route()
}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.address); err != nil {
		log.Fatalf("server not running on address %s, because error %v", s.address, err.Error())
	}
}

func NewServer() *Server {
	tokenConfig := configs.LoadConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta", configs.ENV.DB_Host, configs.ENV.DB_Port, configs.ENV.DB_User, configs.ENV.DB_Password, configs.ENV.DB_Name)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal("server.Connect Err :", err)
	}
	log.Printf("You are now connected to database '%s' as user '%s'", configs.ENV.DB_Name, configs.ENV.DB_User)

  // service
  jwtService := services.NewJwtService(tokenConfig)

	// Repo
	usersRepo := users.NewRepository(conn)
	productsRepo := products.NewRepository(conn)
	transactionsRepo := transactions.NewRepository(conn)
	transactionsDetailRepo := transactions_detail.NewRepository(conn)
	
	// UseCase
	authUC := auth.NewUseCase(usersRepo, jwtService)
	usersUC := users.NewUseCase(usersRepo, jwtService)
	productsUC := products.NewUseCase(productsRepo, usersUC, jwtService)
	transactionsDetailUC := transactions_detail.NewUseCase(transactionsDetailRepo, productsUC)
	transactionsUC := transactions.NewUseCase(transactionsRepo, transactionsDetailUC, productsUC, usersUC)

	engine := gin.Default()
	address := fmt.Sprintf("%s:%d", configs.ENV.API_Host, configs.ENV.API_Port)

	return &Server{
		authUC,
		usersUC,
		productsUC,
		transactionsUC,
		transactionsDetailUC,
		jwtService,
		engine,
		address,
	}
}
