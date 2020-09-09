// VIMTRUN#!
// ROOT_PATH="/home/mihira/c/gec2/deploy_context" LOGS_PATH="/home/mihira/c/gec2/deploy_context/logs" ROLES_PATH="/home/mihira/c/gec2/deploy_context/roles" CREDENTIALS_FILE_PATH="/home/mihira/.ssh/aws-credentials" EC2_REGION="ap-southeast-2" DEPLOY_CONTEXT_PATH="/home/mihira/c/gec2/deploy_context/context" SSH_KEY_PATH=/home/mihira/.ssh/blocksci/blocksci.pem "$GOPATH"/bin/gec2 -v -s cmd -r test-command -a="date" -n appUi
// VIMTRUN#!
package main

import (
	"fmt"
	"gec2/aws"
	"gec2/config"
	"gec2/ec2Query"
	"gec2/log"
	"gec2/nodeContext"
	"gec2/opts"
	"gec2/provision"
	"gec2/roleexecutor"
	"gec2/schemaWriter"
	gec2ssh "gec2/ssh"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

// The config file should always be named config.yaml
var ConfigFileName = "config.yaml"
var SecretsFileName = "secrets.json"

// The roles file should always be named roles.yaml
var RoleFileName = "roles.yaml"

func ConfigPath() string {
	return fmt.Sprintf("%s/%s", viper.GetString("DEPLOY_CONTEXT_PATH"), ConfigFileName)
}

func SecretsPath() string {
	return fmt.Sprintf("%s/%s", viper.GetString("DEPLOY_CONTEXT_PATH"), SecretsFileName)
}

func setupConfig() {
	log.Infof("Deploy context path : %s", viper.GetString("DEPLOY_CONTEXT_PATH"))
	var err error
	// Parse the config
	err = config.ParseConfig(ConfigPath())
	if err != nil {
		log.Fatalf("Parsing config got error: %s", err)
	} else {
		log.Info("Config loaded from %s", ConfigPath())
	}

	// Parse the secrets
	err = config.ParseSecrets(SecretsPath())
	if err != nil {
		log.Fatalf("Parsing secrets got error: %s", err)
	} else {
		log.Info("Secrets parsed from %s", SecretsPath())
	}
}

func setupRoles() {
	var err error
	// Parse the roles
	rolePath := fmt.Sprintf("%s/%s", viper.GetString("ROLES_PATH"), RoleFileName)
	roleexecutor.ParseRoles(rolePath)
	if err != nil {
		log.Fatalf("Parsing roles got error: %s", err)
	}
}

func main() {
	err := opts.ParseAppConfigCmdLine()
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Setup()
	log.Info("Running gec2 v1.9")
	// Parse command line options
	if opts.Opts.Verbose {
		log.SetLevel(logrus.DebugLevel)
	}

	setupConfig()
	setupRoles()

	var runningNodes []nodeContext.NodeContext
	if opts.StageListImages() {
		ec2svc, _ := aws.ConnectAWS()
		f, err := ec2svc.DescribeImages(nil)
		fmt.Println(f, err)
		return
	}

	if !opts.StageCMD() {
		ec2svc, _ := aws.ConnectAWS()

		if opts.DoStageAll() || opts.StageProvision() {
			// Provision nodes
			provision.EnsureConfigProvisioned(ec2svc)
		}

		// Create node context
		for _, name := range config.Names() {
			nodeInst, err := ec2Query.GetInstanceByName(ec2svc, name)
			if err != nil {
				log.Infof("%s: could not be get. Error: %s \n", name, err)
			}

			node, err := config.GetNode(name)
			if err != nil {
				log.Fatalf("Could not fine node %s in config", name)
			}

			runningNodes = append(runningNodes, &nodeContext.Ec2NodeContext{
				InstName: name,
				Node:     node,
				Instance: nodeInst,
			})
		}
	} else {
		log.Info("Loaded schema to use as node context", ConfigPath())

		// Read the deployed_schema configuration
		schema, err := schemaWriter.ReadSchemaObject()
		if err != nil {
			log.Fatal(err)
		}

		for _, node := range schema.Nodes {
			runningNodes = append(runningNodes, &node)
		}
	}

	if opts.DoStageAll() || opts.StageProvision() {
		// Write config information
		err = schemaWriter.WriteSchema(runningNodes)
		if err != nil {
			log.Fatal(err)
		}
	}

	if opts.DoStageAll() || opts.StageSSHCheck() {
		// Wait for SSH access
		allRunning := false
		log.Infof("Waiting for ssh availability...")
		for !allRunning {
			allRunning = true
			resChannel := make(chan gec2ssh.CheckSSHResult)
			for inx, _ := range runningNodes {
				go gec2ssh.CheckSSH(
					viper.GetString("SSH_KEY_PATH"),
					runningNodes[inx],
					resChannel,
				)
			}
			for range runningNodes {
				result := <-resChannel
				log.Infof("ssh status for %s: %v", result.Name, result.DidConnect)
				allRunning = allRunning && result.DidConnect
			}
			time.Sleep(time.Second * 3)
		}
	}

	if opts.DoStageAll() || opts.StageCMD() {
		rolesToRun := []string{}
		if len(opts.RolesToRun()) > 0 {
			rolesToRun = opts.RolesToRun()
		} else {
			rolesToRun = config.RolesToRunInOrder()
		}

		for _, roleName := range rolesToRun {
			log.Info("------------------------------------")
			log.Infof("----- Executing role %s: ", roleName)
			log.Info("-----------------------------------")
			roleexecutor.ExecuteRole(runningNodes, roleName)
		}
	}

	log.Infof("Instance fully provisioned!")
}
