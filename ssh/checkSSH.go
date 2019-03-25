package ssh

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"golang.org/x/crypto/ssh"
	"io"
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

func some(client *ssh.Client) {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	r, w := io.Pipe()
	session.Stdout = w

	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			fmt.Println(scanner.Text()) // Println will add back the final '\n'
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("end")
	}()

	fmt.Println("here")
	if err := session.Run("which htop"); err != nil {
		fmt.Println("Failed to run: " + err.Error())
	}
	fmt.Println("heres")
	w.Close()

	fmt.Println("now here")
}

func CheckSSH(
	keyFilePath string,
	name string,
	instance *ec2.Instance,
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
		*instance.NetworkInterfaces[0].Association.PublicIp,
		"22",
	), sshConfig)

	if errs != nil {
		resChan <- CheckSSHResult{name, false, true, fmt.Sprintf("Faled to Dial %s", errs)}
	} else {
		resChan <- CheckSSHResult{name, true, false, ""}
	}
}
