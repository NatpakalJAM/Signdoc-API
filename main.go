package main

import (
	"log"
	"signdoc_api/config"
	"signdoc_api/custom_cache"
	"signdoc_api/db"
	"signdoc_api/gcp"
	"signdoc_api/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/utahta/go-cronowriter"
)

func init() {
	config.Init()
	db.Init()
	gcp.Init()
	// redis.Init()
	custom_cache.Init()

	// now := time.Now()
	// fmt.Println(common.GenerateFileMeta(now, "test"))
	// fmt.Println(common.GenerateFileMeta(now, "test"))
}

func main() {
	app := fiber.New()
	/*
		if cfg.C.Environment != "development" {
			app.Settings.Prefork = true
		}*/
	// app.Use(helmet.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use(requestid.New())
	app.Use(logger.New(prepareLogger()))

	route.Init(app)

	go func() {
		// queue.StartProcessConsumer()
	}()

	log.Fatal(app.Listen(":3000"))
}

func prepareLogger() (config logger.Config) {
	outfile := cronowriter.MustNew("./logs/%Y-%m-%d.log", cronowriter.WithMutex())
	config.Format = "[${time}] ${status} - ${latency} [${method}] ${path} header:[username=${reqHeader:username}] queryParams:[${queryParams}] body:[${body}]\n"
	config.Output = outfile
	return config
}
