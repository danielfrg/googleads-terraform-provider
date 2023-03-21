package client

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	EndPoint                      = "googleads.googleapis.com:443"
	ListAccessibleCustomersMehtod = "GET"
	GoogleAdsSearchMehtod         = "POST"
)

type GoogleAdsClient struct {
	Connection grpc.ClientConn
	Context    context.Context
	CustomerId string
}

type (
	ContextOption func(*context.Context)
	HeaderOption  func(header *http.Header)
)

// var ctx *context.Context

func NewGoogleAdsClient(dev_token, access_token, customer_id, login_customer_id string) (*GoogleAdsClient, error) {
	ctx := context.Background()

	headers := metadata.Pairs(
		"authorization", "Bearer "+access_token,
		"developer-token", dev_token,
		"login-customer-id", login_customer_id,
	)
	ctx = metadata.NewOutgoingContext(ctx, headers)

	cred := grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	conn, err := grpc.Dial("googleads.googleapis.com:443", cred)
	if err != nil {
		panic(err)
	}

	c := GoogleAdsClient{
		Connection: *conn,
		Context:    ctx,
		CustomerId: customer_id,
	}

	return &c, nil
}
