# Variables definition
variable "url" {
  type        = string
  description = "The URL (connection string) of Couchbase server"
}

variable "username" {
  type        = string
  description = "A Couchbase user's name"
}

variable "password" {
  type        = string
  description = "A Couchbase user's password"
}

variable "bucket_creation_delay" {
  type        = number
  description = "A delay (in seconds) until the bucket is created on a cluster"
}

# A Couchbase provider
provider "couchbase" {
  url = var.url
  username = var.username
  password = var.password
  bucket_creation_delay = var.bucket_creation_delay
}

# Create a bucket
resource "couchbase_bucket" "bucket" {
  name = "test"
  flush_enabled = false
  quota = 150
  type = 0
}

# Create a primary index
resource "couchbase_index" "index" {
  bucket_name = couchbase_bucket.bucket.name
  index_name = "def_${couchbase_bucket.bucket.name}_primary"
}
