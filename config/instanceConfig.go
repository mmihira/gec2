package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type SshParam struct {
	HostName string `json:"hostName"`
	Port     int32  `json:"port"`
}

type EbsMapping struct {
	DeviceName string `json:"deviceName"`
	VolumeSize int64    `json:"volumeSize"`
}

// InstanceConfig
type InstanceConfig struct {
	Ami              string   `json:"ami"`
	Type             string   `json:"type"`
	Placement        string   `json:"placement"`
	AttachVolume     bool     `json:"attach_volume"`
	Volume           string   `json:"volume"`
	VolumeMountPoint string   `json:"volume_mount_point"`
	VolumeMountDir   string   `json:"volume_mount_dir"`
	KeyName          string   `json:"keyname"`
	EnvInjection     []string `json:"env_injection"`
	SecurityGroups   []string `json:"security_groups"`
	Roles            []string `json:"roles"`
	SshParam         `json:"sshParams"`
	EbsMappings      []EbsMapping `json:"ebsMappings"`
}

// SecurityGroupsForAws The SecurityGroupsForAws
func (s *InstanceConfig) SecurityGroupsForAws() []*string {
	var ret []*string
	for _, sg := range s.SecurityGroups {
		ret = append(ret, aws.String(sg))
	}
	return ret
}

// DeviceMappingsForAws Get Device Mappings for provisoning this instance
func (s *InstanceConfig) DeviceMappingsForAws() []*ec2.BlockDeviceMapping {
	var ret []*ec2.BlockDeviceMapping

	if len(s.EbsMappings) == 0 {
		return []*ec2.BlockDeviceMapping{
			&ec2.BlockDeviceMapping {
				DeviceName: aws.String("/dev/sda1"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(100),
				},
			},
		}
	} else {
		for _, mapping := range s.EbsMappings {
			mapping := ec2.BlockDeviceMapping{
					DeviceName: aws.String(mapping.DeviceName),
					Ebs: &ec2.EbsBlockDevice{
						VolumeSize: aws.Int64(mapping.VolumeSize),
					},
			}

			ret = append(ret, &mapping)
		}
		return ret
	}
}

