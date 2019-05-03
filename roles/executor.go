package roles

import (
	"bytes"
	"fmt"
	"gec2/nodeContext"
	"gec2/opts"
	"gec2/schemaWriter"
	gec2ssh "gec2/ssh"
	"gec2/log"
	"io/ioutil"
	"os"
	"sync"
	"text/template"
)

// Should really handle coroutine errors properly

// Executes a copy step
func executeStepCopy(
	keyFilePath string,
	ctx *nodeContext.NodeContext,
	barrier *sync.WaitGroup,
	step *Step,
) error {
	barrier.Add(1)
	defer barrier.Done()

	srcPath := fmt.Sprintf("%s/%s", opts.Opts.DeployContext, step.Src)

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
	ctx *nodeContext.NodeContext,
	barrier *sync.WaitGroup,
	step *Step,
) error {
	tplPath := fmt.Sprintf("%s/%s", opts.Opts.DeployContext, step.Src)

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

// ExecuteRole
func ExecuteRole(nodes []nodeContext.NodeContext, roleName string) {
	role, roleFound := RolesSingleton[roleName]
	if !roleFound {
		log.Fatalf("Role %s doesn't exist in roles", roleName)
	}

	for _, step := range role.Steps {
		switch step.StepType {
		case ROLE_TYPE_COPY:
			log.Infof("Executing step %s: ", step.StepType)

			var wg sync.WaitGroup

			for nodeInx, node := range nodes {
				if node.HasRole(roleName) {
					wg.Add(1)
					go executeStepCopy(opts.Opts.SshKeyPath, &nodes[nodeInx], &wg, &step)
				}
			}

			wg.Wait()

		case ROLE_TYPE_SCRIPT:
			log.Infof("Executing step %s: ", step.StepType)

			var wg sync.WaitGroup

			for nodeInx, node := range nodes {
				if node.HasRole(roleName) {
					wg.Add(1)
					go gec2ssh.RunScripts(step.Scripts, opts.Opts.SshKeyPath, &nodes[nodeInx], &wg)
				}
			}

			wg.Wait()
		case ROLE_TYPE_TEMPLATE:
			log.Infof("Executing step %s: ", step.StepType)

			var wg sync.WaitGroup

			for nodeInx, node := range nodes {
				if node.HasRole(roleName) {
					wg.Add(1)
					go executeStepTemplate(opts.Opts.SshKeyPath, &nodes[nodeInx], &wg, &step)
				}
			}

			wg.Wait()
		default:
			log.Fatalf("Unexpected step %s", step.StepType)
		}
	}
}
