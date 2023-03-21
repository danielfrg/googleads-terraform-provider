package main

import (
	"fmt"
	"os"
	"terraform-provider-googleads/googleads/client"
	"terraform-provider-googleads/googleads/resources"

	"github.com/shenzhencenter/google-ads-pb/services"
)

func main() {
	dev_token := os.Getenv("GOOGLEADS_DEVELOPER_TOKEN")
	access_token := os.Getenv("GOOGLEADS_ACCESS_TOKEN")
	customer_id := os.Getenv("GOOGLEADS_CUSTOMER_ID")
	login_customer_id := os.Getenv("GOOGLEADS_LOGIN_CUSTOMER_ID")

	client, err := client.NewGoogleAdsClient(dev_token, access_token, customer_id, login_customer_id)

	if err != nil {
		panic(err)
	}

	adsService := services.NewGoogleAdsServiceClient(&client.Connection)
	_ = adsService

	image, err := resources.GetImageFromFilePath("test/640x360.png")
	if err != nil {
		panic(err)
	}
	fmt.Println(image.MimeType)
	fmt.Println(image.Hash)
}
