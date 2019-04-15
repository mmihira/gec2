package provision

import (
	"fmt"
	"gec2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gec2/ec2Query"
	"time"
	log "github.com/sirupsen/logrus"
)

type Gec2ProvisionResult struct {
	Reservation *ec2.Reservation
	Node *config.NodeInst
	InstanceId string
}

func provisionNameForNectar(ec2svc *ec2.EC2, name string) (*ec2.Reservation, error) {
	node, getNodeError := config.GetNode(name)
	if getNodeError != nil {
		return nil, fmt.Errorf("In config get node: %s", getNodeError)
	}
	config := node.Config()

	startinput := &ec2.RunInstancesInput{
		ImageId:        &config.Ami,
		InstanceType:   &config.Type,
		KeyName:        aws.String("orca"),
		MaxCount:       aws.Int64(1),
		MinCount:       aws.Int64(1),
		Placement:      &ec2.Placement{AvailabilityZone: &config.Placement},
		SecurityGroups: config.SecurityGroupsForAws(),
	}

	rresult, err := ec2svc.RunInstances(startinput)
	if err != nil {
		return nil, fmt.Errorf("When requesting run instance %s", err)
	}

	return rresult, nil
}

func provisionNameForAws(ec2svc *ec2.EC2, name string) (*ec2.Reservation, error) {
	node, getNodeError := config.GetNode(name)
	if getNodeError != nil {
		return nil, fmt.Errorf("In config get node: %s", getNodeError)
	}
	config := node.Config()

	startinput := &ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sdh"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(100),
				},
			},
		},
		ImageId:        &config.Ami,
		InstanceType:   &config.Type,
		KeyName:        aws.String("blocksci"),
		MaxCount:       aws.Int64(1),
		MinCount:       aws.Int64(1),
		Placement:      &ec2.Placement{AvailabilityZone: &config.Placement},
		SecurityGroups: config.SecurityGroupsForAws(),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: node.Name(),
					},
				},
			},
		},
	}

	rresult, err := ec2svc.RunInstances(startinput)
	if err != nil {
		return nil, fmt.Errorf("When requesting run instance %s", err)
	}

	return rresult, nil
}

// ProvisionName Provision a node by name
func ProvisionName(ec2svc *ec2.EC2, name string) (*ec2.Reservation, error) {
	if config.ProviderIsNectar() {
		return provisionNameForNectar(ec2svc, name)
	}

	return provisionNameForAws(ec2svc, name)
}

// EnsureConfigProvisioned Ensure the config is provisioned
func EnsureConfigProvisioned(ec2svc *ec2.EC2) error {
	// Get names that should be provisioned
	namesToProvision := map[string]bool{}
	namesAimedRunning := map[string]bool{}
	for _, name := range config.Names() {
		namesToProvision[name] = true
		namesAimedRunning[name] = true
	}

	// Get the currently existing nodes relevant to the config
	provisioned, getInstError := ec2Query.ProvisionedInstances(ec2svc)
	if getInstError != nil {
		log.Fatalf("Getting running instances error: %s", getInstError)
	}

	// Determine what needs to be created
	var provisionedNames []string
	for _, prov := range provisioned {
		provisionedNames = append(provisionedNames, prov.Name)
		log.Infof("%s has allready been created.. skipping\n", prov.Name)
		delete(namesToProvision, prov.Name)
	}

	// Create the missing nodes
	var reservations []Gec2ProvisionResult
	var hasProvisionError = false
	for name, _ := range namesToProvision {
		node, err := config.GetNode(name)
		if err != nil {
			log.Fatalf("Error getting node from config %s", node)
		}

		log.Infof("Creating %s\n", name)
		reservation, provError := ProvisionName(ec2svc, name)
		if provError != nil {
			log.Fatalf("%s Could not be provisioned. Error: %s \n", name, provError)
			hasProvisionError = true
		} else {
			reservations = append(reservations, Gec2ProvisionResult{
				Reservation: reservation,
				Node: &node,
				InstanceId: *reservation.Instances[0].InstanceId,
			})
		}
	}

	if hasProvisionError { return fmt.Errorf("Provision error") }

	// Wait for all the nodes to enter running state
	for hasProvisioned := false; !hasProvisioned; {
		hasProvisioned = true
		for _, res := range reservations {
			g, err := ec2Query.GetInstanceById(ec2svc, res.InstanceId)
			if err != nil {
				log.Fatalf("Error getting instance %s: %s \n", res.InstanceId, err)
			}

			if *g.State.Name == "running" {
				hasProvisioned = hasProvisioned && true
			} else {
				hasProvisioned = hasProvisioned && false
			}
			log.Infof("%s state: %s\n", res.InstanceId, *g.State.Name)
		}
		time.Sleep(time.Second * 3)
	}

	if config.ProviderIsNectar() && len(reservations) > 0 {
		log.Info("Setting tags for nodes...")
		for _, res := range reservations {
			log.Infof("Setting tags for %s", *res.Node.Name())
			input := &ec2.CreateTagsInput{
				Resources: []*string{
					aws.String(res.InstanceId),
				},
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(*res.Node.Name()),
					},
				},
			}
			_, err := ec2svc.CreateTags(input)
			if err != nil { log.Fatal(err.Error()) }
		}
	}

	return nil
}
