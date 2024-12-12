package main

import (
	"gabrieleromanato/go-smtp-server/app"
	"log"
	"net"
	"os"

	"github.com/emersion/go-smtp"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	store, err := app.NewEmailStore("emails.db")
	if err != nil {
		log.Fatalf("Errore durante l'inizializzazione del database: %v", err)
	}
	defer store.Db.Close()

	err = godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Errore durante il caricamento del file .env: %v", err)
	}

	go func() {
		backend := &app.Backend{Store: store}

		// Configura il server SMTP
		server := smtp.NewServer(backend)
		server.Addr = ":" + os.Getenv("SMTP_SERVER_PORT") // Porta locale
		server.Domain = os.Getenv("SMTP_SERVER_HOST")
		server.AllowInsecureAuth = true

		// Avvia il server
		listener, err := net.Listen("tcp", server.Addr)
		if err != nil {
			log.Fatalf("Errore nell'avviare il server: %v", err)
		}
		defer listener.Close()

		log.Printf("Server SMTP in ascolto su %s", server.Addr)
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Errore nel server SMTP: %v", err)
		}
	}()

	r := gin.Default()

	r.Static("/attachments", "./attachments")
	r.Use(app.CORSMiddleware())

	r.POST("/login", app.HandleLogin)
	r.Use(app.AuthMiddleware())

	api := r.Group("/api")
	{
		api.GET("/emails", app.GetEmails(store))
		api.DELETE("/emails/:id", app.DeleteEmail(store))
		api.POST("/emails", app.SendEmail(store))
	}
	webServerPort := ":" + os.Getenv("WEB_SERVER_PORT")
	if err := r.Run(webServerPort); err != nil {
		log.Fatalf("Errore durante l'avvio del server: %v", err)
	}

}
