package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"

	"github.com/go-logr/logr"
	util "github.com/kubbee/ccloud-cluster-operator/internal"
)

func newSR(cloud string, geo string, environmentId string, logger *logr.Logger) (*util.SchemaRegistry, error) {
	//confluent sr cluster enable --cloud aws --geo us --environment env-9kkdnv --output json
	logger.Info("Start::newSchemaRegistry")

	cmd := exec.Command("/bin/confluent", "sr", "cluster", "enable", "--cloud", cloud, "--geo", geo, "--environment", environmentId, "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		logger.Error(err, "error to run the command to create SchemaRegistry")
		return &util.SchemaRegistry{}, err
	}

	output := cmdOutput.Bytes()
	message, _ := getOutput(output)

	sr := &util.SchemaRegistry{}

	if err := json.Unmarshal([]byte(message), sr); err != nil {
		logger.Error(err, "error to parse SchemaRegistry Json to Struct")
		return &util.SchemaRegistry{}, err
	}

	return sr, nil

}

/*
 *
 */
func getSR(environmentId string, logger *logr.Logger) (*util.SchemaRegistry, error) {
	logger.Info("func GetSechemaRegistry started")

	cmd := exec.Command("/bin/confluent", "schema-registry", "cluster", "describe", "--environment", environmentId, "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err_ := cmd.Run(); err_ != nil {
		return &util.SchemaRegistry{}, nil
	} else {
		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		sr := &util.FullSchemaRegistry{}

		if jerr := json.Unmarshal([]byte(message), sr); jerr != nil {
			return &util.SchemaRegistry{}, jerr
		}

		smallSR := &util.SchemaRegistry{
			Id:         sr.ClusterId,
			EnpointURL: sr.EndpointURL,
		}

		return smallSR, nil
	}
}

// this function creates a new Schema Registry Api Key
func newSRApiKey(schemaRegistryId string, apiKeyName string, logger *logr.Logger) (*util.ApiKey, error) {
	logger.Info("Start::newSRApiKey")

	cmd := exec.Command("/bin/confluent", "api-key", "create", "--resource", schemaRegistryId, "--description", apiKeyName, "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		logger.Error(err, "error to create the SR api-key for the application")
		return &util.ApiKey{}, err
	} else {

		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		apiKey := &util.ApiKey{}
		json.Unmarshal([]byte(message), apiKey)

		if apiKey.Api != "" && apiKey.Secret != "" {
			return apiKey, nil
		} else {
			return &util.ApiKey{}, errors.New("error to parse the api-key")
		}
	}
}
