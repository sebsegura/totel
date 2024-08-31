package service

import (
	"context"
	"encoding/json"
	"github.com/Bancar/goala/ulog"
	"github.com/Bancar/goala/utel"
	"io"
	"net/http"
)

type Request struct {
	ID   string `dynamodbav:"id" json:"id"`
	Name string `dynamodbav:"name" json:"name"`
}

type Response struct {
	Msg string `json:"msg"`
}

type Service struct {
	client Client
}

func New(c Client) *Service {
	return &Service{
		client: c,
	}
}

func (s *Service) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := utel.NewSpan(ctx, "Handle")
	defer span.End()

	in, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, []byte(`{"error": "bad request"}`))
		return
	}

	var rr Request
	if err = json.Unmarshal(in, &rr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, []byte(`{"error": "internal"}`))
		return
	}

	if err = s.client.Create(ctx, &rr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, []byte(`{"error": "internal"}`))
		return
	}

	res, _ := json.Marshal(&Response{Msg: "ok"})
	err = utel.IncrementCounter(ctx, "fn2-request")
	if err != nil {
		ulog.With(ulog.Str("error", err.Error())).Error("cannot send metric data")
	}

	writeResponse(w, res)
}

func writeResponse(w http.ResponseWriter, body []byte) {
	_, _ = w.Write(body)
}

func (s *Service) Save(ctx context.Context, in *Request) error {
	return s.client.Create(ctx, in)
}
