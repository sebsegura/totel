package service

import (
	"context"
	"encoding/json"
	"github.com/Bancar/goala/ulog"
	"github.com/Bancar/goala/utel"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"io"
	"log"
	"net/http"
	ut "totel/utel"
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
	mp := otel.GetMeterProvider()
	cfg := ut.GetUtelConfig()

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
	err = utel.IncrementCounter(ctx, mp.(*sdkmetric.MeterProvider), &utel.MetricAttributes{
		Name:  "fn2-metric",
		Unit:  "1",
		Flow:  cfg.Flow,
		Owner: cfg.Owner,
	})
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

type ddb struct {
	client *dynamodb.Client
}

func newConn(ctx context.Context) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("cannot connect to dynamodb: %v", err)
	}

	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointResolverV2 = dynamodb.NewDefaultEndpointResolverV2()
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})
}

func NewDDBClient(ctx context.Context) Client {
	return &ddb{
		client: newConn(ctx),
	}
}

func (c *ddb) Create(ctx context.Context, in *Request) error {
	tracer := otel.Tracer(ut.GetUtelConfig().ServiceName)
	ctx, span := tracer.Start(ctx, "DDB")
	defer span.End()

	span.SetAttributes(
		attribute.String("table.name", "violeros"),
		attribute.String("item.id", in.ID),
		attribute.String("item.name", in.Name),
	)
	item, err := attributevalue.MarshalMap(in)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("violeros"),
		Item:      item,
	}

	span.AddEvent("Saving item...")
	_, err = c.client.PutItem(ctx, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
