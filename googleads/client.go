package googleads

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

type Client struct {
	conn *grpc.ClientConn
	ctx  *context.Context
}

type (
	ContextOption func(*context.Context)
	HeaderOption  func(header *http.Header)
)

var conn *grpc.ClientConn
var ctx *context.Context

func NewGoogleAdsClient(dev_token, access_token, customer_id string) (*Client, error) {
	conn_, err := GetGRPCConnection()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	headers := metadata.Pairs(
		"authorization", "Bearer "+access_token,
		"developer-token", dev_token,
		"login-customer-id", customer_id,
	)
	ctx = metadata.NewOutgoingContext(ctx, headers)

	c := Client{
		conn: conn_,
		ctx:  &ctx,
	}

	return &c, nil
}

func GetGRPCConnection() (*grpc.ClientConn, error) {
	if conn != nil {
		return conn, nil
	}

	cred := grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	conn, err := grpc.Dial(EndPoint, cred)

	return conn, err
}
