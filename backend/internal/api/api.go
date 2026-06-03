package api

import (
	"api-vpn/internal/config"
	"api-vpn/internal/handlers"
	"api-vpn/internal/middleware"
	"api-vpn/internal/repository"
	"api-vpn/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupServer(db *pgxpool.Pool, cfg config.Config) *gin.Engine {
	r := gin.Default()

	serversRepo := repository.NewVPNServersRepo(db)
	serversUC := usecase.NewVPNServers(serversRepo)

	clientsRepo := repository.NewVPNClientsRepo(db)
	clientsUC := usecase.NewVPNClients(clientsRepo)

	registry := usecase.NewXUIRegistry(serversUC)
	xuiRepo := repository.NewXUIAccessRepo(db)
	xuiUC := usecase.NewXUIAccess(xuiRepo, clientsUC, registry)

	usersRepo := repository.NewUsersRepo(db)
	subsRepo := repository.NewSubscriptionsRepo(db)
	usersUC := usecase.NewUserService(usersRepo, subsRepo, clientsUC, serversUC, xuiUC, cfg.DefaultVPNServer, cfg.DefaultMaxIPs)

	h := &handlers.Handlers{
		DB:        db,
		Servers:   serversUC,
		Clients:   clientsUC,
		XUIAccess: xuiUC,
		Users:     usersUC,
	}

	RegisterRoutes(r, h, cfg)

	return r
}

func RegisterRoutes(r *gin.Engine, h *handlers.Handlers, cfg config.Config) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", h.Ping)
		v1.GET("/health", h.Health)
	}

	user := v1.Group("/user")
	user.Use(middleware.UserAuth(cfg.InternalToken, cfg.TelegramBotToken))
	{
		user.GET("/me", h.UserMe)
		user.POST("/subscription/mock-activate", h.UserMockActivate)
		user.GET("/config", h.UserGetConfig)
		user.POST("/config/refresh", h.UserRefreshConfig)
	}

	protected := v1.Group("")
	protected.Use(middleware.InternalToken(cfg.InternalToken))
	{
		protected.GET("/servers", h.ListServers)
		protected.POST("/clients", h.CreateClient)
		protected.GET("/clients/resolve", h.ResolveClient)
		protected.GET("/clients/:uuid", h.GetClient)
		protected.POST("/clients/:uuid/deactivate", h.DeactivateClient)
		protected.POST("/clients/:uuid/provision", h.ProvisionAccess)
		protected.GET("/clients/:uuid/access", h.GetAccess)
		protected.POST("/clients/:uuid/revoke", h.RevokeAccess)
		protected.POST("/clients/:uuid/extend", h.ExtendClient)
		protected.POST("/clients/:uuid/max-ips", h.UpdateClientMaxIPs)
		protected.GET("/monitor/targets", h.ListMonitorTargets)
	}
}
