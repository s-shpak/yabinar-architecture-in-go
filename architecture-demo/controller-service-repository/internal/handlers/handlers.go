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

func Serve(cfg config.Config, foobar Foobar) error {
	h := newHandlers(foobar)
	router := newRouter(h)

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router,
	}

	return srv.ListenAndServe()
}

func newRouter(h *handlers) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /fizzbuzz/{n}", h.GetFoobar)

	return mux
}

type Foobar interface {
	GetFoobar(req *service.GetFoobarRequest) (*service.GetFoobarResponse, error)
}

type handlers struct {
	foobar Foobar
}

func newHandlers(foobar Foobar) *handlers {
	return &handlers{
		foobar: foobar,
	}
}

type GetFoobarResponse struct {
	Data []string `json:"data"`
}

func (h *handlers) GetFoobar(w http.ResponseWriter, r *http.Request) {
	nStr := r.PathValue("n")
	n, err := strconv.ParseInt(nStr, 10, 32)
	if err != nil {
		http.Error(w, "n must be an integer", http.StatusBadRequest)
		return
	}

	resp, err := h.foobar.GetFoobar(&service.GetFoobarRequest{
		N: int(n),
	})
	if err != nil {
		if errors.Is(err, service.ErrGetFoobarInvalidRequest) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("failed to get foobar: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to write a get foobar response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
