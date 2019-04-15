package ssh

import (
	"bufio"
	"fmt"
	"gec2/nodeContext"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"io"
	"sync"
	"time"
)

// createScriptRemote Copy the script to the server
func createScriptRemote(client *ssh.Client, fileContents []byte) (bool, error) {
	sftp, err := sftp.NewClient(client)
	if err != nil {
		return false, err
	}
	defer sftp.Close()

	f, err := sftp.Create("/tmp/toRun.sh")
	if err != nil {
		return false, err
	}

	if _, err := f.Write(fileContents); err != nil {
		return false, err
	}
	return true, nil
}

// runCommand Runs a command and also prints the output to a screen prefixed
// by the name of the user
func runCommand(client *ssh.Client, command string, outputPrefix string) {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Infof("Failed to create session: %s", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	r, w := io.Pipe()
	session.Stdout = w

	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			log.WithFields(log.Fields{
				"node": outputPrefix,
			}).Infof("%s\n", scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Info(err)
		}
	}()

	if err := session.Run(command); err != nil {
		log.Infof("Failed to run: " + err.Error())
	}
	w.Close()
}

func RunScripts(
	scriptPaths []string,
	keyFilePath string,
	ctx *nodeContext.NodeContext,
	barrier *sync.WaitGroup,
) (bool, error) {
	name := ctx.Name

	sshConfig := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			KeyFile(keyFilePath),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Until(time.Now().Add(time.Second * 3)),
	}

	for _, scriptPath := range scriptPaths {

		scriptName := fmt.Sprintf("%s/%s", "/home/mihira/c/gec2/deploy_context", scriptPath)

		fileContents, err := ioutil.ReadFile(scriptName)
		if err != nil { log.Fatalf("Read script error: %s", err) }

		client, err := ssh.Dial("tcp", fmt.Sprintf(
			"%s:%s",
			ctx.PublicIpAddress(),
			"22",
		), sshConfig)

		if err != nil { return false, err }

		log.Infof("sftp install script for %s", name)
		createScriptRemote(client, fileContents)

		log.Infof("running script %s for %s", scriptPath, name)
		runCommand(client, "chmod u+x /tmp/toRun.sh", name)
		runCommand(client, "/tmp/toRun.sh", name)
	}

	barrier.Done()
	return true, nil
}
