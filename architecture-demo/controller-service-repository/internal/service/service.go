package service

import (
	"csr/internal/repository"
	"errors"
	"fmt"
	"strconv"
)

type Repository interface {
	GetFoobar(req *repository.GetFoobarRequest) (*repository.GetFoobarResponse, error)
	SetFoobar(req *repository.SetFoobarRequest)
}

type Foobar struct {
	store Repository
}

func NewFoobar(store Repository) *Foobar {
	return &Foobar{
		store: store,
	}
}

type GetFoobarRequest struct {
	N int
}

type GetFoobarResponse struct {
	Data []string
}

var (
	ErrGetFoobarInvalidRequest = errors.New("invalid get foobar request")
	ErrRepoFailed              = errors.New("repo failed")
)

func (f *Foobar) GetFoobar(req *GetFoobarRequest) (*GetFoobarResponse, error) {
	if err := getFoobarValidateRequest(req); err != nil {
		return nil, err
	}

	repositoryResp, err := f.store.GetFoobar(&repository.GetFoobarRequest{
		N: req.N,
	})
	if err != nil {
		if !errors.Is(err, repository.ErrGetFoobarNotFound) {
			return nil, fmt.Errorf("failed to fetch the foobar result from the store: %w", err)
		}
	}

	if repositoryResp != nil {
		return &GetFoobarResponse{
			Data: repositoryResp.Data,
		}, nil
	}

	resp := f.calculateFoobar(req)
	f.store.SetFoobar(&repository.SetFoobarRequest{
		N:    req.N,
		Data: resp.Data,
	})

	return resp, nil
}

func getFoobarValidateRequest(req *GetFoobarRequest) error {
	if req.N <= 0 {
		return fmt.Errorf("%w: n must be greater than 0, got %d", ErrGetFoobarInvalidRequest, req.N)
	}
	nLimit := 10000
	if req.N > nLimit {
		return fmt.Errorf("%w: n must be less than %d, got %d", ErrGetFoobarInvalidRequest, nLimit, req.N)
	}

	// эта проверка синтетическая, она имеет смысл только для демо
	if req.N == 42 {
		return fmt.Errorf("%d is a no-no", req.N)
	}
	return nil
}

func (f *Foobar) calculateFoobar(req *GetFoobarRequest) *GetFoobarResponse {
	resp := &GetFoobarResponse{
		Data: make([]string, 0, req.N),
	}
	for i := 1; i <= req.N; i++ {
		var res string
		if i%3 == 0 && i%5 == 0 {
			res = "foobar"
		} else if i%3 == 0 {
			res = "foo"
		} else if i%5 == 0 {
			res = "bar"
		} else {
			res = strconv.Itoa(i)
		}
		resp.Data = append(resp.Data, res)
	}
	return resp
}
