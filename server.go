package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/georgejdanforth/rckv/kv"
)

type Server struct {
	kvStore kv.KVStore
}

func NewServer(kvStore kv.KVStore) *Server {
	return &Server{kvStore}
}

func (s Server) HandleGet(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query()
	keys, ok := queryString["key"]
	if !ok || len(keys) < 1 {
		http.Error(w, "Key not given", http.StatusBadRequest)
		return
	}
	if len(keys) > 1 {
		http.Error(w, "Multiple values for key given", http.StatusBadRequest)
		return
	}

	key := keys[0]
	if key == "" {
		http.Error(w, "Key must be at least one byte in length", http.StatusBadRequest)
	}

	val, err := s.kvStore.Get(key)
	if err != nil {
		if errors.Is(err, kv.KeyNotFound) {
			http.Error(w, fmt.Sprintf("Key %s not found", key), http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(val))
}

func (s Server) HandleSet(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query()
	if len(queryString) != 1 {
		http.Error(w, "Setting multiple keys is not supported", http.StatusBadRequest)
		return
	}
	var err error
	for key := range queryString {
		vals, _ := queryString[key]
		if len(vals) > 1 {
			http.Error(w, fmt.Sprintf("Multiple values given for key %s", key), http.StatusBadRequest)
			return
		}
		err = s.kvStore.Set(key, vals[0])
		break
	}

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
