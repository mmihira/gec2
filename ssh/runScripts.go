package ssh

import (
	"bufio"
	"fmt"
	"gec2/config"
	"gec2/log"
	"gec2/roles"
	"gec2/nodeContext"
	"gec2/schemaWriter"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

var TMP_SCRIPT_PARENT = "/tmp"
var TMP_SCRIPT_PATH = fmt.Sprintf("%s/%s", TMP_SCRIPT_PARENT, "toRun.sh")

// runCommand Runs a command and also prints the output to a screen prefixed
// by the name of the user
func runCommand(client *ssh.Client, command string, outputPrefix string) error {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Infof("Failed to create session: %s", err)
		return err
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	r, w := io.Pipe()
	er, ew := io.Pipe()
	defer w.Close()
	defer ew.Close()
	session.Stdout = w
	session.Stderr = ew

	// Standard out scanner
	go func() {
		scanner := bufio.NewScanner(r)
		buf := make([]byte, 0, 128*1024)
		scanner.Buffer(buf, 1024*1024)
		for scanner.Scan() {
			log.WithFields(logrus.Fields{
				"node": outputPrefix,
			}).Infof("%s\n", scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Errorf("Scanner to run command got error %s", err)
			session.Stdout = nil
			w.Close()
		}

		defer r.Close()
	}()

	// Error scanner
	go func() {
		scanner := bufio.NewScanner(er)
		buf := make([]byte, 0, 128*1024)
		scanner.Buffer(buf, 1024*1024)
		for scanner.Scan() {
			log.WithFields(logrus.Fields{
				"node": outputPrefix,
			}).Errorf("%s\n", scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Errorf("Scanner to run command got error %s", err)
			session.Stderr = nil
			ew.Close()
		}

		defer er.Close()
	}()

	if err := session.Run(command); err != nil {
		log.Errorf("Failed to run command %s ", err.Error())
		return fmt.Errorf("Failed to run command %s", err)
	}
	return nil
}

func CopyFile(client *ssh.Client, fileContents []byte, location string) error {
	sftp, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftp.Close()

	f, err := sftp.Create(location)
	if err != nil {
		return err
	}

	if _, err := f.Write(fileContents); err != nil {
		return err
	}
	return nil
}

func RunScripts(
	scriptPaths []roles.Script,
	keyFilePath string,
	ctx nodeContext.NodeContext,
	barrier *sync.WaitGroup,
) (bool, error) {
	name := ctx.Name()

	keyFile, err := KeyFile(keyFilePath)
	if err != nil {
		return false, err
	}

	sshConfig := &ssh.ClientConfig{
		User:            "ubuntu",
		Auth:            []ssh.AuthMethod{keyFile},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Until(time.Now().Add(time.Second * 3)),
	}

	for _, scriptPath := range scriptPaths {

		scriptName := fmt.Sprintf("%s/%s", viper.GetString("ROOT_PATH"), scriptPath.FileName())

		fileContents, err := ioutil.ReadFile(scriptName)
		if err != nil {
			log.Fatalf("Read script error: %s", err)
		}

		client, err := ssh.Dial("tcp", fmt.Sprintf(
			"%s:%s",
			ctx.PublicIpAddress(),
			"22",
		), sshConfig)

		if err != nil {
			return false, err
		}

		// Copy the script
		log.Infof("Installing script for %s", name)
		CopyFile(client, fileContents, TMP_SCRIPT_PATH)

		// Copy the schema
		schema, err := schemaWriter.ReadSchemaBytes()
		remoteSchemaPath := fmt.Sprintf(fmt.Sprintf("%s/%s", TMP_SCRIPT_PARENT, schemaWriter.SCHEMA_NAME))
		CopyFile(client, schema, remoteSchemaPath)

		log.Infof("running script %s for %s", scriptPath, name)
		runCommand(client, "chmod u+x /tmp/toRun.sh", name)
		runCommandString := fmt.Sprintf(
			"GECSECRETS='%s' /tmp/toRun.sh %s",
			strings.TrimSpace(config.SecretsMapAsJsonString()),
			scriptPath.Args(),
		)
		runCommand(client, runCommandString, name)

		// Remove script and schema
		log.Infof("Cleanup script for %s", name)
		runCommand(
			client,
			fmt.Sprintf("rm %s; rm %s", remoteSchemaPath, TMP_SCRIPT_PATH),
			name,
		)
	}

	barrier.Done()
	return true, nil
}

func CopyFileRemote(
	fileContents []byte,
	keyFilePath string,
	destination string,
	ctx nodeContext.NodeContext,
	barrier *sync.WaitGroup,
) (bool, error) {
	name := ctx.Name()

	keyFile, err := KeyFile(keyFilePath)
	if err != nil {
		return false, err
	}

	sshConfig := &ssh.ClientConfig{
		User:            "ubuntu",
		Auth:            []ssh.AuthMethod{keyFile},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Until(time.Now().Add(time.Second * 3)),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf(
		"%s:%s",
		ctx.PublicIpAddress(),
		"22",
	), sshConfig)

	if err != nil {
		return false, err
	}

	log.Infof("copying file for %s to %s", name, destination)
	CopyFile(client, fileContents, destination)

	barrier.Done()
	return true, nil
}
