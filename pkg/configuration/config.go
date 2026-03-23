package configuration

import (
	"fmt"
	"os"

	"github.com/vmware-tanzu/velero/pkg/util/boolptr"
	"gopkg.in/yaml.v2"

	"github.com/slackhq/simple-kubernetes-webhook/pkg/mutation"
)

type Configuration struct {
	Mutation *mutation.Mutation `json:"mutation,omitempty"`
}

var sapccOperator = "sapcc-operator-*"
var cleanupCluster = "cleanup-cluster-*"

var DefaultConfig = Configuration{
	Mutation: &mutation.Mutation{
		Pods: []mutation.Pod{
			{
				Name:    &sapccOperator,
				Enabled: boolptr.True(),
				ContainerReplacements: map[string]string{
					"sapcc-operator": "columbus.common.repositories.cloud.sap/com.sap.edgelm/sapcc-operator:0.59.0-test.i330716-latest",
				},
			},
			{
				Name:    &cleanupCluster,
				Enabled: boolptr.True(),
				ContainerReplacements: map[string]string{
					"cleanup": "columbus.common.repositories.cloud.sap/com.sap.edgelm/docker-cleanup-cluster:0.35.0-test.i330716-latest",
				},
			},
		},
	},
}

// ReadConfigurationFromFile reads a YAML file from the given path and unmarshals it into a Configuration struct
func ReadConfigurationFromFile(filePath string) (*Configuration, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	cfg := new(Configuration)
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return cfg, nil
}
