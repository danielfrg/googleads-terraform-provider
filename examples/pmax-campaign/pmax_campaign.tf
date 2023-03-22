# resource "googleads_pmax_campaign" "my_campaign" {
#   name = "My PMax Campaign from TF"

#   headlines = [
#     googleads_text_asset.headline1.resource_name,
#     googleads_text_asset.headline2.resource_name,
#     googleads_text_asset.headline3.resource_name,
#   ]

#   long_headlines = [googleads_text_asset.long_headline1.resource_name]

#   descriptions = [
#     googleads_text_asset.description1.resource_name,
#   ]

#   business_name = googleads_text_asset.business_name.resource_name

#   marketing_images = [
#     googleads_image_asset.marketing_image.resource_name
#   ]

#   logo_images = [
#     googleads_image_asset.logo_image.resource_name
#   ]
# }

resource "googleads_budget" "my_budget" {
  name              = "My PMax Campaign Budget"
  amount_micros     = 10000
  delivery_method   = "STANDARD"
  explicitly_shared = false
}

resource "googleads_text_asset" "headline1" {
  text = "Headline 1"
}

resource "googleads_text_asset" "headline2" {
  text = "Headline 2"
}

resource "googleads_text_asset" "headline3" {
  text = "Headline 3"
}

resource "googleads_text_asset" "long_headline1" {
  text = "Long Headline"
}

resource "googleads_text_asset" "description1" {
  text = "My description"
}

resource "googleads_text_asset" "business_name" {
  text = "My Business"
}

resource "googleads_image_asset" "marketing_image" {
  name = "tf test marketing image asset"
  path = "marketing1.jpeg"
}

resource "googleads_image_asset" "logo_image" {
  name = "tf test image logo asset"
  path = "logo.png"
}

