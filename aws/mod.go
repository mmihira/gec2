package aws

import (
	"fmt"
	"gec2/config"
	"gec2/opts"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// ConnectAWS Connec to a aws ec2 compatible service
// Currently Nectar and AWS are supported
func ConnectAWS() (*ec2.EC2, error) {
	// Get an aws Ec2 session assuming Nectar
	if config.ProviderIsNectar() {
		s := session.Options{
			SharedConfigFiles: []string{opts.Opts.Credentials},
			Config: aws.Config{
				CredentialsChainVerboseErrors: aws.Bool(true),
				Region:                        aws.String(opts.Opts.Region),
				Endpoint:                      aws.String("nova.rc.nectar.org.au:8773/services/Cloud"),
			},
		}
		o, err := session.NewSessionWithOptions(s)
		if err != nil {
			return nil, fmt.Errorf("Error when getting session %s", err)
		}
		ec2svc := ec2.New(o)

		return ec2svc, nil
	}

	// Get an aws Ec2 session assuming AWS
	s := session.Options{
		SharedConfigFiles: []string{opts.Opts.Credentials},
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
			Region:                        aws.String(opts.Opts.Region),
		},
	}

	o, err := session.NewSessionWithOptions(s)
	if err != nil {
		return nil, fmt.Errorf("Error when getting session %s", err)
	}
	ec2svc := ec2.New(o)

	return ec2svc, nil
}
