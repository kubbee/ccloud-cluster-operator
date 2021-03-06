package services

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strconv"

	"github.com/go-logr/logr"
	messagesv1alpha1 "github.com/kubbee/ccloud-cluster-operator/api/v1alpha1"
	util "github.com/kubbee/ccloud-cluster-operator/internal"
)

func newKafka(kafka *messagesv1alpha1.CCloudKafka, logger *logr.Logger) (*util.ClusterKafka, error) {
	//  confluent kafka cluster create lion  --cloud aws --region us-east-1 --availability single-zone --type basic
	logger.Info("Start::newKafka")

	if id, err := selectKafka(kafka.Spec.ClusterName, logger); err != nil {
		logger.Error(err, "error to find kafka cluster")
		return &util.ClusterKafka{}, err
	} else {
		if id != "" {
			return findKafka(id, logger)
		} else {

			cmd := newKafkaClusterType(kafka, logger)

			cmdOutput := &bytes.Buffer{}
			cmd.Stdout = cmdOutput

			if err := cmd.Run(); err != nil {
				logger.Error(err, "Error to create kafka cluster")
				return &util.ClusterKafka{}, err
			} else {

				output := cmdOutput.Bytes()
				message, _ := getOutput(output)

				clusterKafka := &util.ClusterKafka{}

				if jErr := json.Unmarshal([]byte(message), clusterKafka); jErr != nil {
					logger.Error(jErr, "Error to parse Cluster Kafka")
				}

				return clusterKafka, nil
			}
		}
	}
}

func newKafkaClusterType(kafka *messagesv1alpha1.CCloudKafka, logger *logr.Logger) *exec.Cmd {
	logger.Info("Start::newKafkaType")
	if kafka.Spec.CCloudKafkaDedicate.Dedicated {
		logger.Info("Cluster Type Dedicated")

		ckus := strconv.FormatInt(kafka.Spec.CCloudKafkaDedicate.CKU, 10)

		logger.Info("command,  /bin/confluent kafka cluster create " + kafka.Spec.ClusterName + " --cloud " + kafka.Spec.Cloud + " --region " + kafka.Spec.Region + " --availability " + kafka.Spec.Availability + " --cku " + ckus + " --type " + kafka.Spec.ClusterType + " --output json")

		return exec.Command("/bin/confluent", "kafka", "cluster", "create", kafka.Spec.ClusterName, "--cloud", kafka.Spec.Cloud, "--region", kafka.Spec.Region, "--availability", kafka.Spec.Availability, "--cku", ckus, "--type", kafka.Spec.ClusterType, "--output", "json")
	} else {
		logger.Info("Cluster Type Normal")
		logger.Info("command,  /bin/confluent kafka cluster create " + kafka.Spec.ClusterName + " --cloud " + kafka.Spec.Cloud + " --region " + kafka.Spec.Region + " --availability " + kafka.Spec.Availability + " --type " + kafka.Spec.ClusterType + " --output json")

		return exec.Command("/bin/confluent", "kafka", "cluster", "create", kafka.Spec.ClusterName, "--cloud", kafka.Spec.Cloud, "--region", kafka.Spec.Region, "--availability", kafka.Spec.Availability, "--type", kafka.Spec.ClusterType, "--output", "json")
	}
}

func findKafka(kafkaClusteName string, logger *logr.Logger) (*util.ClusterKafka, error) {
	logger.Info("Start::findKafka")
	logger.Info("Run::Command::confluent kafka cluster list --output json")

	cmd := exec.Command("/bin/confluent", "kafka", "cluster", "list", "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		logger.Error(err, "Error to get kafka cluster")
		return &util.ClusterKafka{}, err
	} else {

		var clusterId string

		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		clusters := []util.ClusterKafka{}

		if jErr := json.Unmarshal([]byte(message), &clusters); jErr != nil {
			logger.Error(jErr, "Error to parse kafka cluster")
		}

		for i := 0; i < len(clusters); i++ {
			if kafkaClusteName == clusters[i].Name {
				clusterId = clusters[i].Id
				break
			}
		}

		return findKafkaClusterSettings(clusterId, logger)
	}
}

func selectKafka(kafkaClusteName string, logger *logr.Logger) (string, error) {
	logger.Info("Start::selectKafka")
	logger.Info("Run::Command::confluent kafka cluster list --output json")

	cmd := exec.Command("/bin/confluent", "kafka", "cluster", "list", "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		logger.Error(err, "Error to get kafka cluster")
		return "", err
	} else {

		var clusterId string

		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		clusters := []util.ClusterKafka{}

		if jErr := json.Unmarshal([]byte(message), &clusters); jErr != nil {
			logger.Error(jErr, "Error to parse kafka cluster")
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

func findKafkaClusterSettings(clusterId string, logger *logr.Logger) (*util.ClusterKafka, error) {
	logger.Info("Start::getKafkaClusterSettings")
	logger.Info("command:: /bin/confluent kafka cluster use " + clusterId)

	cmdUse := exec.Command("/bin/confluent", "kafka", "cluster", "use", clusterId)

	if err := cmdUse.Run(); err != nil {
		logger.Error(err, "error to use Kafka Cluster")
		return &util.ClusterKafka{}, err
	} else {

		cmd := exec.Command("/bin/confluent", "kafka", "cluster", "describe", "--output", "json")

		cmdOutput := &bytes.Buffer{}
		cmd.Stdout = cmdOutput

		if err := cmd.Run(); err != nil {
			logger.Error(err, "error to get the kafka cluster reference")
			return &util.ClusterKafka{}, err
		} else {
			output := cmdOutput.Bytes()
			message, _ := getOutput(output)

			cluster := &util.ClusterKafka{}

			if jErr := json.Unmarshal([]byte(message), cluster); jErr != nil {
				logger.Error(jErr, "error to parse kafka cluster")
				return &util.ClusterKafka{}, jErr
			}

			return cluster, nil
		}
	}
}

func newKafkaApiKey(clusterId string, description string, logger *logr.Logger) (*util.ApiKey, error) {
	logger.Info("Creating Api-Key for Kafka Connection")

	cmd := exec.Command("/bin/confluent", "api-key", "create", "--resource", clusterId, "--description", description, "--output", "json")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		logger.Error(err, "Confluent::NewKafkaApiKey:: Error to create the kafka api-key for the application")
		return &util.ApiKey{}, err
	} else {

		output := cmdOutput.Bytes()
		message, _ := getOutput(output)

		kafkaApiKey := &util.ApiKey{}

		if jErr := json.Unmarshal([]byte(message), kafkaApiKey); jErr != nil {
			logger.Error(jErr, "Error to parse Kafka ApiKey")
			return &util.ApiKey{}, jErr
		}

		return kafkaApiKey, nil

	}
}
