package server

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Page struct {
	Title   string
	Slug    string
	Body    []byte
	PubDate time.Time // 2006-01-02 15:04
}

func (p *Page) save() error {
	filename := "data/" + p.Slug + ".md"
	return os.WriteFile(filename, []byte(p.Body), 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

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

		log.Printf("logged in: %v", parts)

		next.ServeHTTP(w, r)
	})
}

// public
func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/articles", http.StatusFound)
}
func (s *Server) GetAllPublishedArticlesHandler(w http.ResponseWriter, r *http.Request) {

}
func (s *Server) GetPublishedArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	w.Write([]byte("slug: " + slug))
}

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

	// Admin subrouter with authentication
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(s.basicAuthentication)

	adminRouter.HandleFunc("/", s.AdminIndexHandler).Methods("GET")
	// with ?action=new
	adminRouter.HandleFunc("/articles", s.AdminGetAllArticlesHandler).Methods("GET")
	adminRouter.HandleFunc("/articles", s.AdminCreateArticleHandler).Methods("POST")
	// with ?action=edit
	adminRouter.HandleFunc("/articles/{slug}", s.AdminUpdateArticleGetHandler).Methods("GET")
	adminRouter.HandleFunc("/articles/{slug}", s.AdminUpdateArticleHandler).Methods("PUT")
	adminRouter.HandleFunc("/articles/{slug}", s.AdminDeleteArticleHandler).Methods("DELETE")
	return r
}
