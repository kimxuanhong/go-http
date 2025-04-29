# Product Gin Module

A module for managing HTTP server using Gin framework.

## Features

- HTTP server management
- Environment-based configuration
- Middleware support
- Graceful shutdown support

## Configuration

The module uses environment variables for server configuration:

```bash
SERVER_HOST=localhost
SERVER_PORT=8080
GIN_MODE=debug
```

## Usage

1. Import the module in your project:
```go
import "github.com/yourusername/product-gin"
```

2. Initialize and start the server:
```go
cfg := config.NewServerConfig()
server := server.NewServer(cfg)

// Add routes
server.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "pong",
    })
})

// Start server
if err := server.Start(); err != nil {
    log.Fatal(err)
}
```

## Dependencies

- Go 1.21 or later
- Gin Web Framework

## License

MIT 