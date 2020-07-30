package schemaWriter

type Tag struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type NodeSchema struct {
	InstName      string   `json:"name"`
	InstKeyName   string   `json:"keyname"`
	InstRoles     []string `json:"roles"`
	Ip        string   `json:"ip"`
	PrivateIp string   `json:"privateIp"`
}

type Schema struct {
	Nodes     map[string]NodeSchema   `json:"nodes"`
	WithRoles map[string][]NodeSchema `json:"withRoles"`
}



// PublicIpAddress The PublicIpAddress
func (n *NodeSchema) PublicIpAddress() string {
	// Otherwise AWS
	return n.Ip
}

func (n *NodeSchema) Roles() []string {
	return n.InstRoles
}

func (n *NodeSchema) KeyName() string {
	return n.InstKeyName
}

// PrivateIpAddress The PrivateIpAddress
func (n *NodeSchema) PrivateIpAddress() string {
	return n.PrivateIp
}

// PrivateIpAddress The PrivateIpAddress
func (n *NodeSchema) Name() string {
	return n.InstName
}

// HasRole  Check if this node has this role
func (n *NodeSchema) HasRole(role string) bool {
	roles := n.Roles()
	for _, nodeRole := range roles {
		if nodeRole == role {
			return true
		}
	}
	return false
}
