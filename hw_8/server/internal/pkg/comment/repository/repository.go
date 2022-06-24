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

func NewRepository(httpCli *http_client.HttpClient) domain.CommentRepository {
	return repository{
		httpClient: httpCli,
	}
}

func (r repository) Create(ctx context.Context, comment domain.Comment) error {
	requestID, err := tools.GetRequestID(ctx)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(comment)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "http://vk-golang.ru:16000/comment", bytes.NewBuffer(reqBody))
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
		return errors.New("failed to create comment remotely")
	}

	return nil
}

func (r repository) Like(ctx context.Context, commentID string) error {
	requestID, err := tools.GetRequestID(ctx)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://vk-golang.ru:16000/comment/like?cid=%s", commentID),
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := r.httpClient.DoRequestToVk(req)
	if err != nil {
		return err
	}
	log.Printf("%s %v Request: %v %v %v; Response: %v %v", requestID, err, req.Method,
		req.URL, req.Body, resp.Status, resp.Body)

	if resp.StatusCode != 200 {
		return errors.New("failed to like comment remotely")
	}

	return nil
}
