package main

import (
	"gabrieleromanato/go-smtp-server/app"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/emersion/go-smtp"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	errEnv := godotenv.Load(".env")

	if errEnv != nil {
		log.Fatalf("Error while loading the .env file: %v", errEnv)
	}
	dbName := os.Getenv("DB_NAME")
	store, err := app.NewEmailStore(dbName)
	if err != nil {
		log.Fatalf("Error while creating the database: %v", err)
	}
	defer store.Db.Close()

	go func() {
		backend := &app.Backend{Store: store}
		maxMessageSize, _ := strconv.Atoi(os.Getenv("MAX_MESSAGE_SIZE"))

		server := smtp.NewServer(backend)
		server.Addr = ":" + os.Getenv("SMTP_SERVER_PORT")
		server.Domain = os.Getenv("SMTP_SERVER_HOST")
		server.MaxMessageBytes = int64(maxMessageSize)
		server.AllowInsecureAuth = true

		listener, err := net.Listen("tcp", server.Addr)
		if err != nil {
			log.Fatalf("Error while starting the SMTP server: %v", err)
		}
		defer listener.Close()

		log.Printf("SMTP server listening at %s", server.Addr)
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Error in the SMTP server: %v", err)
		}
	}()

	r := gin.Default()

	r.Static("/attachments", "./attachments")
	r.Use(app.CORSMiddleware())

	r.GET("/check-token", app.CheckTokenExpiration)
	r.POST("/login", app.HandleLogin)
	r.Use(app.AuthMiddleware())

	api := r.Group("/api")
	{
		api.GET("/emails", app.GetEmails(store))
		api.GET("/search", app.SearchForEmails(store))
		api.DELETE("/emails/:id", app.DeleteEmail(store))
		api.POST("/emails", app.SendEmail(store))
	}
	webServerPort := "0.0.0.0:" + os.Getenv("WEB_SERVER_PORT")
	if err := r.Run(webServerPort); err != nil {
		log.Fatalf("Error while starting the web server: %v", err)
	}

}
