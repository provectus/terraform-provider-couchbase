package couchbase

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/couchbase/gocb.v1"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const (
	nameProperty          = "name"
	passwordProperty      = "password"
	flushEnabledProperty  = "flush_enabled"
	indexReplicasProperty = "index_replicas"
	quotaProperty         = "quota"
	replicasProperty      = "replicas"
	typeProperty          = "type"
)

func resourceBucket() *schema.Resource {
	return &schema.Resource{
		Create: CreateBucket,
		Read:   ReadBucket,
		Update: UpdateBucket,
		Delete: DeleteBucket,
		Schema: map[string]*schema.Schema{
			nameProperty: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			passwordProperty: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			flushEnabledProperty: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			indexReplicasProperty: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			quotaProperty: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
			replicasProperty: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			typeProperty: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validateType,
			},
		},
	}
}

func CreateBucket(data *schema.ResourceData, meta interface{}) (err error) {
	_, manager, err := connect(meta.(*Config))
	if err != nil {
		return err
	}

	settings := fetchSettings(data)
	if err = manager.InsertBucket(&settings); err != nil {
		return
	}

	log.Printf("[INFO] A bucket with the name %q was created", settings.Name)

	data.SetId(settings.Name)

	return ReadBucketDelayed(data, meta, true)
}

func ReadBucket(data *schema.ResourceData, meta interface{}) (err error) {
	return ReadBucketDelayed(data, meta, false)
}

func ReadBucketDelayed(data *schema.ResourceData, meta interface{}, delayed bool) (err error) {
	log.Printf("[INFO] Reading bucket with the name %q...", data.Id())
	cluster, _, err := connect(meta.(*Config))
	if err != nil {
		return
	}

	if delayed {
		time.Sleep(time.Duration(meta.(*Config).BucketCreationDelay) * time.Second)
	}

	bucket, err := cluster.OpenBucket(data.Id(), data.Get(passwordProperty).(string))
	if err != nil {
		log.Printf("[ERROR] Can not read a bucket %q: %s", data.Id(), err)
		return
	}

	if bucket == nil {
		log.Printf("[WARN] Can not find a bucket %q", data.Id())
		data.SetId("")
	}

	return
}

func UpdateBucket(data *schema.ResourceData, meta interface{}) (err error) {
	bucketSettings := fetchSettings(data)

	if err = updateBucket(&bucketSettings, meta.(*Config)); err != nil {
		return
	}

	log.Printf("[INFO] A bucket with the name %q was updated", data.Id())

	return ReadBucketDelayed(data, meta, true)
}

func DeleteBucket(data *schema.ResourceData, meta interface{}) (err error) {
	_, manager, err := connect(meta.(*Config))
	if err != nil {
		return
	}

	err = manager.RemoveBucket(data.Get(nameProperty).(string))
	if err == nil {
		log.Printf("[INFO] A bucket with the name %q was removed", data.Get(nameProperty).(string))
		data.SetId("")
	}

	return
}

func connect(config *Config) (cluster *gocb.Cluster, manager *gocb.ClusterManager, err error) {
	if cluster, err = gocb.Connect(config.Url); err != nil {
		return
	}

	if err = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: config.Username,
		Password: config.Password,
	}); err != nil {
		return
	}
	manager = cluster.Manager(config.Username, config.Password)

	return
}

func fetchSettings(data *schema.ResourceData) gocb.BucketSettings {
	return gocb.BucketSettings{
		Name:          data.Get(nameProperty).(string),
		Type:          gocb.BucketType(data.Get(typeProperty).(int)),
		Quota:         data.Get(quotaProperty).(int),
		Replicas:      data.Get(replicasProperty).(int),
		FlushEnabled:  data.Get(flushEnabledProperty).(bool),
		IndexReplicas: data.Get(indexReplicasProperty).(bool),
	}
}

func validateType(val interface{}, key string) (warns []string, errs []error) {
	value := val.(int)
	switch value {
	case 0, 1, 2:
		break
	default:
		errs = append(errs, fmt.Errorf("%q contains an invalid value: %v. Valid values are: "+
			"0 (Couchbase), 1 (Memcached), 2 (Ephemeral)", key, value))
	}

	return
}

func updateBucket(settings *gocb.BucketSettings, config *Config) (err error) {

	resp, err := updateRequest(config, settings)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if err = resp.Body.Close(); err != nil {
			log.Printf("[ERROR] Failed to close socket (%s)", err)
		}

		return errors.New(string(data))
	}

	return nil
}

func getData(settings *gocb.BucketSettings) []byte {
	posts := url.Values{}

	switch settings.Type {
	case gocb.Couchbase:
		posts.Add("bucketType", "couchbase")
	case gocb.Memcached:
		posts.Add("bucketType", "memcached")
	case gocb.Ephemeral:
		posts.Add("bucketType", "ephemeral")
	}

	if settings.FlushEnabled {
		posts.Add("flushEnabled", "1")
	} else {
		posts.Add("flushEnabled", "0")
	}

	posts.Add("replicaNumber", fmt.Sprintf("%d", settings.Replicas))
	posts.Add("authType", "sasl")
	posts.Add("saslPassword", settings.Password)
	posts.Add("ramQuotaMB", fmt.Sprintf("%d", settings.Quota))

	return []byte(posts.Encode())
}

func updateRequest(config *Config, settings *gocb.BucketSettings) (*http.Response, error) {
	reqUri := fmt.Sprintf("%s/pools/default/buckets/%s", config.Hosts[rand.Intn(len(config.Hosts))], settings.Name)
	req, err := http.NewRequest("POST", reqUri, bytes.NewReader(getData(settings)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.Username, config.Password)

	return config.HttpClient.Do(req)
}
