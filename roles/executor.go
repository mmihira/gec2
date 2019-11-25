package roles

import (
	"bytes"
	"fmt"
	"gec2/log"
	"gec2/nodeContext"
	"gec2/opts"
	"gec2/schemaWriter"
	gec2ssh "gec2/ssh"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"sync"
	"text/template"
)

// Should really handle coroutine errors properly

// Executes a copy step
func executeStepCopy(
	keyFilePath string,
	ctx nodeContext.NodeContext,
	barrier *sync.WaitGroup,
	step *Step,
) error {
	barrier.Add(1)
	defer barrier.Done()

	srcPath := fmt.Sprintf("%s/%s", viper.GetString("ROOT_PATH"), step.Src)

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		log.Errorf("%s", err)
		return err
	}

	dat, err := ioutil.ReadFile(srcPath)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("Error reading file to template %s", err)
	}

	go gec2ssh.CopyFileRemote(
		dat,
		keyFilePath,
		step.Dst,
		ctx,
		barrier,
	)

	return nil
}

// Executes a template step
func executeStepTemplate(
	keyFilePath string,
	ctx nodeContext.NodeContext,
	barrier *sync.WaitGroup,
	step *Step,
) error {
	tplPath := fmt.Sprintf("%s/%s", viper.GetString("ROOT_PATH"), step.Src)

	dat, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return fmt.Errorf("Error reading file to template %s", err)
	}

	schema, err := schemaWriter.ReadSchema()
	if err != nil {
		return err
	}

	tpl, err := template.New("tpl").Parse(string(dat))
	if err != nil {
		return fmt.Errorf("Error initialising template %s", err)
	}

	log.Debugf("Schema was :: %#v", schema)
	var buff bytes.Buffer
	err = tpl.Execute(&buff, schema)
	if err != nil {
		return fmt.Errorf("Error executing template %s", err)
	}

	log.Debugf("Template was :: %s", buff.String())

	go gec2ssh.CopyFileRemote(
		buff.Bytes(),
		keyFilePath,
		step.Dst,
		ctx,
		barrier,
	)

	return nil
}

// ExecuteRole Execute the role on the node
func ExecuteRole(nodes []nodeContext.NodeContext, roleName string) {
	role, roleFound := RolesSingleton[roleName]
	if !roleFound {
		log.Fatalf("Role %s doesn't exist in roles", roleName)
	}

	if len(role.Steps) == 0 {
		log.Fatalf("Role %s had no steps - aborting", roleName)
	}

	for _, step := range role.Steps {
		switch step.StepType {
		case ROLE_TYPE_COPY:
			log.Infof("Executing step %s: ", step.StepType)

			var wg sync.WaitGroup

			for nodeInx, node := range nodes {
				if node.HasRole(roleName) || opts.HasSpecifiedNode(node.Name()) {
					wg.Add(1)
					go executeStepCopy(viper.GetString("SSH_KEY_PATH"), nodes[nodeInx], &wg, &step)
				}
			}

			wg.Wait()

		case ROLE_TYPE_SCRIPT:
			log.Infof("Executing step %s: ", step.StepType)

			var wg sync.WaitGroup

			for nodeInx, node := range nodes {
				if node.HasRole(roleName) || opts.HasSpecifiedNode(node.Name()) {
					wg.Add(1)
					go gec2ssh.RunScripts(step.Scripts, viper.GetString("SSH_KEY_PATH"), nodes[nodeInx], &wg)
				}
			}

			wg.Wait()
		case ROLE_TYPE_TEMPLATE:
			log.Infof("Executing step %s: ", step.StepType)

			var wg sync.WaitGroup

			for nodeInx, node := range nodes {
				if node.HasRole(roleName) || opts.HasSpecifiedNode(node.Name()) {
					wg.Add(1)
					go executeStepTemplate(viper.GetString("SSH_KEY_PATH"), nodes[nodeInx], &wg, &step)
				}
			}

			wg.Wait()
		default:
			log.Fatalf("Unexpected step %s", step.StepType)
		}
	}
}
