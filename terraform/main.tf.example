# 配置阿里云 Provider
provider "alicloud" {
  access_key = "Your AccessKey ID Here"
  secret_key = "Your AccessKey Secret Here"
  region     = "cn-hangzhou"
}

resource "alicloud_oss_bucket" "bucket-acl" {
  bucket = "bucket-laracom-acl"
  acl    = "private"
}