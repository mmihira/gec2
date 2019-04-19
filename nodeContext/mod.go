package nodeContext

import (
	"gec2/config"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Get latest node information
type NodeContext struct {
	Name     string
	Node     config.NodeInst
	Instance *ec2.Instance
}

// PublicIpAddress The PublicIpAddress
func (n *NodeContext) PublicIpAddress() string {
	// For some reason the PrivateIpAddress is the public ip address in Nectar
	if config.ProviderIsNectar() {
		return *n.Instance.PrivateIpAddress
	}

	// Otherwise AWS
	return *n.Instance.PublicIpAddress
}

// HasRole  Check if this node has this role
func (n *NodeContext) HasRole(role string) bool {
	roles := n.Node.Roles()
	for _, nodeRole := range  roles {
		if nodeRole == role {
				return true
		}
  }
  return false
}
