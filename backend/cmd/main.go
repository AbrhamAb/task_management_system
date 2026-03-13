package main

import (
	"log"
	"net/http" // this is go's http server library, not to be confused with the internal/handler/rest/http package
	"os"

	"task-management-backend/internal/glue/routing"
	"task-management-backend/internal/handler/rest"
	"task-management-backend/internal/module"
	"task-management-backend/platform/database"
	"task-management-backend/platform/logger"

	"go.uber.org/zap" // for structured logging, used by our custom logger package
)

func main() {
	appLogger, err := logger.New() 
	if err != nil {
		log.Fatalf("init logger: %v", err)
	} // If logger initialization fails: Stop application and print error
	defer func() { _ = appLogger.Sync() }() // Ensure all logs are flushed before application exits

	dbURL := os.Getenv("DATABASE_URL") /* := is a short variable declaration, it declares 
	and initializes dbURL in one step. 
	It's only used within the main function, so it's appropriate here. 
	If DATABASE_URL is not set,
	log a fatal error and exit the application.*/
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	dbConn, err := database.Connect(dbURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer dbConn.Close() // defer means run this when function exits, ensuring the database connection is properly closed when the application shuts down.
 
	modules := module.New(dbConn) // They interact with repositories and database.
	handlers := rest.New(modules) //Handlers process HTTP requests.
	router := routing.NewRouter(handlers) //Router sends requests to handlers.

	port := os.Getenv("PORT") //Reads server port from environment.
	if port == "" {
		port = "8080"
	}

	appLogger.Info("server started", zap.String("port", port)) 
	if err = http.ListenAndServe(":"+port, router); err != nil { //This starts the HTTP API server.
		log.Fatalf("start server: %v", err)
	}
}
