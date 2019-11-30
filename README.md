# Terraform Couchbase provider
This provider helps to manage the Couchbase resources. Currently supports the next resources:
- [Bucket](#Bucket)
- [Index](#Index)

## Quick start
Download a suitable binary from the release page and put it into your `~/.terraform.d/plugins/`. After that, you can configure this provider in your terraform code, for the example of the code please follow to `./example` folder.
``` hcl
provider "couchbase" {
  url = "couchbase://localhost:11210"
  username = "Administrator"
  password = "password"
  bucket_creation_delay = 5
}
```

## Development Guide
If you wish to work on the provider, you'll first need Go installed on your machine. You'll also need to install Docker and Docker-compose for testing provider locally.

To compile the provider, run build. This will build the provider and put the provider binary in the current working directory.
``` bash
$ go build
$ ls ./terraform-provider-couchbase
```
Additionally, you can start docker-compose stack, it will build provider binary and put it into your Terraform plugins directory. Also, Couchbase service would be started. 
``` bash
$ docker-compose up -d
$ ls ~/.terraform.d/plugins/terraform-provider-couchbase
```
## Resources
### Bucket

| Property | Type | Description | Default |
|----------|------------|----------------|-----------|
| **name** | `string` | a bucket's name |  |
| **password** | `string` | a bucket's password | `""` |
| **flush_enabled** | `bool` | is flush enabled for the bucket | `false` |
| **index_replicas**  | `bool` | should index be replicated on replicas | `false` |
| **quota** | `integer` | a memory quota for the bucket in megabytes |`100` | 
| **replicas** | `integer` | replicas count for the bucket| `1`|
| **type** | `integer` | a bucket type *(Couchbase (0), Memcached(1), Ephemeral(2))* | `0`|

### Index
| Property | Type | Description | Default |
|----------|------------|----------------|-----------|
| **bucket_name** | `string` | a bucket's name | |
| **bucket_password** | `string` | a bucket's password | `""` |
| **index_name** | `string` | an index name | |
| **index_fields** | `string` | bucket's fields for index (e.g. *"field1, field2"*). _Keep empty to create a primary index_| `""` |

### User
| Property | Type | Description | Default |
|----------|------------|----------------|-----------|
| **user_name** | `string` | a user's name | |
| **bucket_password** | `string` | a user's password | |
| **user_roles** | `string` | a user's permissions for bucket (e.g. *"data_reader:test, data_writer:test"*) | |
