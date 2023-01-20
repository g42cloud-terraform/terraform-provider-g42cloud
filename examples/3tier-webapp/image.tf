# Feching the images with the specificed filters
data "g42cloud_images_image" "image" {
  name = "CentOS 7.9 64bit"
  visibility = "public"
  most_recent = true
}
