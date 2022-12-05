package services

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strconv"

	messagesv1alpha1 "github.com/kubbee/ccloud-cluster-operator/api/v1alpha1"
	util "github.com/kubbee/ccloud-cluster-operator/internal"
)

func newKafka(kafka *messagesv1alpha1.CCloudKafka) (*util.ClusterKafka, error) {
	//  confluent kafka cluster create lion  --cloud aws --region us-east-1 --availability single-zone --type basic

	if id, err := selectKafka(kafka.Spec.ClusterName); err != nil {
		return &util.ClusterKafka{}, err
	} else {
		if id != "" {
			return findKafka(id)
		} else {

			cmd := newKafkaClusterType(kafka)

			cmdOutput := &bytes.Buffer{}
			cmd.Stdout = cmdOutput

			if err := cmd.Run(); err != nil {
				return &util.ClusterKafka{}, err
			} else {

				output := cmdOutput.Bytes()
				message, _ := getOutput(output)

				clusterKafka := &util.ClusterKafka{}

				if jErr := json.Unmarshal([]byte(message), clusterKafka); jErr != nil {
					return &util.ClusterKafka{}, jErr
				}

				return clusterKafka, nil
			}
		}
	}
}

func newKafkaClusterType(kafka *messagesv1alpha1.CCloudKafka) *exec.Cmd {
	if kafka.Spec.CCloudKafkaDedicate.Dedicated {
		ckus := strconv.FormatInt(kafka.Spec.CCloudKafkaDedicate.CKU, 10)
		return exec.Command("/bin/confluent", "kafka", "cluster", "create", kafka.Spec.ClusterName, "--cloud", kafka.Spec.Cloud, "--region", kafka.Spec.Region, "--availability", kafka.Spec.Availability, "--cku", ckus, "--type", kafka.Spec.ClusterType, "--output", "json")
	} else {
		return exec.Command("/bin/confluent", "kafka", "cluster", "create", kafka.Spec.ClusterName, "--cloud", kafka.Spec.Cloud, "--region", kafka.Spec.Region, "--availability", kafka.Spec.Availability, "--type", kafka.Spec.ClusterType, "--output", "json")
	}
}

func findKafka(kafkaClusteName string) (*util.ClusterKafka, error) {
	cmd := exec.Command("/bin/confluent", "kafka", "cluster", "list", "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		return &util.ClusterKafka{}, err
	} else {

		var clusterId string

		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		clusters := []util.ClusterKafka{}

		if jErr := json.Unmarshal([]byte(message), &clusters); jErr != nil {
			return &util.ClusterKafka{}, jErr
		}

		for i := 0; i < len(clusters); i++ {
			if kafkaClusteName == clusters[i].Name {
				clusterId = clusters[i].Id
				break
			}
		}
		return findKafkaClusterSettings(clusterId)
	}
}

func selectKafka(kafkaClusteName string) (string, error) {

	cmd := exec.Command("/bin/confluent", "kafka", "cluster", "list", "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		return "", err
	} else {

		var clusterId string

		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		clusters := []util.ClusterKafka{}

		if jErr := json.Unmarshal([]byte(message), &clusters); jErr != nil {
			return "", jErr
		}

		for i := 0; i < len(clusters); i++ {
			if kafkaClusteName == clusters[i].Name {
				clusterId = clusters[i].Id
				break
			}
		}

		return clusterId, nil
	}
}

func findKafkaClusterSettings(clusterId string) (*util.ClusterKafka, error) {
	cmdUse := exec.Command("/bin/confluent", "kafka", "cluster", "use", clusterId)

	if err := cmdUse.Run(); err != nil {
		return &util.ClusterKafka{}, err
	} else {
		cmd := exec.Command("/bin/confluent", "kafka", "cluster", "describe", "--output", "json")

		cmdOutput := &bytes.Buffer{}
		cmd.Stdout = cmdOutput

		if err := cmd.Run(); err != nil {
			return &util.ClusterKafka{}, err
		} else {
			output := cmdOutput.Bytes()
			message, _ := getOutput(output)

			cluster := &util.ClusterKafka{}

			if jErr := json.Unmarshal([]byte(message), cluster); jErr != nil {
				return &util.ClusterKafka{}, jErr
			}

			return cluster, nil
		}
	}
}

func newKafkaApiKey(clusterId string, description string, serviceAccount string) (*util.ApiKey, error) {
	cmd := exec.Command("/bin/confluent", "api-key", "create", "--resource", clusterId, "--description", description, "--service-account", serviceAccount, "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		return &util.ApiKey{}, err
	} else {

		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		kafkaApiKey := &util.ApiKey{}

		if jErr := json.Unmarshal([]byte(message), kafkaApiKey); jErr != nil {
			return &util.ApiKey{}, jErr
		}

		return kafkaApiKey, nil

	}
}
