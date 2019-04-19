package roles

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
)

// Step Generic step type
type Step struct {
	StepType string   `json:"stepType"`
	Scripts  []string `json:"scripts"`
	Src      string   `json:"src"`
	Dst      string   `json:"dst"`
}

type Role struct {
	Steps []Step `json:"steps"`
}

type Roles map[string]Role

var RolesSingleton Roles

const (
	ROLE_TYPE_SCRIPT   = "script"
	ROLE_TYPE_COPY     = "copy"
	ROLE_TYPE_TEMPLATE = "template"
)

func ParseRoles(path string) error {
	dat, _ := ioutil.ReadFile(path)
	return yaml.Unmarshal(dat, &RolesSingleton)
}

// GetConfig Get the config
func GetRoleInst() (Roles, error) {
	return RolesSingleton, nil
}
