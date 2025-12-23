package service

import (
	"errors"
	"fmt"
	"strconv"

	"csr/internal/repository"
)

type Repository interface {
	GetFizzBuzz(req *repository.GetFizzBuzzRequest) (*repository.GetFizzBuzzResponse, error)
	SetFizzBuzz(req *repository.SetFizzBuzzRequest)
}

type FizzBuzz struct {
	store Repository
}

func NewFizzBuzz(store Repository) *FizzBuzz {
	return &FizzBuzz{
		store: store,
	}
}

type GetFizzBuzzRequest struct {
	N int
}

type GetFizzBuzzResponse struct {
	Data []string
}

var (
	ErrGetFizzBuzzInvalidRequest = errors.New("invalid get fizzbuzz request")
	ErrRepoFailed                = errors.New("repo failed")
)

func (f *FizzBuzz) GetFizzBuzz(req *GetFizzBuzzRequest) (*GetFizzBuzzResponse, error) {
	if err := getFizzBuzzValidateRequest(req); err != nil {
		return nil, err
	}

	repositoryResp, err := f.store.GetFizzBuzz(&repository.GetFizzBuzzRequest{
		N: req.N,
	})
	if err != nil {
		if !errors.Is(err, repository.ErrGetFizzBuzzNotFound) {
			return nil, fmt.Errorf("failed to fetch the fizzbuzz result from the store: %w", err)
		}
	}

	if repositoryResp != nil {
		return &GetFizzBuzzResponse{
			Data: repositoryResp.Data,
		}, nil
	}

	resp := f.calculateFizzBuzz(req)
	f.store.SetFizzBuzz(&repository.SetFizzBuzzRequest{
		N:    req.N,
		Data: resp.Data,
	})

	return resp, nil
}

func getFizzBuzzValidateRequest(req *GetFizzBuzzRequest) error {
	if req.N <= 0 {
		return fmt.Errorf("%w: n must be greater than 0, got %d", ErrGetFizzBuzzInvalidRequest, req.N)
	}
	nLimit := 10000
	if req.N > nLimit {
		return fmt.Errorf("%w: n must be less than %d, got %d", ErrGetFizzBuzzInvalidRequest, nLimit, req.N)
	}

	// эта проверка синтетическая, она имеет смысл только для демо
	if req.N == 42 {
		return fmt.Errorf("%d is a no-no", req.N)
	}
	return nil
}

func (f *FizzBuzz) calculateFizzBuzz(req *GetFizzBuzzRequest) *GetFizzBuzzResponse {
	resp := &GetFizzBuzzResponse{
		Data: make([]string, 0, req.N),
	}
	for i := 1; i <= req.N; i++ {
		var res string
		if i%3 == 0 && i%5 == 0 {
			res = "fizzbuzz"
		} else if i%3 == 0 {
			res = "fizz"
		} else if i%5 == 0 {
			res = "buzz"
		} else {
			res = strconv.Itoa(i)
		}
		resp.Data = append(resp.Data, res)
	}
	return resp
}

// Cyclic dependency
// var (
// 	ErrGetFizzBuzzNotFound = errors.New("data not found")
// )
