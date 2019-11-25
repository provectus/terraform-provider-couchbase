package couchbase

import (
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/couchbase/gocb.v1"
)

const (
	readIndexDelay         = "3s"
	bucketNameProperty     = "bucket_name"
	bucketPasswordProperty = "bucket_password"
	indexNameProperty      = "index_name"
	indexFieldsProperty    = "index_fields"
)

func resourceIndex() *schema.Resource {
	return &schema.Resource{
		Create: createIndex,
		Read:   readIndex,
		Update: updateIndex,
		Delete: deleteIndex,
		Schema: map[string]*schema.Schema{
			bucketNameProperty: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			bucketPasswordProperty: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			indexNameProperty: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			indexFieldsProperty: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func createIndex(data *schema.ResourceData, meta interface{}) (err error) {
	manager, err := getBucketManager(data, meta)
	if err != nil {
		return
	}
	err = doCreateIndex(data, manager)
	if err != nil {
		return
	}
	return doReadIndex(data, manager, readIndexDelay)
}

func readIndex(data *schema.ResourceData, meta interface{}) (err error) {
	manager, err := getBucketManager(data, meta)
	if err != nil {
		return
	}
	doReadIndex(data, manager, readIndexDelay)
	return
}

func updateIndex(data *schema.ResourceData, meta interface{}) (err error) {
	manager, err := getBucketManager(data, meta)
	if err != nil {
		return
	}
	err = doDropIndex(data, manager)
	if err != nil {
		return
	}
	err = doCreateIndex(data, manager)
	if err != nil {
		return
	}
	return doReadIndex(data, manager, readIndexDelay)
}

func deleteIndex(data *schema.ResourceData, meta interface{}) (err error) {
	manager, err := getBucketManager(data, meta)
	if err != nil {
		return
	}
	return doDropIndex(data, manager)
}

func getBucketManager(data *schema.ResourceData, meta interface{}) (manager *gocb.BucketManager, err error) {
	bucketName := data.Get(bucketNameProperty).(string)
	log.Printf("[INFO] Reading bucket with the name %q...", bucketName)
	config := meta.(*Config)
	cluster, _, err := connect(config)
	if err != nil {
		return
	}
	bucket, err := cluster.OpenBucket(bucketName, data.Get(bucketPasswordProperty).(string))
	if err != nil {
		log.Printf("[ERROR] Can not read a bucket %q: %s", bucketName, err)
		return
	}
	if bucket == nil {
		log.Printf("[WARN] Can not find a bucket %q", bucketName)
	}
	manager = bucket.Manager(config.Username, config.Password)
	if manager == nil {
		log.Printf("[WARN] Can not get a manager for bucket %q", bucketName)
	}
	return
}

func doReadIndex(data *schema.ResourceData, manager *gocb.BucketManager, delay string) (err error) {
	d, _ := time.ParseDuration(delay)
	indexName := data.Get(indexNameProperty).(string)
	err = manager.WatchIndexes([]string{indexName}, false, d)
	if err != nil {
		log.Printf("[WARN] Can not find index %q", indexName)
		data.SetId("")
	}
	return
}

func doDropIndex(data *schema.ResourceData, manager *gocb.BucketManager) (err error) {
	indexName := data.Get(indexNameProperty).(string)
	err = manager.DropIndex(indexName, false)
	if err == nil {
		log.Printf("[INFO] An index with the name %q was removed", indexName)
		data.SetId("")
	}
	return
}

func doCreateIndex(data *schema.ResourceData, manager *gocb.BucketManager) (err error) {
	indexName := data.Get(indexNameProperty).(string)
	indexFields := data.Get(indexFieldsProperty).(string)
	if indexFields == "" {
		err = manager.CreatePrimaryIndex(indexName, true, false)
	} else {
		err = manager.CreateIndex(indexName, strings.Split(indexFields, ","), true, false)
	}
	if err == nil {
		log.Printf("[INFO] An index with the name %q was created", indexName)
		data.SetId(indexName)
	}
	return
}
