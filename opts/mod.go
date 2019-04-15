package opts

import (
	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"os"
)

// Command line options
type AppOpt struct {
	SshKeyPath    string `long:"sshkey" description:"SSH private key file path" required:"true"`
	Region        string `long:"region" description:"Provider Region" required:"true"`
	Credentials   string `long:"credentials" description:"Credentials path" required:"true"`
	DeployContext string `long:"context" description:"Deploy context path" required:"true"`
}

var Opts AppOpt

func ParseOpts() {
	_, err := flags.ParseArgs(&Opts, os.Args)
	if err != nil {
		log.Fatalf("Parsing args got error: %s", err)
	}
}
