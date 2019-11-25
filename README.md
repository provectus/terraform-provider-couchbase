# Terraform Couchbase provider
This provider helps to execute CRUD operation over Couchbase buckets.

### Build
- To build a provider's executable file, run: `go build`.

### Usage
- Copy a file `cp terraform-provider-couchbase $HOME/.terraform.d/plugins/`.
- Start Couchbase using docker-compose: `docker-compose up -d`
- Create terraform configuration. See examples: `example/couchbase.*`
- Run `terrafrom init && terraform apply` to see results

### Bucket resource properties
- **name** `string`  a bucket's name (REQUIRED);
- **password** `string` - a bucket's password *(default - "")*;
- **flush_enabled** `bool`- is flush enabled for the bucket *(default - false)*;
- **index_replicas** `bool` - should index be replicated on replicas *(default - false)*;
- **quota** `integer` - a memory quota for the bucket in megabytes *(default - 100)*; 
- **replicas** `integer` - replicas count for the bucket;
- **type** `integer` - a bucket type *(Couchbase (0) - default, Memcached(1), Ephemeral(2))*

### Index resource properties
- **bucket_name** `string` - a bucket's name (REQUIRED);
- **bucket_password** `string` - a bucket's password *(default - "")*;
- **index_name** `string` - an index name (REQUIRED);
- **index_fields** `string` - bucket's fields *(e.g. "field1, field2")* for index. Primary index will be created by default *(default - "")*;