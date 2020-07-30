package nodeContext

import (
	"gec2/config"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type NodeContext interface {
	PublicIpAddress() string
	PrivateIpAddress() string
	HasRole(string) bool
	Name() string
	Roles() []string
	KeyName() string
}

type Ec2NodeContext struct {
	InstName string
	Node     config.NodeInst
	Instance *ec2.Instance
}

// PublicIpAddress The PublicIpAddress
func (n *Ec2NodeContext) PublicIpAddress() string {
	// For some reason the PrivateIpAddress is the public ip address in Nectar
	if config.ProviderIsNectar() {
		return *n.Instance.PrivateIpAddress
	}

	// Otherwise AWS
	return *n.Instance.PublicIpAddress
}

func (n *Ec2NodeContext) Roles() []string {
	return n.Node[n.InstName].Roles
}

func (n *Ec2NodeContext) KeyName() string {
	return n.Node[n.InstName].KeyName
}

// PrivateIpAddress The PrivateIpAddress
func (n *Ec2NodeContext) PrivateIpAddress() string {
	return *n.Instance.PrivateIpAddress
}

// PrivateIpAddress The PrivateIpAddress
func (n *Ec2NodeContext) Name() string {
	return n.InstName
}

// HasRole  Check if this node has this role
func (n *Ec2NodeContext) HasRole(role string) bool {
	roles := n.Node.Roles()
	for _, nodeRole := range roles {
		if nodeRole == role {
			return true
		}
	}
	return false
}

