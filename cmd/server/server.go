package main

import (
	"encoding/json"
	"log"
	"net/http"

	"Skillfactory/Censorship/pkg/censor"
	"Skillfactory/Censorship/pkg/middleware"
)

type CheckRequest struct {
	Text string `json:"text"`
}

type CheckResponse struct {
	Allowed bool   `json:"allowed"`
	Message string `json:"message,omitempty"`
}

func main() {
	censorService := censor.NewCensor()

	mux := http.NewServeMux()
	mux.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CheckRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Text == "" {
			http.Error(w, "Text is required", http.StatusBadRequest)
			return
		}

		allowed := censorService.CheckText(req.Text)

		response := CheckResponse{
			Allowed: allowed,
		}

		if !allowed {
			response.Message = "Comment contains banned words"
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	handler := middleware.CORSMiddleware(
		middleware.RequestIDMiddleware(
			middleware.LoggingMiddleware(mux),
		),
	)

	log.Println("Censorship Service starting on :8083")
	log.Fatal(http.ListenAndServe(":8083", handler))
}
