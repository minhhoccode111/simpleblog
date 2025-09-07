package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()
	r.Use(s.corsMiddleware)

	// Root level routes (no versioning)
	r.HandleFunc("/", s.HelloWorldHandler)

	// API v1 routes
	v1 := r.PathPrefix("/v1").Subrouter()
	s.registerV1Routes(v1)

	return r
}

func (s *Server) registerV1Routes(r *mux.Router) {
	r.HandleFunc("/home", s.HomeHandler).Methods("GET")
	r.HandleFunc("/about", s.AboutHandler).Methods("GET")
	r.HandleFunc("/articles", s.GetAllArticlesHandler).Methods("GET")
	r.HandleFunc("/articles/{id}", s.GetArticleHandler).Methods("GET")

	r.Use(s.basicAuthentication)

	r.HandleFunc("/admin", s.AdminIndexHandler).Methods("GET")
	r.HandleFunc("/admin/articles", s.AdminGetAllArticlesHandler).Methods("GET")
	r.HandleFunc("/admin/articles", s.AdminCreateArticleHandler).Methods("POST")
	r.HandleFunc("/admin/articles/{id}", s.AdminGetArticleHandler).Methods("GET")
	r.HandleFunc("/admin/articles/{id}", s.AdminEditArticleHandler).Methods("PUT")
	r.HandleFunc("/admin/articles/{id}", s.AdminDeleteArticleHandler).Methods("DELETE")
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

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request)           {}
func (s *Server) AboutHandler(w http.ResponseWriter, r *http.Request)          {}
func (s *Server) GetAllArticlesHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) GetArticleHandler(w http.ResponseWriter, r *http.Request)     {}

func (s *Server) AdminIndexHandler(w http.ResponseWriter, r *http.Request)          {}
func (s *Server) AdminGetAllArticlesHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) AdminCreateArticleHandler(w http.ResponseWriter, r *http.Request)  {}
func (s *Server) AdminGetArticleHandler(w http.ResponseWriter, r *http.Request)     {}
func (s *Server) AdminEditArticleHandler(w http.ResponseWriter, r *http.Request)    {}
func (s *Server) AdminDeleteArticleHandler(w http.ResponseWriter, r *http.Request)  {}
