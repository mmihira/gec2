package ec2Query

import (
	"fmt"
	"gec2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// This package holds ec2 query functions.
// GetTaggedRunningInstances: getting the instance which are currently running.
// GetInstanceByName: get an instance by name

// Instance state names can be
// The state of the instance (
// pending |
// running |
// shutting-down |
// terminated |
// stopping |
// stopped

// GetTaggedRunningInstances
func GetTaggedRunningInstances(ec2svc *ec2.EC2) (*ec2.DescribeInstancesOutput, error) {
	p2, _ := config.GetConfig()

	var names []string
	for _, node := range p2.Nodes {
		names = append(names, *node.Name())
	}

	var tagFilters []*string
	for _, name := range names {
		tagFilters = append(tagFilters, aws.String(name))
	}

	filters := []*ec2.Filter{
		&ec2.Filter{
			Name: aws.String("instance-state-name"),
			Values: []*string{
				aws.String("running"),
				aws.String("pending"),
			},
		},
		&ec2.Filter{
			Name:   aws.String("tag:Name"),
			Values: tagFilters,
		},
	}

	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	return ec2svc.DescribeInstances(input)
}

// GetInstanceByName Get instance by name
func GetInstanceByName(ec2svc *ec2.EC2, name string) (*ec2.Instance, error) {
	filters := []*ec2.Filter{
		&ec2.Filter{
			Name: aws.String("instance-state-name"),
			Values: []*string{
				aws.String("running"),
				aws.String("pending"),
			},
		},
		&ec2.Filter{
			Name: aws.String("tag:Name"),
			Values: []*string{
				aws.String(name),
			},
		},
	}

	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	result, err := ec2svc.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	if len(result.Reservations) == 0 {
		return nil, fmt.Errorf("Could not find instance %s", name)
	}

	if len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("Could not find instance %s", name)
	}

	return result.Reservations[0].Instances[0], nil
}
