package main

import (
	"encoding/json"
	"fmt"
	envoy_config_bootstrap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"log"
	"os"
)

func run() error {
	bootstrap := envoy_config_bootstrap_v3.Bootstrap{
		Node: &envoy_config_core_v3.Node{
			Cluster: "test-cluster",
			Id:      "test-id",
		},
		DynamicResources: &envoy_config_bootstrap_v3.Bootstrap_DynamicResources{
			CdsConfig: &envoy_config_core_v3.ConfigSource{
				ResourceApiVersion: envoy_config_core_v3.ApiVersion_V3,
				ConfigSourceSpecifier: &envoy_config_core_v3.ConfigSource_Path{
					Path: "/etc/envoy/cds.yaml",
				},
			},
		},
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	if err := e.Encode(&bootstrap); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("error: %s", err)
	}
}
