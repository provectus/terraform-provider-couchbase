package couchbase

import "gopkg.in/couchbase/gocb.v1"

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
