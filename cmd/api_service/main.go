package main

import (
	"log"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/asalvi0/bond-trading/internal/api/authn"
	"github.com/asalvi0/bond-trading/internal/api/order"
	"github.com/asalvi0/bond-trading/internal/api/user"
	"github.com/asalvi0/bond-trading/internal/config"
)

func main() {
	app := fiber.New(fiber.Config{
		// Prefork:     true,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	setupMiddleware(app)

	err := authn.RegisterRoutes(app)
	if err != nil {
		log.Fatal(err)
	}

	err = user.RegisterRoutes(app)
	if err != nil {
		log.Fatal(err)
	}

	err = order.RegisterRoutes(app)
	if err != nil {
		log.Fatal(err)
	}

	port := config.Config("API_SERVICE_PORT")
	log.Fatal(app.Listen(":" + port))
}

func setupMiddleware(app *fiber.App) {
	// storage := sqlite3.New()

	app.Use(favicon.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		TimeFormat: "02-Jan-2006 15:04:05",
	}))

	app.Use(idempotency.New())
	app.Use(etag.New())
	app.Use(requestid.New())

	app.Use(pprof.New())
	app.Get("/metrics", monitor.New())

	app.Use(healthcheck.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(limiter.New(limiter.Config{
		// Storage:    storage,
		Max:        1000,
		Expiration: 1 * time.Minute,
	}))

	app.Use(helmet.New())
	// app.Use(csrf.New())
}
