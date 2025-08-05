package server

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Rizkyyullah/pay-simple/auth"
	"github.com/Rizkyyullah/pay-simple/configs"
	"github.com/Rizkyyullah/pay-simple/middlewares"
	"github.com/Rizkyyullah/pay-simple/products"
	"github.com/Rizkyyullah/pay-simple/shared/services"
	"github.com/Rizkyyullah/pay-simple/transactions"
	"github.com/Rizkyyullah/pay-simple/transactions_detail"
	"github.com/Rizkyyullah/pay-simple/users"

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
  logMiddleware := middlewares.NewLogMiddleware()
  v1.Use(logMiddleware.ActivityLogs())

  auth.NewController(v1, s.authUC, s.jwtService).Route()
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

func getDatabaseConnectionString() string {
	if configs.ENV.DATABASE_URL != "" {
		log.Print("server: Using connection string from DATABASE_URL")
		return configs.ENV.DATABASE_URL
	}
	
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		configs.ENV.DB_Host,
		configs.ENV.DB_Port,
		configs.ENV.DB_User,
		configs.ENV.DB_Password,
		configs.ENV.DB_Name,
		configs.ENV.Timezone,
	)
	
	log.Printf("Using constructed connection string for host: %s", configs.ENV.DB_Host)
	return dsn
}

func NewServer() *Server {
	tokenConfig := configs.LoadConfig()
	
	dsn := getDatabaseConnectionString()
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("server.Connect Err: %v", err)
	}
	
	if err := conn.Ping(context.Background()); err != nil {
		log.Fatalf("server.Ping Err: %v", err)
	}

	log.Printf("Successfully connected to database")

  jwtService := services.NewJwtService(tokenConfig)

	usersRepo := users.NewRepository(conn)
	productsRepo := products.NewRepository(conn)
	transactionsRepo := transactions.NewRepository(conn)
	transactionsDetailRepo := transactions_detail.NewRepository(conn)
	
	authUC := auth.NewUseCase(usersRepo, jwtService)
	usersUC := users.NewUseCase(usersRepo, jwtService)
	productsUC := products.NewUseCase(productsRepo, usersUC, jwtService)
	transactionsDetailUC := transactions_detail.NewUseCase(transactionsDetailRepo, productsUC)
	transactionsUC := transactions.NewUseCase(transactionsRepo, transactionsDetailUC, productsUC, usersUC)

	engine := gin.Default()
	
	port := os.Getenv("PORT")
	if port == "" {
		if configs.ENV.API_Port != "" {
			port = fmt.Sprint(configs.ENV.API_Port)
		} else {
			port = "5000"
		}
	}
	
	address := fmt.Sprintf(":%s", port)
	log.Printf("Server will start on %s", address)

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