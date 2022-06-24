package http_client

import (
	"net/http"
	"server/internal/api/middleware"
	"time"
)

type HttpClient struct {
	monitoring *middleware.Monitoring
}

func NewHttpClient(mnt *middleware.Monitoring) *HttpClient {
	return &HttpClient{
		monitoring: mnt,
	}
}

func (c *HttpClient) DoRequestToVk(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	duration := time.Since(start)

	method := req.Method
	URL := req.URL.String()
	status := resp.Status

	c.monitoring.HitsVk.WithLabelValues(method, URL, status).Inc()
	c.monitoring.DurationVk.WithLabelValues(method, URL, status).Observe(duration.Seconds())
	return resp, err
}
