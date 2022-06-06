package services

import (
	"github.com/go-logr/logr"
	util "github.com/kubbee/ccloud-cluster-operator/internal"

	messagesv1alpha1 "github.com/kubbee/ccloud-cluster-operator/api/v1alpha1"
)

func BuildKafka(kafka *messagesv1alpha1.CCloudKafka, environmentId string, logger *logr.Logger) (*util.ClusterKafka, error) {
	logger.Info("Start::BuildKafka")
	if setEnvironment(environmentId, logger) {
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

func CreateKafkaApiKey(clusterId string, description string, logger *logr.Logger) (*util.ApiKey, error) {
	return newKafkaApiKey(clusterId, description, logger)
}

func GetKafkaClusterSettings(clusterId string, environmentId string, logger *logr.Logger) (*util.ClusterKafka, error) {
	logger.Info("Start::GetKafkaClusterSettings")
	logger.Info("ClusterId >>>> " + clusterId)
	if setEnvironment(environmentId, logger) {
		return findKafkaClusterSettings(clusterId, logger)
	} else {
		logger.Info("Was not possible select the environment check the name.")
		return &util.ClusterKafka{}, nil
	}
}

func createKafka(kafka *messagesv1alpha1.CCloudKafka, logger *logr.Logger) (*util.ClusterKafka, error) {
	logger.Info("Start::createKafka")
	return newKafka(kafka, logger)
}

func getKafka(kafka *messagesv1alpha1.CCloudKafka, logger *logr.Logger) (*util.ClusterKafka, error) {
	logger.Info("Start::getKafka")
	return findKafka(kafka.Spec.ClusterName, logger)
}
