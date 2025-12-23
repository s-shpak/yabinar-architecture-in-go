package repository

import (
	"errors"
	"fmt"
	"sync"
)

type Repository interface {
	GetFizzBuzz(req *GetFizzBuzzRequest) (*GetFizzBuzzResponse, error)
	SetFizzBuzz(req *SetFizzBuzzRequest)
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

type GetFizzBuzzRequest struct {
	N int
}

type GetFizzBuzzResponse struct {
	Data []string
}

var (
	ErrGetFizzBuzzNotFound = errors.New("data not found")
)

func newErrGetFizzBuzzNotFound(n int) error {
	return fmt.Errorf("%w for n = %d", ErrGetFizzBuzzNotFound, n)
}

func (s *Store) GetFizzBuzz(req *GetFizzBuzzRequest) (*GetFizzBuzzResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	res, ok := s.s[req.N]
	if !ok {
		return nil, newErrGetFizzBuzzNotFound(req.N)
		// Cyclic dependecy
		// return nil, fmt.Errorf("fizzbuzz not found: %w", service.ErrGetFizzBuzzNotFound)
	}
	return &GetFizzBuzzResponse{
		Data: res,
	}, nil
}

type SetFizzBuzzRequest struct {
	N    int
	Data []string
}

func (s *Store) SetFizzBuzz(req *SetFizzBuzzRequest) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.s[req.N] = req.Data
}
