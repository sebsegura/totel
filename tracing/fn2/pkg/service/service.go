package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"log"
	"totel/utel"
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

func (s *Service) Handle(ctx context.Context, in *Request) (*Response, error) {
	tracer := otel.Tracer(utel.GetUtelConfig().ServiceName)
	ctx, span := tracer.Start(ctx, "HandleRequest")
	defer span.End()

	if err := s.client.Create(ctx, in); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &Response{Msg: "ok"}, nil
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
	fmt.Println("hey")
	tracer := otel.Tracer(utel.GetUtelConfig().ServiceName)
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
