package ec2Query

import (
	"fmt"
	"gec2/config"
	"gec2/opts"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// This package holds ec2 query functions.
// Instance state names can be
// The state of the instance (
// pending |
// running |
// shutting-down |
// terminated |
// stopping |
// stopped

type ProvInst struct {
	Name       string
	Properties ec2.Instance
}

func inputForGetTaggedRunningInstances() ec2.DescribeInstancesInput {
	var names []string
	for _, node := range config.GetConfig().Nodes {
		names = append(names, *node.Name())
	}

	var tagFilters []*string
	for _, name := range names {
		tagFilters = append(tagFilters, aws.String(name))
	}

	tagFilterLabel := fmt.Sprintf("tag:%s", opts.TaggedKeyName())

	filters := []*ec2.Filter{
		&ec2.Filter{
			Name: aws.String("instance-state-name"),
			Values: []*string{
				aws.String("running"),
				aws.String("pending"),
			},
		},
		&ec2.Filter{
			Name:   aws.String(tagFilterLabel),
			Values: tagFilters,
		},
	}

	return ec2.DescribeInstancesInput{
		Filters: filters,
	}
}

// GetTaggedRunningInstances Pending and running instaances are returned
// Only instances tagged with their names are returned
func GetTaggedRunningInstances(ec2svc *ec2.EC2) (*ec2.DescribeInstancesOutput, error) {
	input := inputForGetTaggedRunningInstances()
	return ec2svc.DescribeInstances(&input)
}

// GetInstanceById
func GetInstanceById(ec2svc *ec2.EC2, awsInstanceId string) (*ec2.Instance, error) {
	filters := []*ec2.Filter{
		&ec2.Filter{
			Name: aws.String("instance-id"),
			Values: []*string{
				aws.String(awsInstanceId),
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
		return nil, fmt.Errorf("Could not find instance %s", awsInstanceId)
	}

	if len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("Could not find instance %s", awsInstanceId)
	}

	return result.Reservations[0].Instances[0], nil
}

func inputFiltersForGetInstanceByName(name string) (ec2.DescribeInstancesInput) {
	tagFilterLabel := fmt.Sprintf("tag:%s", opts.TaggedKeyName())
	filters := []*ec2.Filter{
		&ec2.Filter{
			Name: aws.String("instance-state-name"),
			Values: []*string{
				aws.String("running"),
				aws.String("pending"),
			},
		},
		&ec2.Filter{
			Name: aws.String(tagFilterLabel),
			Values: []*string{
				aws.String(name),
			},
		},
	}

	return ec2.DescribeInstancesInput{
		Filters: filters,
	}
}

// GetInstanceByName Get instance by name
func GetInstanceByName(ec2svc *ec2.EC2, name string) (*ec2.Instance, error) {
	inputFilters := inputFiltersForGetInstanceByName(name)

	result, err := ec2svc.DescribeInstances(&inputFilters)
	if err != nil {
		return nil, err
	}

	if len(result.Reservations) == 0 {
		return nil, fmt.Errorf(
			"Could not find instance %s. Expected Reservations in %+v to have length greater than 0",
			name,
			result,
		)
	}

	if len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf(
			"Could not find instance %s. Expected Reservations[0].Instances in %+v to have length greater than 0",
			name,
			result,
		)
	}

	return result.Reservations[0].Instances[0], nil
}

// ProvisionedInstances Returned fully provisionedInstances that
// are part of the config
func ProvisionedInstances(ec2svc *ec2.EC2) ([]ProvInst, error) {
	currentRunningTaggedInstances, err := GetTaggedRunningInstances(ec2svc)
	if err != nil {
		return nil, err
	}

	var ret []ProvInst

	for _, reservation := range currentRunningTaggedInstances.Reservations {
		for _, instance := range reservation.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == opts.TaggedKeyName() {
					ret = append(ret, ProvInst{
						*tag.Value,
						*instance,
					})
				}
			}
		}
	}

	return ret, nil
}
