package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/fakhripraya/user-service/config"
	"github.com/fakhripraya/user-service/data"
	"github.com/fakhripraya/user-service/entities"
	"github.com/fakhripraya/user-service/handlers"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/srinathgs/mysqlstore"
)

var err error

// Session Store based on MYSQL database
var sessionStore *mysqlstore.MySQLStore

// Adapter is an alias
type Adapter func(http.Handler) http.Handler

// Adapt takes Handler funcs and chains them to the main handler.
func Adapt(handler http.Handler, adapters ...Adapter) http.Handler {
	// The loop is reversed so the adapters/middleware gets executed in the same
	// order as provided in the array.
	for i := len(adapters); i > 0; i-- {
		handler = adapters[i-1](handler)
	}
	return handler
}

func main() {

	// creates a structured logger for logging the entire program
	logger := hclog.Default()

	// load configuration from env file
	err = godotenv.Load(".env")

	if err != nil {
		// log the fatal error if load env failed
		log.Fatal(err)
	}

	// Initialize app configuration
	var appConfig entities.Configuration
	data.ConfigInit(&appConfig)

	// Open the database connection based on DB configuration
	logger.Info("Establishing database connection on " + appConfig.Database.Host + ":" + strconv.Itoa(appConfig.Database.Port))
	config.DB, err = gorm.Open("mysql", config.DbURL(config.BuildDBConfig(&appConfig.Database)))
	if err != nil {
		log.Fatal(err)
	}

	defer config.DB.Close()

	// Creates a session store based on MYSQL database
	// If table doesn't exist, creates a new one
	logger.Info("Building session store based on " + appConfig.Database.Host + ":" + strconv.Itoa(appConfig.Database.Port))
	sessionStore, err = mysqlstore.NewMySQLStore(config.DbURL(config.BuildDBConfig(&appConfig.Database)), "dbMasterSession", "/", 3600*24*7, []byte(appConfig.MySQLStore.Secret))
	if err != nil {
		log.Fatal(err)
	}

	defer sessionStore.Close()

	// creates a user instance
	user := data.NewUser(logger)

	// creates the user handler
	userHandler := handlers.NewUserHandler(logger, user, sessionStore)

	// creates a new serve mux
	serveMux := mux.NewRouter()

	// handlers for the API
	logger.Info("Setting handlers for the API")

	// get handlers
	getRequest := serveMux.Methods(http.MethodGet).Subrouter()

	// get user handler
	getRequest.HandleFunc("/", userHandler.GetUser)

	// get global middleware
	getRequest.Use(userHandler.MiddlewareValidateAuth)

	// patch handlers
	patchRequest := serveMux.Methods(http.MethodPatch).Subrouter()

	// update signed in user patch handler
	patchRequest.HandleFunc("/user/update/signed", userHandler.UpdateSignedUser)

	// patch global middleware
	patchRequest.Use(
		userHandler.MiddlewareValidateAuth,
		userHandler.MiddlewareParseUserRequest,
	)

	// CORS
	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// creates a new server
	server := http.Server{
		Addr:         appConfig.API.Host + ":" + strconv.Itoa(appConfig.API.Port), // configure the bind address
		Handler:      corsHandler(serveMux),                                       // set the default handler
		ErrorLog:     logger.StandardLogger(&hclog.StandardLoggerOptions{}),       // set the logger for the server
		ReadTimeout:  5 * time.Second,                                             // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                            // max time to write response to the client
		IdleTimeout:  120 * time.Second,                                           // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		logger.Info("Starting server on port " + appConfig.API.Host + ":" + strconv.Itoa(appConfig.API.Port))

		err = server.ListenAndServe()
		if err != nil {

			if strings.Contains(err.Error(), "http: Server closed") == true {
				os.Exit(0)
			} else {
				logger.Error("Error starting server", "error", err.Error())
				os.Exit(1)
			}
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	signal.Notify(channel, os.Kill)

	// Block until a signal is received.
	sig := <-channel
	logger.Info("Got signal", "info", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
