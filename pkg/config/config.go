package config

import (
	"fmt"

	"github.com/cloudquery/cloudquery/pkg/policy"

	"github.com/cloudquery/cloudquery/internal/logging"
	"github.com/cloudquery/cloudquery/pkg/client/history"
	"github.com/hashicorp/hcl/v2"
)

type Providers []*Provider

func (pp Providers) Names() []string {
	pNames := make([]string, len(pp))
	for i, p := range pp {
		pNames[i] = p.Name
	}
	return pNames
}

type Config struct {
	CloudQuery CloudQuery       `hcl:"cloudquery,block"`
	Providers  Providers        `hcl:"provider,block"`
	Policies   []*policy.Policy `hcl:"policy,block"`
	Modules    hcl.Body         `hcl:"modules,block"`
}

func (c Config) GetProvider(name string) (*Provider, error) {
	for _, p := range c.Providers {
		if name == p.Alias {
			return p, nil
		}
	}
	return nil, fmt.Errorf("provider %s does not exist", name)
}

type CloudQuery struct {
	PluginDirectory string              `hcl:"plugin_directory,optional"`
	PolicyDirectory string              `hcl:"policy_directory,optional"`
	Logger          *logging.Config     `hcl:"logging,block"`
	Providers       []*RequiredProvider `hcl:"provider,block"`
	Connection      *Connection         `hcl:"connection,block"`
	History         *history.Config     `hcl:"history,block"`
}

func (c CloudQuery) GetRequiredProvider(name string) (*RequiredProvider, error) {
	for _, p := range c.Providers {
		if name == p.Name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("provider %s does not exist", name)
}

type Connection struct {
	DSN string `hcl:"dsn,attr"`
}

type RequiredProvider struct {
	Name    string  `hcl:"name,label"`
	Source  *string `hcl:"source,optional"`
	Version string  `hcl:"version"`
}

func (r RequiredProvider) String() string {
	var source string
	if r.Source != nil {
		source = *r.Source
	}
	return fmt.Sprintf("%s/cq-provider-%s@%s", source, r.Name, r.Version)
}

// configFileSchema is the schema for the top-level of a config file. We use
// the low-level HCL API for this level so we can easily deal with each
// block type separately with its own decoding logic.
var configFileSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "cloudquery",
		},
		{
			Type:       "provider",
			LabelNames: []string{"name"},
		},
		{
			Type:       "policy",
			LabelNames: []string{"name"},
		},
		{
			Type: "modules",
		},
	},
}
