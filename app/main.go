package main

import (
	"context"
	"goKreditPintar/config"
	"goKreditPintar/helper"
	"strconv"

	_DeliveryHTTPAuth "goKreditPintar/kp/delivery/http"
	_RepoMySQLAuth "goKreditPintar/kp/repository/mysql"
	_RepoRedisAuth "goKreditPintar/kp/repository/redis"
	_UsecaseAuth "goKreditPintar/kp/usecase"

	"database/sql"
	"flag"
	"fmt"
	"net/url"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// CLI options parse
	configFile := flag.String("c", "config.yaml", "Config file")
	flag.Parse()

	// Config file
	config.ReadConfig(*configFile)

	// Set log level
	switch viper.GetString("server.log_level") {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	// Initialize database
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", viper.GetString("mysql.user"), viper.GetString("mysql.password"), viper.GetString("mysql.host"), viper.GetString("mysql.port"), viper.GetString("mysql.database"))
	val := url.Values{}
	val.Add("multiStatements", "true")
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Migrate database if any new schema
	driver, err := mysql.WithInstance(dbConn, &mysql.Config{})
	if err == nil {
		mig, err := migrate.NewWithDatabaseInstance(viper.GetString("mysql.path_migrate"), viper.GetString("mysql.database"), driver)
		log.Info(viper.GetString("mysql.path_migrate"))
		if err == nil {
			err = mig.Up()
			if err != nil {
				if err == migrate.ErrNoChange {
					log.Debug("No database migration")
				} else {
					log.Error(err)
				}
			} else {
				log.Info("Migrate database success")
			}
			version, dirty, err := mig.Version()
			if err != nil && err != migrate.ErrNilVersion {
				log.Error(err)
			}
			log.Debug("Current DB version: " + strconv.FormatUint(uint64(version), 10) + "; Dirty: " + strconv.FormatBool(dirty))
		} else {
			log.Warn(err)
		}
	} else {
		log.Warn(err)
	}

	// Initialize Redis
	dbRedis := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
		Username: viper.GetString("redis.username"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.database"),
		PoolSize: viper.GetInt("redis.max_connection"),
	})

	_, err = dbRedis.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Redis connection established")

	repoMySQLAuth := _RepoMySQLAuth.NewMySQLAuthRepository(dbConn)
	repoRedisAuth := _RepoRedisAuth.NewRedisAuthRepository(dbRedis)

	usecaseAuth := _UsecaseAuth.NewAuthUsecase(repoMySQLAuth, repoRedisAuth)

	// Initialize HTTP web framework
	app := fiber.New(fiber.Config{
		Prefork:       viper.GetBool("server.prefork"),
		StrictRouting: viper.GetBool("server.strict_routing"),
		CaseSensitive: viper.GetBool("server.case_sensitive"),
		BodyLimit:     viper.GetInt("server.body_limit"),
	})

	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		Done: func(c *fiber.Ctx, logString []byte) {
			// Get the client IP using the GetClientIP function
			clientIP, err := helper.GetClientIP(c)
			if err != nil {
				log.Errorf("Error getting client IP: %s", err)
				return
			}
			log.Infof("| Client IP: %s |", clientIP)
		},
	}))
	app.Use(recover.New())

	// HTTP routing
	app.Get(viper.GetString("server.base_path")+"/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	if viper.GetBool("api_spec") {
		_DeliveryHTTPAuth.RouterOpenAPI(app)
	}
	_DeliveryHTTPAuth.RouterAPI(app, usecaseAuth)

	// go func() {
	if err := app.Listen(":" + viper.GetString("server.port")); err != nil {
		log.Fatal(err)
	}
	// }()

	// Wait for interrupt signal to gracefully shutdown the server
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// <-quit
	// log.Info("Gracefully shutdown")
	// app.Shutdown()
}
