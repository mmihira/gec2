package opts

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// Command line options
type AppOpt struct {
	Verbose    bool     `short:"v" long:"verbose" description:"Verbox output" required:"false"`
	Stages     []string `short:"s" long:"stages" descrption:"Stages to run"`
	Roles      []string `short:"r" long:"roles" description:"Roles to run"`
	Nodes      []string `short:"n" long:"nodes" description:"Nodes to run"`
	Args       []string `short:"a" long:"args" description:"Add additional args to scripts that run"`
	ConfigPath string   `long:"configpath" description:"Config file path"`
}

var Opts AppOpt

const (
	STAGE_PROVISION   = "provision"
	STAGE_SSH         = "ssh"
	STAGE_CMD         = "cmd"
	STAGE_LIST_IMAGES = "listImages"
)

func ScriptArgs() string {
	return strings.Join(Opts.Args, " ")
}

func DoStageAll() bool {
	return len(Opts.Stages) == 0
}

func StageListImages() bool {
	for _, s := range Opts.Stages {
		if s == STAGE_LIST_IMAGES {
			return true
		}
	}
	return false
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
		if s == STAGE_SSH {
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
	viper.SetDefault("ROLES_PATH", "/roles")
	viper.SetDefault("LOGS_PATH", "/logs")
	viper.SetDefault("ROOT_PATH", "/")

	viper.BindEnv("SSH_KEY_PATH", "SSH_KEY_PATH")
	viper.BindEnv("EC2_REGION", "EC2_REGION")
	viper.BindEnv("CREDENTIALS_FILE_PATH", "CREDENTIALS_FILE_PATH")
	viper.BindEnv("DEPLOY_CONTEXT_PATH", "DEPLOY_CONTEXT_PATH")
	viper.BindEnv("ROLES_PATH", "ROLES_PATH")
	viper.BindEnv("LOGS_PATH", "LOGS_PATH")
	viper.BindEnv("ROOT_PATH", "ROOT_PATH")

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

// TaggedKeyName This is the tag key used to identify the name of the node
// It must be equal to "Name" because AWS uses that key name to name an instance
func TaggedKeyName() string {
	return "Name"
}

func ParseAppConfigCmdLine() error {
	err := ParseOpts()
	err = SetupViper()
	return err
}
