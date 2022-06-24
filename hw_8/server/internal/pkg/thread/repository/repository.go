package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"server/internal/pkg/domain"
	"server/internal/pkg/http_client"
	"server/tools"
)

type repository struct {
	httpClient *http_client.HttpClient
}

func NewRepository(httpCli *http_client.HttpClient) domain.ThreadRepository {
	return repository{
		httpClient: httpCli,
	}
}

func (r repository) Create(ctx context.Context, thread domain.Thread) error {
	requestID, err := tools.GetRequestID(ctx)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(thread)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "http://vk-golang.ru:15000/thread", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.DoRequestToVk(req)
	if err != nil {
		return err
	}
	log.Printf("%s %v Request: %v %v %v; Response: %v %v", requestID, err, req.Method,
		req.URL, req.Body, resp.Status, resp.Body)

	if resp.StatusCode != 200 {
		return errors.New("failed to create thread remotely")
	}

	return nil
}

func (r repository) Get(ctx context.Context, id string) (domain.Thread, error) {
	requestID, err := tools.GetRequestID(ctx)
	if err != nil {
		return domain.Thread{}, err
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://vk-golang.ru:15000/thread?id=%s", id), nil)
	if err != nil {
		return domain.Thread{}, err
	}

	resp, err := r.httpClient.DoRequestToVk(req)
	if err != nil {
		return domain.Thread{}, err
	}
	log.Printf("%s %v Request: %v %v %v; Response: %v %v", requestID, err, req.Method,
		req.URL, req.Body, resp.Status, resp.Body)

	if resp.StatusCode != 200 {
		return domain.Thread{}, errors.New("failed to fetch thread remotely")
	}

	var thread domain.Thread
	err = json.NewDecoder(resp.Body).Decode(&thread)
	if err != nil {
		return domain.Thread{}, err
	}

	return thread, nil
}
