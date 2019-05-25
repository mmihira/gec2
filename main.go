// VIMTRUN#!
// CREDENTIALS_FILE_PATH="/home/mihira/.ssh/aws-credentials" EC2_REGION="ap-southeast-2" DEPLOY_CONTEXT_PATH="/home/mihira/c/gec2/deploy_context" SSH_KEY_PATH=/home/mihira/.ssh/blocksci/blocksci.pem  "$GOPATH"/bin/gec2 2  -v -r echo -n appUi -s cmd -n minio
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
	"gec2/roles"
	"gec2/schemaWriter"
	gec2ssh "gec2/ssh"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

// The config file should always be names config.yaml
var ConfigFileName = "config.yaml"
var SecretsFileName = "secrets.json"

// The roles file should always be names roles.yaml
var RoleFileName = "roles.yaml"

func main() {
	err := opts.ParseOpts()
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Setup()
	log.Info("Running gec2 v0.1.0")
	// Parse command line options
	if opts.Opts.Verbose {
		log.SetLevel(logrus.DebugLevel)
	}

	// Parse the config
	configPath := fmt.Sprintf("%s/%s", viper.GetString("DEPLOY_CONTEXT_PATH"), ConfigFileName)
	err = config.ParseConfig(configPath)
	if err != nil {
		log.Fatalf("Parsing config got error: %s", err)
	} else {
		log.Info("Config loaded from %s", configPath)
	}

	// Parse the secrets config
	secretsPath := fmt.Sprintf("%s/%s", viper.GetString("DEPLOY_CONTEXT_PATH"), SecretsFileName)
	er := config.ParseSecrets(secretsPath)
	if er != nil {
		log.Infof("Could not parse secrets : %s\n", err)
	} else {
		log.Info("Secrets loaded from %s", secretsPath)
	}

	// Parse the roles
	rolePath := fmt.Sprintf("%s/%s", viper.GetString("DEPLOY_CONTEXT_PATH"), RoleFileName)
	roles.ParseRoles(rolePath)
	if err != nil {
		log.Fatalf("Parsing roles got error: %s", err)
	}
	ec2svc, _ := aws.ConnectAWS()

	if opts.DoStageAll() || opts.StageProvision() {
		// Provision nodes
		provision.EnsureConfigProvisioned(ec2svc)
	}

	// Create node context
	var runningNodes []nodeContext.NodeContext
	for _, name := range config.Names() {
		nodeInst, err := ec2Query.GetInstanceByName(ec2svc, name)
		if err != nil {
			log.Infof("%s: could not be get. Error: %s \n", name, err)
		}

		node, err := config.GetNode(name)
		if err != nil {
			log.Fatalf("Could not fine node %s in config", name)
		}

		runningNodes = append(runningNodes, nodeContext.NodeContext{
			Name:     name,
			Node:     node,
			Instance: nodeInst,
		})
	}

	// Write config information
	err = schemaWriter.WriteSchema(runningNodes)
	if err != nil {
		log.Fatal(err)
	}

	if opts.DoStageAll() || opts.StageSSHCheck() {
		// Wait for SSH access
		allRunning := false
		log.Infof("Waiting for ssh availability...")
		for !allRunning {
			allRunning = true
			resChannel := make(chan gec2ssh.CheckSSHResult)
			for inx, _ := range runningNodes {
				go gec2ssh.CheckSSH(viper.GetString("SSH_KEY_PATH"), &runningNodes[inx], resChannel)
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
			roles.ExecuteRole(runningNodes, roleName)
		}
	}

	log.Infof("Instance fully provisioned!")
}
