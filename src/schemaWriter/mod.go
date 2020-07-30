package schemaWriter

import (
	"encoding/json"
	"fmt"
	"gec2/config"
	"gec2/nodeContext"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

// The name of the file that is written
var SCHEMA_NAME = "deployed_schema.json"

// Build the schema
func buildSchema(instances []nodeContext.NodeContext) (schema *Schema, err error) {
	schema = &Schema{
		Nodes:     map[string]NodeSchema{},
		WithRoles: map[string][]NodeSchema{},
	}

	for _, ctxt := range instances {
		(*schema).Nodes[ctxt.Name()] = NodeSchema{
			InstName:    ctxt.Name(),
			InstKeyName: ctxt.KeyName(),
			InstRoles:   ctxt.Roles(),
			Ip:          ctxt.PublicIpAddress(),
			PrivateIp:   ctxt.PrivateIpAddress(),
		}
	}

	for _, role := range config.GetAllRoles() {
		nodesInRole := []NodeSchema{}
		for _, ctxt := range instances {
			if ctxt.HasRole(role) {
				nodesInRole = append(nodesInRole, (*schema).Nodes[ctxt.Name()])
			}
		}

		(*schema).WithRoles[role] = nodesInRole
	}

	return
}

// WriteSchema write the schema to the context dir
func WriteSchema(instances []nodeContext.NodeContext) error {
	schemaPath := fmt.Sprintf("%s/%s", viper.GetString("DEPLOY_CONTEXT_PATH"), SCHEMA_NAME)

	err := os.Remove(schemaPath)
	f, err := os.Create(schemaPath)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("Could not create output file: %s", err)
	}

	schema, err := buildSchema(instances)
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

func ReadSchemaBytes() ([]byte, error) {
	path := fmt.Sprintf("%s/%s", viper.GetString("DEPLOY_CONTEXT_PATH"), SCHEMA_NAME)
	dat, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("Error reading schema %s", err)
	}

	return dat, nil
}

func ReadSchemaObject() (*Schema, error) {
	data, err := ReadSchemaBytes()
	if err != nil {
		return nil, err
	}

	var schema Schema
	if err = json.Unmarshal(data, &schema); err != nil {
		return nil, err
	} else {
		return &schema, nil
	}
}

// ReadSchema Read the schema
func ReadSchema() (map[string]interface{}, error) {
	path := fmt.Sprintf("%s/%s", viper.GetString("DEPLOY_CONTEXT_PATH"), SCHEMA_NAME)
	dat, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("Error reading schema %s", err)
	}

	var d map[string]interface{}
	err = json.Unmarshal(dat, &d)
	if err != nil {
		return nil, fmt.Errorf("Error decoding schema %s", err)
	}

	return d, nil
}
