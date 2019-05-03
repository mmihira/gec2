package opts

import (
	flags "github.com/jessevdk/go-flags"
	"fmt"
	"os"
)

// Command line options
type AppOpt struct {
	SshKeyPath    string `long:"sshkey" description:"SSH private key file path" required:"true"`
	Region        string `long:"region" description:"Provider Region" required:"true"`
	Credentials   string `long:"credentials" description:"Credentials path" required:"true"`
	DeployContext string `long:"context" description:"Deploy context path" required:"true"`
	Verbose       bool   `short:"v" long:"verbose" description:"Verbox output" required:"false"`
}

var Opts AppOpt

func ParseOpts() error {
	_, err := flags.ParseArgs(&Opts, os.Args)
	if err != nil {
		return fmt.Errorf("Parsing args got error: %s", err)
	}

	return nil
}
