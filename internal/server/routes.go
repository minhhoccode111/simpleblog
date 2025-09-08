package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()
	r.Use(s.corsMiddleware)

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
		decodeCredentials(r.Header.Get("Authorization"))

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
