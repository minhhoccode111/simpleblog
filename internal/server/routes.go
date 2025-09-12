package server

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

type Page struct {
	Title   string
	Slug    string
	Body    []byte
	PubDate time.Time
}

const timeFormat = "2006-01-02"

var templates = template.Must(template.ParseGlob("template/*.html"))
var linkRegex = regexp.MustCompile(`\[(.*?)\]`)

func (p *Page) save() error {
	filename := "data/" + p.Slug + ".md"
	var body []byte
	metadata := []string{
		"---",
		p.Title,
		p.PubDate.Format(timeFormat),
		"---",
	}

	body = append(body, []byte(strings.Join(metadata, "\n")+"\n")...)
	body = append(body, p.Body...)

	return os.WriteFile(filename, body, 0600)
}

func loadPage(slug string) (*Page, error) {
	filename := "data/" + slug + ".md"
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var title string
	var pubDate time.Time

	// read the first 4 lines
	for range 4 {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		if line == "---" {
			continue
		}
		if title == "" {
			title = line
			continue
		}

		lineTime, err := time.Parse(timeFormat, line)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		pubDate = lineTime
	}

	rest, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &Page{title, slug, rest, pubDate}, nil
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
