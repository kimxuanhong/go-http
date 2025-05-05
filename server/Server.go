package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server defines server operations.
type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
	Engine() *gin.Engine
	RegisterMiddleware(middleware ...gin.HandlerFunc)
	RegisterRoutes(register func(rg *gin.RouterGroup))
	RegisterPrivateRoutes(register func(rg *gin.RouterGroup), middleware ...gin.HandlerFunc)
	RegisterRoute(method, path string, handler gin.HandlerFunc)
	Routes(routes []RouteConfig)
}

type server struct {
	engine     *gin.Engine
	config     *Config
	httpServer *http.Server
}

// NewServer initializes and returns a new Server instance.
//
// Example:
//
//	cfg := config.NewServerConfig()
//	srv := server.NewServer(cfg)
//	srv.RegisterRoutes(func(rg *gin.RouterGroup) {
//	    rg.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
//	})
//	if err := srv.Start(); err != nil {
//	    log.Fatal(err)
//	}
func NewServer(cfg *Config) Server {
	gin.SetMode(cfg.Mode)

	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	return &server{
		engine: engine,
		config: cfg,
	}
}

// Start runs the HTTP server.
//
// Example:
//
//	if err := srv.Start(); err != nil {
//	    log.Fatal(err)
//	}
func (s *server) Start() error {
	addr := s.config.GetAddr()
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	log.Printf("Server is running at %s", addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server with context.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	if err := srv.Shutdown(ctx); err != nil {
//	    log.Fatal(err)
//	}
func (s *server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

// Engine returns the underlying *gin.Engine.
//
// Example:
//
//	engine := srv.Engine()
//	engine.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
func (s *server) Engine() *gin.Engine {
	return s.engine
}

// RegisterMiddleware adds custom middleware.
//
// Example:
//
//	srv.RegisterMiddleware(cors.Default())
func (s *server) RegisterMiddleware(middleware ...gin.HandlerFunc) {
	s.engine.Use(middleware...)
}

// RegisterRoutes registers public routes.
//
// Example:
//
//	srv.RegisterRoutes(func(rg *gin.RouterGroup) {
//	    rg.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
//	})
func (s *server) RegisterRoutes(register func(rg *gin.RouterGroup)) {
	public := s.engine.Group("/")
	register(public)
}

// RegisterPrivateRoutes registers authenticated routes with middleware.
//
// Example:
//
//	srv.RegisterPrivateRoutes(func(rg *gin.RouterGroup) {
//	    rg.GET("/profile", func(c *gin.Context) { c.JSON(200, gin.H{"user": "current user"}) })
//	}, authMiddleware)
func (s *server) RegisterPrivateRoutes(register func(rg *gin.RouterGroup), middleware ...gin.HandlerFunc) {
	private := s.engine.Group("/private")
	private.Use(middleware...)
	register(private)
}

// RegisterRoute quickly registers a single route.
//
// Example:
//
//	srv.RegisterRoute("GET", "/ping", func(c *gin.Context) {
//	    c.JSON(200, gin.H{"message": "pong"})
//	})
func (s *server) RegisterRoute(method, path string, handler gin.HandlerFunc) {
	switch method {
	case "GET":
		s.engine.GET(path, handler)
	case "POST":
		s.engine.POST(path, handler)
	case "PUT":
		s.engine.PUT(path, handler)
	case "PATCH":
		s.engine.PATCH(path, handler)
	case "DELETE":
		s.engine.DELETE(path, handler)
	default:
		log.Printf("Unsupported method: %s", method)
	}
}

// Routes registers a list of routes into the Gin engine of the server.
//
// Each route in the list will:
//   - Be grouped by `Path` using `engine.Group(Path)`
//   - Apply any provided middleware to that group
//   - Register the corresponding handler for the HTTP method through `RegisterRoute`
//
// Note: `RegisterRoute` should handle adding the route to the Gin engine properly.
//
// Parameters:
//   - routes: A list of route configurations, including Path, Method, Middleware, and Handler.
//
// Example:
//
//	routes := []RouteConfig{
//	    {
//	        Path:       "/users/:id",
//	        Method:     http.MethodGet,
//	        HandleFunc: userHandler.GetUser,
//	        Middleware: []gin.HandlerFunc{AuthMiddleware},
//	    },
//	}
//	s.Routes(routes)
func (s *server) Routes(routes []RouteConfig) {
	for _, r := range routes {
		group := s.engine.Group(r.Path)
		group.Use(r.Middleware...)
		s.RegisterRoute(r.Method, r.Path, r.HandleFunc)
	}
}
