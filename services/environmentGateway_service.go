package services

import (
	"github.com/go-logr/logr"
	messagesv1alpha1 "github.com/kubbee/ccloud-cluster-operator/api/v1alpha1"
	util "github.com/kubbee/ccloud-cluster-operator/internal"
)

func BuildEnvironment(ccloudEnvironment *messagesv1alpha1.CCloudEnvironment, logger *logr.Logger) (*util.Environment, error) {
	if ccloudEnvironment.Spec.CCloudEnvironmentResource.ResourceExist {
		return getEnvironment(ccloudEnvironment, logger)
	} else {
		return createEnvironment(ccloudEnvironment, logger)
	}

}

func createEnvironment(ccloudEnvironment *messagesv1alpha1.CCloudEnvironment, logger *logr.Logger) (*util.Environment, error) {
	if environment, e := newEnvironment(ccloudEnvironment.Spec.Name); e == nil {
		if setEnvironment(environment.Id) {
			logger.Info("The environment was selected with sucess")
		}
		return environment, nil
	} else {
		return &util.Environment{}, e
	}
}

func getEnvironment(ccloudEnvironment *messagesv1alpha1.CCloudEnvironment, logger *logr.Logger) (*util.Environment, error) {
	if environment, e := findEnviroment(ccloudEnvironment.Spec.Name); e == nil {
		if setEnvironment(environment.Id) {
			logger.Info("The environment was selected with sucess")
		}
		return environment, nil
	} else {
		return &util.Environment{}, e
	}
}
