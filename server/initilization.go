package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

var adminToken = "some_token"

type Server struct {
	db *DB
}

func createServer(db *DB) *Server {
	return &Server{db: db}
}

func (s *Server) makeEndPoints() {
	r := mux.NewRouter()
	r.HandleFunc("/user_banner", s.UserBannerHandler)
	r.HandleFunc("/banner", s.BannerHandler)
	r.HandleFunc("/banner/{id:[0-9]+}", s.IdBannerHandler)
	http.Handle("/", r)
}
