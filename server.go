package main

import (
	"net/http"

	"github.com/georgejdanforth/rckv/kv"
)

type Server struct {
	kvStore *kv.KVStore
}

func NewServer(kvStore *kv.KVStore) *Server {
	return &Server{kvStore}
}

func (s Server) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func (s Server) HandleSet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
