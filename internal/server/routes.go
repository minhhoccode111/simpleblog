package server

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()
	r.Use(s.corsMiddleware)
	r.Use(s.basicAuthentication)

	r.HandleFunc("/", s.HelloWorldHandler)

	return r
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().
			Set("Access-Control-Allow-Origin", "*")
			// Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().
			Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().
			Set("Access-Control-Allow-Credentials", "false")
			// Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) basicAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// first request will have no authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// remove prefix "Authorization: Basic " from header
		encodedBase64 := authHeader[len("Basic "):]

		payload, err := base64.StdEncoding.DecodeString(encodedBase64)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(string(payload), ":", 2)
		if parts[0] != "admin" || parts[1] != "admin" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello, world"))
}

func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request)           {}
func (s *Server) AboutHandler(w http.ResponseWriter, r *http.Request)          {}
func (s *Server) GetAllArticlesHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) GetArticleHandler(w http.ResponseWriter, r *http.Request)     {}

func (s *Server) AdminIndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/articles", http.StatusFound)
}

func (s *Server) AdminGetAllArticlesHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) AdminCreateArticleHandler(w http.ResponseWriter, r *http.Request)  {}
func (s *Server) AdminGetArticleHandler(w http.ResponseWriter, r *http.Request)     {}
func (s *Server) AdminEditArticleHandler(w http.ResponseWriter, r *http.Request)    {}
func (s *Server) AdminDeleteArticleHandler(w http.ResponseWriter, r *http.Request)  {}
