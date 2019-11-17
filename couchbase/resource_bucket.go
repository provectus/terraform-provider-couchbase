package couchbase

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/couchbase/gocb.v1"
)

const (
	nameProperty     = "name"
	passwordProperty = "password"
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
		},
	}
}

func CreateBucket(data *schema.ResourceData, meta interface{}) error {
	manager, err := getManager(meta.(*Config))
	if err != nil {
		return err
	}

	bucketSettings := gocb.BucketSettings{
		FlushEnabled:  true,
		IndexReplicas: true,
		Name:          data.Get(nameProperty).(string),
		Password:      data.Get(passwordProperty).(string),
		Quota:         120,
		Replicas:      1,
		Type:          gocb.Couchbase,
	}

	err = manager.InsertBucket(&bucketSettings)
	if err != nil {
		return err
	}

	data.SetId(data.Get(nameProperty).(string))

	return nil
}

func ReadBucket(data *schema.ResourceData, meta interface{}) error {
	cluster, err := gocb.Connect(meta.(*Config).Url)
	if err != nil {
		return err
	}

	_ = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: meta.(*Config).Username,
		Password: meta.(*Config).Password,
	})

	bucket, err := cluster.OpenBucket(data.Get(nameProperty).(string), data.Get(passwordProperty).(string))
	if err != nil {
		return err
	}

	if bucket == nil {
		data.SetId("")
	}

	return nil
}

func UpdateBucket(data *schema.ResourceData, meta interface{}) error {
	manager, err := getManager(meta.(*Config))
	if err != nil {
		return err
	}

	bucketSettings := gocb.BucketSettings{
		FlushEnabled:  true,
		IndexReplicas: true,
		Name:          data.Get(nameProperty).(string),
		Password:      data.Get(passwordProperty).(string),
		Quota:         120,
		Replicas:      1,
		Type:          gocb.Couchbase,
	}

	err = manager.UpdateBucket(&bucketSettings)
	if err != nil {
		return err
	}

	return nil
}

func DeleteBucket(data *schema.ResourceData, meta interface{}) error {
	manager, err := getManager(meta.(*Config))
	if err != nil {
		return err
	}

	err = manager.RemoveBucket(data.Get(nameProperty).(string))
	if err == nil {
		data.SetId("")
	}

	return err
}

func getManager(config *Config) (*gocb.ClusterManager, error) {
	cluster, err := gocb.Connect(config.Url)
	if err != nil {
		return nil, err
	}

	return cluster.Manager(config.Username, config.Password), nil
}
