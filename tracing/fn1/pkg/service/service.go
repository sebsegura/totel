package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
	"os"
	"strings"
	"totel/utel"
)

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	Msg string `json:"msg"`
}

type Service struct {
	client *http.Client
}

func New(c *http.Client) *Service {
	return &Service{
		client: c,
	}
}

func (s *Service) Handle(ctx context.Context, in *Request) (*Response, error) {
	return s.call(ctx, in)
}

func id() string {
	return ulid.Make().String()
}

func (s *Service) call(ctx context.Context, in *Request) (*Response, error) {
	tracer := otel.Tracer(utel.GetUtelConfig().ServiceName)
	ctx, span := tracer.Start(ctx, "Fn2Request")
	defer span.End()

	type R struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	p := R{
		ID:   id(),
		Name: in.Name,
	}
	b, _ := json.Marshal(p)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint(), strings.NewReader(string(b)))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	propagator := xray.Propagator{}
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	req.Header.Set("Content-Type", "application/json")
	span.AddEvent("Making request...")
	res, err := s.client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var r Response
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

func endpoint() string {
	e := os.Getenv("ENDPOINT")
	if e == "" {
		e = "http://127.0.0.1:9001/lambda"
	}

	return e
}
