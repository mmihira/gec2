package config

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

// InstanceConfig struct to represent an instance config
type InstanceConfig struct {
	Ami               string   `json:"ami"`
	Type              string   `json:"type"`
	Placement         string   `json:"placement"`
	AttachVolume      string   `json:"attach_volume"`
	Volume            string   `json:"volume"`
	VolumeMountPoint  string   `json:"volume_mount_point"`
	VolumeMountDir    string   `json:"volume_mount_dir"`
	EnvInjection      []string `json:"env_injection"`
	SecurityGroups    []string `json:"security_groups"`
	Roles             []string `json:"roles"`
}

type NodeInst map[string]InstanceConfig
type Config struct {
	Provider string     `json:"provider"`
	Nodes    []NodeInst `json:"nodes"`
}

const NECTAR_PROVIDER = "Nectar"
const AWS_PROVIDER = "AWS"

// ConfigSingleton The config singleton
var ConfigSingleton Config

// ProviderIsNectar
func ProviderIsNectar() bool {
	return ConfigSingleton.Provider == NECTAR_PROVIDER
}

// ProviderIsAws
func ProviderIsAWS() bool {
	return ConfigSingleton.Provider == AWS_PROVIDER
}

// Name get the name of a node
func (s *NodeInst) Name() *string {
	keys := make([]string, 0, len(*s))
	for k := range *s {
		keys = append(keys, k)
	}
	return &keys[0]
}

func (n *NodeInst) Roles() []string {
	return (*n)[*n.Name()].Roles
}

// GetAllRoles Get all the roles which need to be run for
// the nodes
func GetAllRoles() []string {
	var ret []string
	set := make(map[string]bool)
	for _, inst := range ConfigSingleton.Nodes {
		for _, role := range inst.Roles() {
			set[role] = true
		}
	}

	for key := range set {
		ret = append(ret, key)
	}

	return ret
}

// Names Get names in the config
func Names() []string {
	var ret []string
	for _, inst := range ConfigSingleton.Nodes {
		ret = append(ret, *inst.Name())
	}
	return ret
}

// GetNode Get a node by name
func GetNode(name string) (NodeInst, error) {
	for _, inst := range ConfigSingleton.Nodes {
		if *inst.Name() == name {
			return inst, nil
		}
	}
	return NodeInst{}, errors.New("Could not find node")
}

// Config Get the config of a node
func (s *NodeInst) Config() InstanceConfig {
	keys := make([]string, 0, len(*s))
	for k := range *s {
		keys = append(keys, k)
	}
	return (*s)[keys[0]]
}

// SecurityGroupsForAws The SecurityGroupsForAws
func (s *InstanceConfig) SecurityGroupsForAws() []*string {
	var ret []*string
	for _, sg := range s.SecurityGroups {
		ret = append(ret, aws.String(sg))
	}
	return ret
}

// ParseConfig Initially parse the config
// Call this first of all
func ParseConfig(path string) error {
	dat, _ := ioutil.ReadFile(path)
	err := yaml.Unmarshal(dat, &ConfigSingleton)

	if err != nil {
		return err
	}

	if ConfigSingleton.Provider != AWS_PROVIDER && ConfigSingleton.Provider != NECTAR_PROVIDER {
		return fmt.Errorf("Provider must be one of %s or %s", AWS_PROVIDER, NECTAR_PROVIDER)
	}

	return nil
}

// GetConfig Get the config
func GetConfig() Config {
	return ConfigSingleton
}
