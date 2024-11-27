package repository

import (
	"errors"
	"fmt"
	"sync"
)

type Repository interface {
	GetFoobar(req *GetFoobarRequest) (*GetFoobarResponse, error)
	SetFoobar(req *SetFoobarRequest)
}

type Store struct {
	mux *sync.Mutex
	s   map[int][]string
}

func NewStore() *Store {
	return &Store{
		mux: &sync.Mutex{},
		s:   make(map[int][]string),
	}
}

type GetFoobarRequest struct {
	N int
}

type GetFoobarResponse struct {
	Data []string
}

var (
	ErrGetFoobarNotFound = errors.New("data not found")
)

func newErrGetFoobarNotFound(n int) error {
	return fmt.Errorf("%w for n = %d", ErrGetFoobarNotFound, n)
}

func (s *Store) GetFoobar(req *GetFoobarRequest) (*GetFoobarResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	res, ok := s.s[req.N]
	if !ok {
		return nil, newErrGetFoobarNotFound(req.N)
	}
	return &GetFoobarResponse{
		Data: res,
	}, nil
}

type SetFoobarRequest struct {
	N    int
	Data []string
}

func (s *Store) SetFoobar(req *SetFoobarRequest) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.s[req.N] = req.Data
}
