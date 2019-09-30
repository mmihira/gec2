package opts

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
	"os"
)

// Command line options
type AppOpt struct {
	Verbose    bool     `short:"v" long:"verbose" description:"Verbox output" required:"false"`
	Stages     []string `short:"s" long:"stages" descrption:"Stages to run"`
	Roles      []string `short:"r" long:"roles" description:"Roles to run"`
	Nodes      []string `short:"n" long:"nodes" description:"Nodes to run"`
	ConfigPath string   `long:"configpath" description:"Config file path"`
}

var Opts AppOpt

const (
	STAGE_PROVISION = "provision"
	STAGE_SSH       = "ssh"
	STAGE_CMD       = "cmd"
)

func DoStageAll () bool {
	return len(Opts.Stages) == 0
}

func StageProvision() bool {
	for _, s := range Opts.Stages {
		if s == STAGE_PROVISION {
			return true
		}
	}
	return false
}

func StageSSHCheck() bool {
	for _, s := range Opts.Stages {
		if s ==  STAGE_SSH {
			return true
		}
	}
	return false
}

func StageCMD() bool {
	for _, s := range Opts.Stages {
		if s == STAGE_CMD {
			return true
		}
	}
	return false
}

func RolesToRun() []string {
	return Opts.Roles
}

func HasSpecifiedNode(node string) bool {
	for _, n := range Opts.Nodes {
		if n == node {
			return true
		}
	}
	return false
}

func ParseOpts() error {
	_, err := flags.ParseArgs(&Opts, os.Args)
	if err != nil {
		return fmt.Errorf("Parsing args got error: %s", err)
	}

	return nil
}

func SetupViper() error {
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

func ParseAppConfigCmdLine() error {
	err := ParseOpts()
	err = SetupViper()
	return err
}
