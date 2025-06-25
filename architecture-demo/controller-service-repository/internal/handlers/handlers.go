package handlers

import (
	"csr/internal/handlers/config"
	"csr/internal/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

func Serve(cfg config.Config, fizzBuzz FizzBuzz) error {
	h := newHandlers(fizzBuzz)
	router := newRouter(h)

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router,
	}

	return srv.ListenAndServe()
}

func newRouter(h *handlers) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /fizzbuzz/{n}", h.GetFizzBuzz)

	return mux
}

type FizzBuzz interface {
	GetFizzBuzz(req *service.GetFizzBuzzRequest) (*service.GetFizzBuzzResponse, error)
}

type handlers struct {
	fizzBuzz FizzBuzz
}

func newHandlers(fizzBuzz FizzBuzz) *handlers {
	return &handlers{
		fizzBuzz: fizzBuzz,
	}
}

type GetFizzBuzzResponse struct {
	Data []string `json:"data"`
}

func (h *handlers) GetFizzBuzz(w http.ResponseWriter, r *http.Request) {
	nStr := r.PathValue("n")
	n, err := strconv.ParseInt(nStr, 10, 32)
	if err != nil {
		http.Error(w, "n must be an integer", http.StatusBadRequest)
		return
	}

	resp, err := h.fizzBuzz.GetFizzBuzz(&service.GetFizzBuzzRequest{
		N: int(n),
	})
	if err != nil {
		if errors.Is(err, service.ErrGetFizzBuzzInvalidRequest) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("failed to get fizzbuzz: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to write a get fizzbuzz response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
