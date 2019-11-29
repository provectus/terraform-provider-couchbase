package couchbase

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"gopkg.in/couchbase/gocb.v1"
	"gopkg.in/couchbase/gocbcore.v7"
)

const (
	providerUrl                 = "url"
	providerName                = "username"
	providerPassword            = "password"
	providerBucketCreationDelay = "bucket_creation_delay"
)

type Config struct {
	Url                 string
	Username            string
	Password            string
	BucketCreationDelay int
	AgentConfig         gocbcore.AgentConfig
	HttpClient          http.Client
	Hosts               []string
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			providerUrl: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COUCHBASE_URL", nil),
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "" {
						errors = append(errors, fmt.Errorf("url must not be an empty string"))
					}
					return
				},
				Description: "The URL (connection string) of Couchbase server",
			},
			providerName: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COUCHBASE_USERNAME", nil),
				Description: "A Couchbase user's name",
			},

			providerPassword: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("COUCHBASE_PASSWORD", nil),
				Description: "A Couchbase user's password",
			},
			providerBucketCreationDelay: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("COUCHBASE_PASSWORD", nil),
				Description: "A delay (in seconds) until the bucket is created on a cluster",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"couchbase_bucket": resourceBucket(),
			"couchbase_index":  resourceIndex(),
			"couchbase_user":   resourceUser(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	agentConfig := gocbcore.AgentConfig{
		UserString:           "gocb/" + gocb.Version(),
		ConnectTimeout:       60000 * time.Millisecond,
		ServerConnectTimeout: 7000 * time.Millisecond,
		NmvRetryDelay:        100 * time.Millisecond,
		UseKvErrorMaps:       true,
		UseDurations:         true,
		NoRootTraceSpans:     true,
		UseCompression:       true,
		UseZombieLogger:      true,
	}

	if err := agentConfig.FromConnStr(d.Get(providerUrl).(string)); err != nil {
		return nil, err
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: agentConfig.TlsConfig,
		},
	}

	delay, err := strconv.Atoi(d.Get(providerBucketCreationDelay).(string))
	if err != nil {
		return nil, err
	}

	return &Config{
		Url:                 d.Get(providerUrl).(string),
		Username:            d.Get(providerName).(string),
		Password:            d.Get(providerPassword).(string),
		BucketCreationDelay: delay,
		AgentConfig:         agentConfig,
		HttpClient:          httpClient,
		Hosts:               getHosts(&agentConfig),
	}, nil
}

func getHosts(agentConfig *gocbcore.AgentConfig) (hosts []string) {
	for _, host := range agentConfig.HttpAddrs {
		if agentConfig.TlsConfig != nil {
			hosts = append(hosts, "https://"+host)
		} else {
			hosts = append(hosts, "http://"+host)
		}
	}
	return
}
