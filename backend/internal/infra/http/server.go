package http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kvloginov/cup-of-team/backend/internal/api/model"
)

// Config holds server configuration
type Config struct {
	Port string
}

// Server represents the HTTP server
type Server struct {
	config   Config
	router   *mux.Router
	handlers map[string]http.HandlerFunc
}

// NewServer creates a new server instance
func NewServer(config Config) *Server {
	return &Server{
		config:   config,
		router:   mux.NewRouter(),
		handlers: make(map[string]http.HandlerFunc),
	}
}

// Handle registers a handler with CORS middleware
func (s *Server) Handle(pattern string, handler http.HandlerFunc) {
	s.handlers[pattern] = handler
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// API routes - all API handlers registered without /api prefix
	apiRouter := s.router.PathPrefix("/api").Subrouter()
	for pattern, handler := range s.handlers {
		// Skip /health as it should be on root level
		if pattern == "/health" {
			continue
		}
		apiRouter.HandleFunc(pattern, corsMiddleware(handler))
	}

	// Health check - outside /api prefix
	if healthHandler, ok := s.handlers["/health"]; ok {
		s.router.HandleFunc("/health", corsMiddleware(healthHandler)).Methods("GET")
	}

	// Get frontend path from environment or use default
	frontendPath := os.Getenv("FRONTEND_PATH")
	if frontendPath == "" {
		frontendPath = "./frontend/dist"
	}

	log.Printf("Serving frontend from: %s", frontendPath)

	// Serve static frontend files
	// Serve index.html for root path
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, frontendPath+"/index.html")
	}).Methods("GET")

	// Serve other static files
	fileServer := http.FileServer(http.Dir(frontendPath))
	s.router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Setup routes before starting
	s.setupRoutes()

	addr := s.config.Port
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, s.router)
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Log origin for debugging
		if origin != "" {
			log.Printf("[CORS] Request from origin: %s", origin)
		}

		// Allow requests from any origin (including Capacitor apps)
		// Capacitor typically uses origins like:
		// - capacitor://localhost
		// - http://localhost
		// - ionic://localhost
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// SendJSON sends a JSON response
func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// SendError sends an error response
func SendError(w http.ResponseWriter, status int, message string) {
	SendJSON(w, status, model.ErrorResponse{
		Success: false,
		Error:   message,
	})
}
