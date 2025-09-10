package server

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func unauthorizedResponse(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func (s *Server) basicAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			unauthorizedResponse(w)
			return
		}

		encodedBase64 := authHeader[len("Basic "):]
		payload, err := base64.StdEncoding.DecodeString(encodedBase64)
		if err != nil {
			unauthorizedResponse(w)
			return
		}

		parts := strings.SplitN(string(payload), ":", 2)
		if parts[0] != "admin" || parts[1] != "admin" {
			unauthorizedResponse(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// public
func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/articles", http.StatusFound)
}
func (s *Server) GetAllPublishedArticlesHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) GetPublishedArticleHandler(w http.ResponseWriter, r *http.Request)     {}

// auth
func (s *Server) AdminIndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/articles", http.StatusFound)
}
func (s *Server) AdminGetAllArticlesHandler(w http.ResponseWriter, r *http.Request)   {}
func (s *Server) AdminCreateArticleHandler(w http.ResponseWriter, r *http.Request)    {}
func (s *Server) AdminUpdateArticleGetHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) AdminUpdateArticleHandler(w http.ResponseWriter, r *http.Request)    {}
func (s *Server) AdminDeleteArticleHandler(w http.ResponseWriter, r *http.Request)    {}

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", s.IndexHandler).Methods("GET")
	r.HandleFunc("/articles", s.GetAllPublishedArticlesHandler).Methods("GET")
	r.HandleFunc("/articles/{slug}", s.GetPublishedArticleHandler).Methods("GET")

	r.Use(s.basicAuthentication)
	r.HandleFunc("/admin", s.AdminIndexHandler).Methods("GET")
	// with ?action=new
	r.HandleFunc("/admin/articles", s.AdminGetAllArticlesHandler).Methods("GET")
	r.HandleFunc("/admin/articles", s.AdminCreateArticleHandler).Methods("POST")

	// with ?action=edit
	r.HandleFunc("/admin/articles/{slug}", s.AdminUpdateArticleGetHandler).Methods("GET")
	r.HandleFunc("/admin/articles/{slug}", s.AdminUpdateArticleHandler).Methods("PUT")

	r.HandleFunc("/admin/articles/{slug}", s.AdminDeleteArticleHandler).Methods("DELETE")
	return r
}
