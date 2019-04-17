package schemaWriter

import (
	"encoding/json"
	"fmt"
	"gec2/config"
	"gec2/opts"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
)

// The name of the file that is written
var SCHEMA_NAME = "deployed_schema"

// Get the ip address of the instance
func getIpAddress(inst *ec2.Instance) string {
	if config.ProviderIsNectar() {
		return *inst.PrivateIpAddress
	} else {
		return *inst.PublicIpAddress
	}
}

// Build the schema
func buildSchema(instanceMap map[string]*ec2.Instance) (schema *Schema, err error) {
	schema = &Schema{}
	for name, instance := range instanceMap {
		node, err := config.GetNode(name)
		if err != nil {
			return nil, err
		}

		(*schema)[name] = NodeSchema{
			Name:    name,
			KeyName: node[name].KeyName,
			Roles:   node[name].Roles,
			Ip:      getIpAddress(instance),
		}
	}

	return
}

// WriteSchema write the schema to the context dir
func WriteSchema(instanceMap map[string]*ec2.Instance) error {
	f, err := os.Create(fmt.Sprintf("%s/%s.json", opts.Opts.DeployContext, SCHEMA_NAME))
	defer f.Close()

	if err != nil {
		return fmt.Errorf("Could not create output file: %s", err)
	}

	schema, err := buildSchema(instanceMap)
	if err != nil {
		return fmt.Errorf("Could not build schema: %s", err)
	}

	b, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("Could not marshal output file: %s", err)
	}
	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("Could not write output file: %s", err)
	}

	return nil
}
