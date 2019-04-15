package nodeContext
import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"gec2/config"
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
	return *n.Instance.NetworkInterfaces[0].Association.PublicIp
}
