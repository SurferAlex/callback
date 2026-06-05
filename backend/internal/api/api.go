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
	r.Use(middleware.CORS(cfg.CORSOrigins))
	r.Use(middleware.RequestLog())

	serversRepo := repository.NewVPNServersRepo(db)
	serversUC := usecase.NewVPNServers(serversRepo)

	clientsRepo := repository.NewVPNClientsRepo(db)
	clientsUC := usecase.NewVPNClients(clientsRepo)

	registry := usecase.NewXUIRegistry(serversUC)
	xuiRepo := repository.NewXUIAccessRepo(db)
	xuiUC := usecase.NewXUIAccess(xuiRepo, clientsUC, registry)

	usersRepo := repository.NewUsersRepo(db)
	subsRepo := repository.NewSubscriptionsRepo(db)
	trialsRepo := repository.NewTrialActivationsRepo(db)
	refreshRepo := repository.NewAuthRefreshRepo(db)
	authUC := usecase.NewAuthSession(usersRepo, refreshRepo, cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	usersUC := usecase.NewUserService(usersRepo, subsRepo, trialsRepo, clientsUC, serversUC, xuiUC, cfg.DefaultVPNServer, cfg.DefaultMaxIPs)

	authH := &handlers.AuthHandlers{
		Auth:         authUC,
		Users:        usersUC,
		BotToken:     cfg.TelegramBotToken,
		CookieDomain: cfg.CookieDomain,
		CookieSecure: cfg.CookieSecure,
		RefreshTTL:   cfg.JWTRefreshTTL,
	}

	h := &handlers.Handlers{
		DB:        db,
		Servers:   serversUC,
		Clients:   clientsUC,
		XUIAccess: xuiUC,
		Users:     usersUC,
	}

	RegisterRoutes(r, h, authH, authUC, cfg)

	return r
}

func RegisterRoutes(r *gin.Engine, h *handlers.Handlers, authH *handlers.AuthHandlers, authUC *usecase.AuthSession, cfg config.Config) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", h.Ping)
		v1.GET("/health", h.Health)
	}

	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/session/webapp", authH.SessionTelegramWebApp)
		authGroup.POST("/session/widget", authH.SessionTelegramWidget)
		authGroup.POST("/refresh", authH.Refresh)
		authGroup.POST("/logout", authH.Logout)
	}

	userAuth := middleware.UserAuth(cfg.InternalToken, cfg.TelegramBotToken, authUC)
	user := v1.Group("/user")
	user.Use(userAuth)
	{
		user.GET("/me", h.UserMe)
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
	}

	// Bot/admin: internal token only (no UserAuth JWT/tma) — same paths, separate middleware chain.
	internalAuth := middleware.InternalUserAuth(cfg.InternalToken)
	v1.POST("/user/trial/activate", internalAuth, h.UserTrialActivate)
	v1.POST("/user/subscription/mock-activate", internalAuth, h.UserMockActivate)
}
