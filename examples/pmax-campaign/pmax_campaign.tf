
# resource "googleads_pmax_campaign" "my_campaign" {
#   name = "My PMax Campaign from TF"

#   headlines = [
#     googleads_text_asset.headline1,
#     googleads_text_asset.headline2,
#     googleads_text_asset.headline3,
#   ]

#   long_headlines = ["Long Headline"]

#   descriptions = [
#     googleads_text_asset.long_headline1,
#   ]

#   business_name = googleads_text_asset.business_name

#   marketing_images = [
#     googleads_image_asset.marketing_image.resource_name
#   ]

#   logo_images = [
#     googleads_image_asset.logo_image.resource_name
#   ]
# }

resource "googleads_text_asset" "headline1" {
  text = "Headline 1 123"
}

# resource "googleads_text_asset" "headline2" {
#   text = "Headline 2"
# }

# resource "googleads_text_asset" "headline3" {
#   text = "Headline 3"
# }

# resource "googleads_text_asset" "long_headline1" {
#   text = "Long Headline"
# }

# resource "googleads_text_asset" "business_name" {
#   text = "My Business"
# }

# resource "googleads_image_asset" "marketing_image" {
#   name = "tf test image asset"
#   path = "marketing1.jpeg"
# }

resource "googleads_image_asset" "logo_image" {
  name = "tf test image logo asset"
  path = "logo.png"
}

