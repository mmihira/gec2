package config

import (
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

const SECRETS_FILE = "secrets.yaml"
const NECTAR_PROVIDER = "Nectar"
const AWS_PROVIDER = "AWS"

// ConfigSingleton The config singleton
var ConfigSingleton Config

type Config struct {
	Provider string     `json:"provider"`
	Nodes    []NodeInst `json:"nodes"`
	Roles    []string   `json:"roles"`
}

type NodeInst map[string]InstanceConfig

var NodesMap map[string]InstanceConfig
var secretsMap *gabs.Container
var isSecretsValid bool

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

// RolesToRunInOrder The roles to run in order
func RolesToRunInOrder() []string {
	return ConfigSingleton.Roles
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

// ParseConfig Initially parses the config
// This function should be called initially to populate ConfigSingleton
func ParseConfig(pathToConfig string) error {
	dat, _ := ioutil.ReadFile(pathToConfig)
	return createConfig(dat)
}

func ParseFromString(dat string) error {
	return createConfig([]byte(dat))
}

func createConfig(dat []byte) error {
	ConfigSingleton = Config{}
	err := yaml.Unmarshal(dat, &ConfigSingleton)

	if err != nil {
		return err
	}

	if ConfigSingleton.Provider != AWS_PROVIDER &&
		ConfigSingleton.Provider != NECTAR_PROVIDER {
		return fmt.Errorf("Provider must be one of %s or %s", AWS_PROVIDER, NECTAR_PROVIDER)
	}

	NodesMap = map[string]InstanceConfig{}
	for nIdx := range ConfigSingleton.Nodes {
		node := ConfigSingleton.Nodes[nIdx]
		keys := []string{}
		for k := range node {
			keys = append(keys, k)
		}
		if len(keys) != 1 {
			return fmt.Errorf(
				"Incorrect node config format. Only expect one key, the name of the node to define the node configuration",
			)
		}

		NodesMap[keys[0]] = node[keys[0]]
	}

	return nil
}

func ParseSecrets(path string) error {
	secretsDat, readErr := ioutil.ReadFile(path)

	if readErr != nil {
		return readErr
	}

	gabsC, secretsErr := gabs.ParseJSON(secretsDat)
	secretsMap = gabsC

	if secretsErr == nil {
		isSecretsValid = true
	} else {
		return secretsErr
	}
	return nil
}

func SecretsMapAsJsonString() string {
	return secretsMap.String()
}

func IsSecrestValid() bool {
	return isSecretsValid
}

// GetConfig Get the config
func GetConfig() Config {
	return ConfigSingleton
}

// ProviderIsNectar
func ProviderIsNectar() bool {
	return ConfigSingleton.Provider == NECTAR_PROVIDER
}

// ProviderIsAws
func ProviderIsAWS() bool {
	return ConfigSingleton.Provider == AWS_PROVIDER
}
