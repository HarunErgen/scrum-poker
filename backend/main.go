package main

import (
	"fmt"
	"github.com/scrum-poker/backend/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/handlers"
)

func main() {
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	websocket.Init()

	r := mux.NewRouter()

	r.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	r.HandleFunc("/rooms", handlers.CreateRoomHandler).Methods("POST")
	r.HandleFunc("/rooms/{roomId}", handlers.GetRoomHandler).Methods("GET")
	r.HandleFunc("/rooms/{roomId}/join", handlers.JoinRoomHandler).Methods("POST")

	r.HandleFunc("/sessions/{userId}/{roomId}", handlers.CreateSessionHandler).Methods("POST")
	r.HandleFunc("/sessions", handlers.GetSessionHandler).Methods("GET")
	r.HandleFunc("/sessions", handlers.DeleteSessionHandler).Methods("DELETE")

	r.HandleFunc("/ws/{roomId}", handlers.WebSocketHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		fmt.Printf("Server is running on port %s...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
