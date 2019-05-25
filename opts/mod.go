package opts

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
	"fmt"
	"os"
)

// Command line options
type AppOpt struct {
	Verbose       bool   `short:"v" long:"verbose" description:"Verbox output" required:"false"`
	ConfigPath 		string `long:"configpath" description:"Config file path"`
}

var Opts AppOpt

func ParseOpts() error {
	_, err := flags.ParseArgs(&Opts, os.Args)
	if err != nil {
		return fmt.Errorf("Parsing args got error: %s", err)
	}

	viper.SetDefault("EC2_REGION", "")
	viper.SetDefault("SSH_KEY_PATH", "/sshKey")
	viper.SetDefault("CREDENTIALS_FILE_PATH", "/credentials")
	viper.SetDefault("DEPLOY_CONTEXT_PATH", "/context")

	viper.BindEnv("SSH_KEY_PATH", "SSH_KEY_PATH")
	viper.BindEnv("EC2_REGION", "EC2_REGION")
	viper.BindEnv("CREDENTIALS_FILE_PATH", "CREDENTIALS_FILE_PATH")
	viper.BindEnv("DEPLOY_CONTEXT_PATH", "DEPLOY_CONTEXT_PATH")

	viper.SetConfigType("yaml")

	if len(Opts.ConfigPath) != 0 {
		fmt.Printf("Loading config from %s", Opts.ConfigPath)
		viper.SetConfigFile(Opts.ConfigPath)
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Printf("Fatal error reading config file")
		}
	}

	return nil
}
