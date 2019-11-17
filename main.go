package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/provectus/terraform-provider-couchbase/couchbase"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: couchbase.Provider})
}
