package services

import (
	"bytes"
	"encoding/json"
	"os/exec"

	util "github.com/kubbee/ccloud-cluster-operator/internal"
)

func newEnvironment(environmentName string) (*util.Environment, error) {
	cmd := exec.Command("/bin/confluent", "environment", "create", environmentName, "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		return &util.Environment{}, err
	} else {
		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		env := &util.Environment{}

		if err := json.Unmarshal([]byte(message), env); err != nil {
			return &util.Environment{}, err
		}

		return env, nil
	}
}

func findEnviroment(environmentName string) (*util.Environment, error) {
	cmd := exec.Command("/bin/confluent", "environment", "list", "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		return &util.Environment{}, err
	} else {
		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		environments := []util.Environment{}

		if err := json.Unmarshal([]byte(message), &environments); err != nil {
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

func setEnvironment(environmentId string) bool {
	cmd := exec.Command("/bin/confluent", "environment", "use", environmentId)

	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
