package roles

import (
	"strings"
)

type Script string

// Name get the name of a node
func (s *Script) FileName() string {
	g := strings.Fields(string(*s))
	return g[0]
}

// Name get the name of a node
func (s *Script) Args() string {
	g := strings.Fields(string(*s))
	if len(g) > 1 {
		return g[1]
	}
	return ""
}

// Step Generic step type
type Step struct {
	StepType string   `json:"stepType"`
	Scripts  []Script `json:"scripts"`
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

// GetConfig Get the config
func GetRoleInst() (Roles, error) {
	return RolesSingleton, nil
}
