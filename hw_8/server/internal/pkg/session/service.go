package session

import (
	"context"
	"errors"
	"log"
	"net/http"
	"server/internal/pkg/domain"
	"server/internal/pkg/http_client"
	"server/tools"
)

type service struct {
	httpClient *http_client.HttpClient
}

func NewService(httpCli *http_client.HttpClient) domain.SessionService {
	return service{
		httpClient: httpCli,
	}
}

func (s service) CheckSession(ctx context.Context, headers http.Header) (domain.Session, error) {
	requestID, err := tools.GetRequestID(ctx)
	if err != nil {
		return domain.Session{}, err
	}

	req, err := http.NewRequest(http.MethodGet, "http://vk-golang.ru:17000/int/CheckSession", nil)
	if err != nil {
		return domain.Session{}, err
	}

	req.Header = headers

	resp, err := s.httpClient.DoRequestToVk(req)
	if err != nil {
		return domain.Session{}, err
	}
	log.Printf("%s %v Request: %v %v %v; Response: %v %v", requestID, err, req.Method,
		req.URL, req.Body, resp.Status, resp.Body)

	switch resp.StatusCode {
	case 500:
		return domain.Session{}, errors.New("failed to request check session")
	case 200:
		return domain.Session{}, nil
	default:
		return domain.Session{}, domain.ErrNoSession
	}
}
