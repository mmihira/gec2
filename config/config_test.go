package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var _ = Describe("ConfigSingleton", func() {
	var configFile = `
provider: "Nectar"
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
roles:
  - "init"
  `

  firstNode :=  func() *InstanceConfig {
    fn := ConfigSingleton.Nodes[0]["node1"]
    return &fn
  }

	Context("Initially parsing the config", func() {
		It("Should create the config", func() {
			dat := []byte(configFile)
			err := createConfig(dat)
			Expect(err).To(BeNil())
		})
	})

	Context("With config parsed", func() {
		BeforeEach(func() {
			dat := []byte(configFile)
			err := createConfig(dat)
			Expect(err).To(BeNil())
		})

		It("Provider parsed", func() {
    	Expect(ConfigSingleton.Provider).To(Equal("Nectar"))
    })

		It("Parsed 1 node", func() {
    	Expect(ConfigSingleton.Nodes).Should(HaveLen(1))
    })

    Context("With node", func() {
      It("Has correct key", func() {
        Expect(ConfigSingleton.Nodes[0]).Should(HaveKey("node1"))
      })

      It("Roles are correct", func() {
        node := firstNode()
        Expect(node.Roles).Should(ConsistOf(
          "init",
        ))
      })

      It("SecurityGroups are correct", func() {
        node := firstNode()
        Expect(node.SecurityGroups).Should(ConsistOf(
          "ssh",
          "www",
        ))
      })

      It("SshParams.HostName is correct", func() {
        node := firstNode()
        Expect(node.SshParam.HostName).Should(Equal("ubuntu"))
      })

      It("SshParams.Port is correct", func() {
        node := firstNode()
        Expect(node.SshParam.Port).Should(Equal(int32(80)))
      })

			It("DeviceMappingsForAws is default", func() {
        node := firstNode()
        Expect(node.DeviceMappingsForAws()).Should(BeEquivalentTo(
					[]*ec2.BlockDeviceMapping{
						&ec2.BlockDeviceMapping {
							DeviceName: aws.String("/dev/sda1"),
							Ebs: &ec2.EbsBlockDevice{
								VolumeSize: aws.Int64(100),
							},
						},
					},
				))
			})
    })
	})

	Context("With config with device mappings", func() {
	var withDeviceConfig = `
provider: "Nectar"
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
      ebsMappings:
        - deviceName: "/dev/gda1"
          VolumeSize: 200
        - deviceName: "/dev/nda1"
          VolumeSize: 300
roles:
  - "init"
  `

		BeforeEach(func() {
			dat := []byte(withDeviceConfig)
			err := createConfig(dat)
			Expect(err).To(BeNil())
		})

		It("Parsed 1 node", func() {
    	Expect(ConfigSingleton.Nodes).Should(HaveLen(1))
    })

		It("DeviceMappingsForAws is as configured", func() {
			node := firstNode()
			Expect(node.DeviceMappingsForAws()).Should(BeEquivalentTo(
				[]*ec2.BlockDeviceMapping{
					&ec2.BlockDeviceMapping {
						DeviceName: aws.String("/dev/gda1"),
						Ebs: &ec2.EbsBlockDevice{
							VolumeSize: aws.Int64(200),
						},
					},
					&ec2.BlockDeviceMapping {
						DeviceName: aws.String("/dev/nda1"),
						Ebs: &ec2.EbsBlockDevice{
							VolumeSize: aws.Int64(300),
						},
					},
				},
			))
		})
	})
})

var _ = Describe("NodesMap", func() {
	var configFile = `
provider: "Nectar"
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
roles:
  - "init"
  `

  firstNode :=  func() *InstanceConfig {
    fn := NodesMap["node1"]
    return &fn
  }

	Context("With config parsed", func() {
		BeforeEach(func() {
			dat := []byte(configFile)
			err := createConfig(dat)
			Expect(err).To(BeNil())
		})

    Context("With node", func() {
      It("Roles are correct", func() {
        node := firstNode()
        Expect(node.Roles).Should(ConsistOf(
          "init",
        ))
      })

      It("SecurityGroups are correct", func() {
        node := firstNode()
        Expect(node.SecurityGroups).Should(ConsistOf(
          "ssh",
          "www",
        ))
      })

      It("SshParams.HostName is correct", func() {
        node := firstNode()
        Expect(node.SshParam.HostName).Should(Equal("ubuntu"))
      })

      It("SshParams.Port is correct", func() {
        node := firstNode()
        Expect(node.SshParam.Port).Should(Equal(int32(80)))
      })
    })
	})
})
