package ec2Query

import (
	"gec2/config"
	"gec2/opts"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("inputForGetTaggedRunningInstances", func() {
	Context("With generic config", func() {
		BeforeEach(func() {
			var configFile = `
provider: "AWS"
nodes:
  - node1:
      ami: "ami-0b76c3b150c6b1423"
      type: "t2.small"
      placement: "ap-southeast-2a"
      keyname: "schedulermq"
      attach_volume: false
      security_groups:
        - "ssh"
        - "www"
      roles:
        - "init"
      sshParams:
        hostName: "ubuntu"
        port: 80
  - node2:
      ami: "ami-0b76c3b150c6b1423"
      type: "t2.small"
      placement: "ap-southeast-2a"
      keyname: "schedulermq"
      attach_volume: false
      security_groups:
        - "ssh"
        - "www"
      roles:
        - "init"
      sshParams:
        hostName: "ubuntu"
        port: 80
roles:
 - "init"
`
			err := config.ParseFromString(configFile)
			Expect(err).To(BeNil())

			err = opts.SetupViper()
			Expect(err).To(BeNil())
		})

		It("Creates filter params correctly from config", func() {
			e := inputForGetTaggedRunningInstances()
			Expect(e).To(BeEquivalentTo(
				ec2.DescribeInstancesInput{
					Filters: []*ec2.Filter{
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
								aws.String("node1"),
								aws.String("node2"),
							},
						},
					},
				},
			))
		})
	})
})

var _ = Describe("inputFiltersForGetInstanceByName", func() {
	Context("With generic config", func() {
		BeforeEach(func() {
			err := opts.SetupViper()
			Expect(err).To(BeNil())
		})

		It("Creates filter params correctly from config", func() {
			name := "SomeNode"
			e := inputFiltersForGetInstanceByName(name)
			Expect(e).To(BeEquivalentTo(
				ec2.DescribeInstancesInput{
					Filters: []*ec2.Filter{
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
					},
				},
			))
		})
	})
})
