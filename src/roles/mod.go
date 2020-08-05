package roles

import (
	"strings"
)

type Script string

func (s *Script) FileName() string {
	g := strings.Fields(string(*s))
	return g[0]
}

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
	Cmd      string   `json:"cmd"`
}

type Role struct {
	Before []string `json:"before"`
	Steps  []Step   `json:"steps"`
}

type Roles map[string]Role

var RolesSingleton Roles

const (
	ROLE_TYPE_SCRIPT   = "script"
	ROLE_TYPE_COPY     = "copy"
	ROLE_TYPE_TEMPLATE = "template"
	ROLE_TYPE_COMMAND  = "command"
)

// GetConfig Get the config
func GetRoleInst() (Roles, error) {
	return RolesSingleton, nil
}
