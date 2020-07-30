package aws

import (
	"github.com/spf13/viper"
	"fmt"
	"gec2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// ConnectAWS Connec to a aws ec2 compatible service
// Currently Nectar and AWS are supported
func ConnectAWS() (*ec2.EC2, error) {
	sessionOption := sessionOptions()

	o, err := session.NewSessionWithOptions(sessionOption)
	if err != nil {
		return nil, fmt.Errorf("Error when getting session %s", err)
	}

	ec2svc := ec2.New(o)
	return ec2svc, nil
}

func sessionOptions() session.Options {
	// Get an aws Ec2 session assuming Nectar
	if config.ProviderIsNectar() {
		return session.Options{
			SharedConfigFiles: []string{viper.GetString("CREDENTIALS_FILE_PATH")},
			Config: aws.Config{
				CredentialsChainVerboseErrors: aws.Bool(true),
				Region:                        aws.String(viper.GetString("EC2_REGION")),
				Endpoint:                      aws.String("nova.rc.nectar.org.au:8773/services/Cloud"),
			},
		}
	}

	return session.Options{
		SharedConfigFiles: []string{viper.GetString("CREDENTIALS_FILE_PATH")},
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
			Region:                        aws.String(viper.GetString("EC2_REGION")),
		},
	}
}
