package main

import (
	"log"
	"net/http"

	"github.com/clearclown/HaiLanGo/backend/internal/api/websocket"
	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/service/notification"
)

func main() {
	// Create WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Create notification service
	notificationService := notification.NewService(hub)

	// Create WebSocket handler
	wsHandler := websocket.NewHandler(hub)

	// Setup routes
	http.HandleFunc("/api/v1/ws", wsHandler.ServeWS)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Example endpoint to trigger notifications (for testing)
	http.HandleFunc("/api/v1/test/notify", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}

		// Send a test notification
		_ = notificationService.NotifyOCRProgress(userID, &models.OCRProgressData{
			BookID:         "test-book-123",
			TotalPages:     100,
			ProcessedPages: 50,
			CurrentPage:    50,
			Progress:       50.0,
			Status:         "processing",
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Notification sent"))
	})

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
