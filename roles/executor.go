package roles

import (
	"gec2/nodeContext"
	"sync"
	"gec2/opts"
	gec2ssh "gec2/ssh"
	log "github.com/sirupsen/logrus"
)

func ExecuteRole(nodes []nodeContext.NodeContext, role *Role) {
	for _, step := range role.Steps {
		switch step.StepType {
		case ROLE_TYPE_SCRIPT:
			log.Infof("Executing step %s: ", step.StepType)

			var wg sync.WaitGroup

			for _, node := range  nodes {
				wg.Add(1)
				go gec2ssh.RunScripts(step.Scripts, opts.Opts.SshKeyPath, &node, &wg)
			}
			wg.Wait()
		default:
			log.Fatalf("Unexpected step %s", step.StepType)
		}
	}
}
