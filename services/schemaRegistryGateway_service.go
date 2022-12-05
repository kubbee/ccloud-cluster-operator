package services

import (
	"github.com/go-logr/logr"
	messagesv1alpha1 "github.com/kubbee/ccloud-cluster-operator/api/v1alpha1"
	util "github.com/kubbee/ccloud-cluster-operator/internal"
)

/*
 *
 */
func BuildSechameRegistry(ccloudSR *messagesv1alpha1.CCloudSchemaRegistry, environmentId string, logger *logr.Logger) (*util.SchemaRegistry, error) {
	logger.Info("Building schema registry")
	if setEnvironment(environmentId) {
		if ccloudSR.Spec.CCloudSchemaRegistryResource.ResourceExist {
			return getSechameRegistry(environmentId, logger)
		} else {
			return createSechameRegistry(ccloudSR, environmentId, logger)
		}
	} else {
		logger.Info("Was not possible select the environment check the name.")
		return &util.SchemaRegistry{}, nil
	}
}

/*
 *
 */
func createSechameRegistry(ccloudSR *messagesv1alpha1.CCloudSchemaRegistry, environmentId string, logger *logr.Logger) (*util.SchemaRegistry, error) {
	logger.Info("Creating schema registry")
	return newSR(ccloudSR.Spec.Cloud, ccloudSR.Spec.Geo, environmentId, logger)
}

/*
 *
 */
func getSechameRegistry(environmentId string, logger *logr.Logger) (*util.SchemaRegistry, error) {
	logger.Info("Getting schema registry")
	return getSR(environmentId, logger)
}

/*
 *
 */
func CreateSRApiKey(sr *util.SchemaRegistry, apiKeyName string, serviceAccount string) (*util.ApiKey, error) {
	return newSRApiKey(sr.Id, apiKeyName, serviceAccount)
}
