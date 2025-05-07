package main

import (
	"fmt"
	"github.com/scrum-poker/backend/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/scrum-poker/backend/db"
	"github.com/scrum-poker/backend/handlers"
)

func main() {
	port := os.Getenv("PORT")
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
	apiRouter := r.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	apiRouter.HandleFunc("/rooms", handlers.CreateRoomHandler).Methods("POST")
	apiRouter.HandleFunc("/rooms/{roomID}", handlers.GetRoomHandler).Methods("GET")
	apiRouter.HandleFunc("/rooms/{roomID}/join", handlers.JoinRoomHandler).Methods("POST")
	apiRouter.HandleFunc("/rooms/{roomID}/leave", handlers.LeaveRoomHandler).Methods("POST")

	apiRouter.HandleFunc("/rooms/{roomID}/vote", handlers.SubmitVoteHandler).Methods("POST")
	apiRouter.HandleFunc("/rooms/{roomID}/reveal", handlers.RevealVotesHandler).Methods("POST")
	apiRouter.HandleFunc("/rooms/{roomID}/reset", handlers.ResetVotesHandler).Methods("POST")
	apiRouter.HandleFunc("/rooms/{roomID}/transfer", handlers.TransferScrumMasterHandler).Methods("POST")

	r.HandleFunc("/ws/{roomID}", handlers.WebSocketHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
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
