package main

import (
	"context"
	"terraform-provider-googleads/googleads"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), googleads.New, providerserver.ServeOpts{
		Address: "github.com/danielfrg/googleads-tf-provider",
	})
}
