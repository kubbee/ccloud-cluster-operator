package services

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"github.com/go-logr/logr"
	util "github.com/kubbee/ccloud-cluster-operator/internal"
)

func newEnvironment(environmentName string, logger *logr.Logger) (*util.Environment, error) {
	logger.Info("Start::NewConfluentEnvironment")

	cmd := exec.Command("/bin/confluent", "environment", "create", environmentName, "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		logger.Error(err, "error to create environment")
		return &util.Environment{}, err
	} else {
		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		env := &util.Environment{}

		if err := json.Unmarshal([]byte(message), env); err != nil {
			logger.Error(err, "error to parse the environment json to struct")
			return &util.Environment{}, err
		}

		return env, nil
	}
}

func findEnviroment(environmentName string, logger *logr.Logger) (*util.Environment, error) {
	logger.Info("Start::getEnvironment")
	cmd := exec.Command("/bin/confluent", "environment", "list", "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		logger.Error(err, "Error to retrive clonfluent cloud environments")
		return &util.Environment{}, err
	} else {
		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		environments := []util.Environment{}

		if err := json.Unmarshal([]byte(message), &environments); err != nil {
			logger.Error(err, "error to parse environment json to struct")
			return &util.Environment{}, err
		}

		environment := util.Environment{}

		for i := 0; i < len(environments); i++ {
			if environments[i].Name == environmentName {
				environment.Id = environments[i].Id
				environment.Name = environments[i].Name
				break
			}
		}

		return &environment, nil
	}
}

func setEnvironment(environmentId string, logger *logr.Logger) bool {
	logger.Info("Start::setEnvironment")
	logger.Info("environemtId >> " + environmentId)

	cmd := exec.Command("/bin/confluent", "environment", "use", environmentId)

	if err := cmd.Run(); err != nil {
		logger.Error(err, "error was not possible select the environment")
		return false
	}

	return true
}
