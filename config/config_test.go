package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
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
		dat := []byte(configFile)
		createConfig(dat)

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
    })
	})
})
