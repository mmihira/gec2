package ssh

import (
	"fmt"
	"gec2/nodeContext"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"time"
)

const DEFAULT_SSH_PORT = 22
const DEFAULT_USER = "ubuntu"

type CheckSSHResult struct {
	Name        string
	DidConnect  bool
	DidError    bool
	ErrorString string
}

func KeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err

	}

	return ssh.PublicKeys(key), nil
}

// CheckSSH Check if ssh access is available
func CheckSSH(
	keyFilePath string,
	ctx nodeContext.NodeContext,
	resChan chan CheckSSHResult,
) {

	keyFile, err := KeyFile(keyFilePath)
	if err != nil {
		resChan <- CheckSSHResult{ctx.Name(), false, true, fmt.Sprintf("Faled to Dial %s", err)}
		return
	}

	sshConfig := &ssh.ClientConfig{
		User:            "ubuntu",
		Auth:            []ssh.AuthMethod{keyFile},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Until(time.Now().Add(time.Second * 3)),
	}

	_, errs := ssh.Dial("tcp", fmt.Sprintf(
		"%s:%s",
		ctx.PublicIpAddress(),
		"22",
	), sshConfig)

	if errs != nil {
		resChan <- CheckSSHResult{ctx.Name(), false, true, fmt.Sprintf("Faled to Dial %s", errs)}
	} else {
		resChan <- CheckSSHResult{ctx.Name(), true, false, ""}
	}
}
