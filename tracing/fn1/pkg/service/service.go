package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Bancar/goala/utel"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
	"os"
	"strings"
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
	ctx, span := utel.NewSpan(ctx, "Fn2Request")
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
		fmt.Println("error 1")
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
		err = fmt.Errorf("invalid status code: %d", res.StatusCode)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var r Response
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.String("process.status", "ok"))

	return &r, nil
}

func endpoint() string {
	e := os.Getenv("ENDPOINT")
	if e == "" {
		e = "http://127.0.0.1:9001/lambda"
	} else {
		e = "https://t69dj24z8h.execute-api.us-east-1.amazonaws.com/test/otel"
	}

	return e
}
