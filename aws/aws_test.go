package aws

import (
	"gec2/config"
	"gec2/opts"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aws", func() {
	Context("Nectar config", func() {
		BeforeEach(func() {
			var configFile = `
provider: "Nectar"
nodes:
roles:
 - "init"
`
			err := config.ParseFromString(configFile)
			Expect(err).To(BeNil())

			err = opts.SetupViper()
			Expect(err).To(BeNil())
		})

		It("Uses the correct session options", func() {
			options := sessionOptions()
			Expect(options).To(BeEquivalentTo(session.Options{
				SharedConfigFiles: []string{"/credentials"},
				Config: aws.Config{
					CredentialsChainVerboseErrors: aws.Bool(true),
					Region:                        aws.String(""),
					Endpoint:                      aws.String("nova.rc.nectar.org.au:8773/services/Cloud"),
				},
			}))
		})
	})

	Context("Aws config", func() {
		BeforeEach(func() {
			var configFile = `
provider: "AWS"
nodes:
roles:
 - "init"
`
			err := config.ParseFromString(configFile)
			Expect(err).To(BeNil())

			err = opts.SetupViper()
			Expect(err).To(BeNil())
		})

		It("Uses the correct session options", func() {
			options := sessionOptions()
			Expect(options).To(BeEquivalentTo(session.Options{
				SharedConfigFiles: []string{"/credentials"},
				Config: aws.Config{
					CredentialsChainVerboseErrors: aws.Bool(true),
					Region:                        aws.String(""),
				},
			}))
		})
	})
})
