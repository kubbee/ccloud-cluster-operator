package services

import (
	"github.com/go-logr/logr"
	util "github.com/kubbee/ccloud-cluster-operator/internal"

	messagesv1alpha1 "github.com/kubbee/ccloud-cluster-operator/api/v1alpha1"
)

func BuildKafka(kafka *messagesv1alpha1.CCloudKafka, environmentId string, logger *logr.Logger) (*util.ClusterKafka, error) {
	logger.Info("Building kafka cluster.")
	if setEnvironment(environmentId) {
		if kafka.Spec.CCloudKafkaResource.ResourceExist {
			return getKafka(kafka, logger)
		} else {
			return createKafka(kafka, logger)
		}
	} else {
		logger.Info("Was not possible select the environment check the name.")
		return &util.ClusterKafka{}, nil
	}
}

func CreateKafkaApiKey(clusterId string, description string, serviceAccount string) (*util.ApiKey, error) {
	return newKafkaApiKey(clusterId, description, serviceAccount)
}

func GetKafkaClusterSettings(clusterId string, environmentId string, logger *logr.Logger) (*util.ClusterKafka, error) {
	logger.Info("Getting kafka cluster settings for ClusterId=" + clusterId)
	if setEnvironment(environmentId) {
		return findKafkaClusterSettings(clusterId)
	} else {
		logger.Info("Was not possible select the environment check the name.")
		return &util.ClusterKafka{}, nil
	}
}

func createKafka(kafka *messagesv1alpha1.CCloudKafka, logger *logr.Logger) (*util.ClusterKafka, error) {
	logger.Info("Creating kafka cluster")
	return newKafka(kafka)
}

func getKafka(kafka *messagesv1alpha1.CCloudKafka, logger *logr.Logger) (*util.ClusterKafka, error) {
	logger.Info("Finding kafka cluster")
	return findKafka(kafka.Spec.ClusterName)
}
