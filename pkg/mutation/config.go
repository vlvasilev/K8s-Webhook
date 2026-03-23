package mutation

type Mutation struct {
	Pods []Pod `json:"pods,omitempty"`
}

type Pod struct {
	Enabled               *bool             `json:"enabled,omitempty"`
	Name                  *string           `json:"name,omitempty"`
	Annotations           map[string]string `json:"annotations,omitempty"`
	Labels                map[string]string `json:"labels,omitempty"`
	ContainerReplacements map[string]string `json:"containerReplacements,omitempty"`
}

var sapccOperator = "sapcc-operator-*"
var cleanupCluster = "cleanup-cluster-*"
