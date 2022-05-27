package services

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"regexp"

	"github.com/go-logr/logr"
	util "github.com/kubbee/ccloud-cluster-operator/internal"
)

func newConfluentEnvironment(Environment string, logger *logr.Logger) (*util.Environment, error) {
	logger.Info("Start::NewConfluentEnvironment")

	truth, err := existsResource(exec.Command("/bin/confluent", "environment", "list", "--output", "json"), Environment, logger)

	if !truth {
		cmd := exec.Command("/bin/confluent", "environment", "create", Environment, "--output", "json")

		cmdOutput := &bytes.Buffer{}
		cmd.Stdout = cmdOutput

		if err := cmd.Run(); err != nil {

			output := cmdOutput.Bytes()
			message, _ := getOutput(output)

			env := &util.Environment{}

			json.Unmarshal([]byte(message), env)

			return env, nil
		}
	} else {
		if err != nil {
			logger.Error(err, "Error to verify if the environment already exists")
			return &util.Environment{}, err
		}
	}

	return &util.Environment{}, nil
}

func existsResource(cmd *exec.Cmd, matcher string, logger *logr.Logger) (bool, error) {

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		logger.Error(err, "Error to verify if the environment already exists")
		return false, err
	} else {

		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		return regexp.MatchString("\\b"+matcher+"\\b", message)
	}
}

/*
 * This function get the output console and converts to string
 */
func getOutput(outs []byte) (string, bool) {
	if len(outs) > 0 {
		return string(outs), true
	}
	return "", false
}
