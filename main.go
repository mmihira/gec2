// VIMTRUN#!
// "$GOPATH"/bin/gec2 --credentials="/home/mihira/.ssh/aws-credentials" --region=ap-southeast-2 --sshkey=/home/mihira/.ssh/blocksci/blocksci.pem -v --context=/home/mihira/c/gec2/deploy_context
// VIMTRUN#!

// "$GOPATH"/bin/gec2 --credentials="/home/mihira/.ssh/orca/aws_creds" --region=NeCTAR --sshkey=/home/mihira/.ssh/orca/orca.pem --context=/home/mihira/c/gec2/deploy_context
// "$GOPATH"/bin/gec2 --credentials="/home/mihira/.ssh/aws-credentials" --region=ap-southeast-2 --sshkey=/home/mihira/.ssh/blocksci/blocksci.pem --context=/home/mihira/c/gec2/deploy_context
package main

import (
	"fmt"
	"gec2/aws"
	"gec2/config"
	"gec2/ec2Query"
	"gec2/nodeContext"
	"gec2/opts"
	"gec2/provision"
	"gec2/roles"
	"gec2/schemaWriter"
	gec2ssh "gec2/ssh"
	"github.com/sirupsen/logrus"
	"gec2/log"
	"time"
)

// The config file should always be names config.yaml
var ConfigFileName = "config.yaml"

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
	configPath := fmt.Sprintf("%s/%s", opts.Opts.DeployContext, ConfigFileName)
	err = config.ParseConfig(configPath)
	if err != nil {
		log.Fatalf("Parsing config got error: %s", err)
	}

	// Parse the roles
	rolePath := fmt.Sprintf("%s/%s", opts.Opts.DeployContext, RoleFileName)
	roles.ParseRoles(rolePath)
	if err != nil {
		log.Fatalf("Parsing roles got error: %s", err)
	}
	ec2svc, _ := aws.ConnectAWS()

	// Provision nodes
	provision.EnsureConfigProvisioned(ec2svc)

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

	// Wait for SSH access
	allRunning := false
	log.Infof("Waiting for ssh availability...")
	for !allRunning {
		allRunning = true
		resChannel := make(chan gec2ssh.CheckSSHResult)
		for inx, _ := range runningNodes {
			go gec2ssh.CheckSSH(opts.Opts.SshKeyPath, &runningNodes[inx], resChannel)
		}
		for range runningNodes {
			result := <-resChannel
			log.Infof("ssh status for %s: %v", result.Name, result.DidConnect)
			allRunning = allRunning && result.DidConnect
		}
		time.Sleep(time.Second * 3)
	}

	rolesToRun := config.RolesToRunInOrder()
	for _, roleName := range rolesToRun {
		log.Info("------------------------------------")
		log.Infof("----- Executing role %s: ", roleName)
		log.Info("-----------------------------------")
		roles.ExecuteRole(runningNodes, roleName)
	}

	log.Infof("Instance fully provisioned!")
}
