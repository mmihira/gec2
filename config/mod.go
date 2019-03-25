package config

import (
	"errors"
	"github.com/ghodss/yaml"
	"github.com/aws/aws-sdk-go/aws"
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
	AnsibleHostGroups []string `json:"ansible_host_groups"`
	EnvInjection      []string `json:"env_injection"`
	SecurityGroups    []string `json:"security_groups"`
}

type NodeInst map[string]InstanceConfig
type Config struct {
	Nodes []NodeInst `json:"nodes"`
}

// Name get the name of a node
func (s *NodeInst) Name() *string {
	keys := make([]string, 0, len(*s))
	for k := range *s {
		keys = append(keys, k)
	}
	return &keys[0]
}

// Names Get names in the config
func (c *Config) Names() []string {
	var ret []string
	for _, inst := range c.Nodes {
		ret = append(ret, *inst.Name())
	}
	return ret
}

// GetNode Get a node by name
func (c *Config) GetNode(name string) (NodeInst, error) {
	for _, inst := range c.Nodes {
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

func (s *InstanceConfig) SecurityGroupsForAws() []*string {
	var ret []*string
	for _, sg := range s.SecurityGroups {
		ret = append(ret, aws.String(sg))
	}
	return ret
}

// GetConfig Get the config
func GetConfig() (Config, error) {
	dat, _ := ioutil.ReadFile("./config.yaml")
	var p2 Config
	marshallError := yaml.Unmarshal(dat, &p2)
	if marshallError != nil {
		return Config{}, marshallError
	}
	return p2, nil
}
