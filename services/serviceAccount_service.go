package services

import (
	"bytes"
	"encoding/json"
	"os/exec"

	util "github.com/kubbee/ccloud-cluster-operator/internal"
)

// createServiceAccount this func create a service account into confluent cloud envrionment
func createServiceAccount(sa, description string) (*util.ServiceAccount, error) {
	//confluent iam service-account create CadatralServiceAccount --description "This is a text"

	cmd := exec.Command("/bin/confluent", "iam", "service-account", "create", sa, "--description", description, "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		return &util.ServiceAccount{}, err
	} else {

		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		serviceAccount := util.ServiceAccount{}

		json.Unmarshal([]byte(message), &serviceAccount)

		return &serviceAccount, nil
	}

}
