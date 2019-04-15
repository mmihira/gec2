package ssh

import (
	"fmt"
	"gec2/nodeContext"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"time"
)

type CheckSSHResult struct {
	Name        string
	DidConnect  bool
	DidError    bool
	ErrorString string
}

func KeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}


// CheckSSH Check if ssh access is available
func CheckSSH(
	keyFilePath string,
	ctx *nodeContext.NodeContext,
	resChan chan CheckSSHResult,
) {
	sshConfig := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			KeyFile(keyFilePath),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Until(time.Now().Add(time.Second * 3)),
	}

	_, errs := ssh.Dial("tcp", fmt.Sprintf(
		"%s:%s",
		ctx.PublicIpAddress(),
		"22",
	), sshConfig)

	if errs != nil {
		resChan <- CheckSSHResult{ctx.Name, false, true, fmt.Sprintf("Faled to Dial %s", errs)}
	} else {
		resChan <- CheckSSHResult{ctx.Name, true, false, ""}
	}
}
