package services

import (
	"github.com/go-logr/logr"
	messagesv1alpha1 "github.com/kubbee/ccloud-cluster-operator/api/v1alpha1"
	util "github.com/kubbee/ccloud-cluster-operator/internal"
)

func CreateEnvironment(ccloudEnvironment *messagesv1alpha1.CCloudEnvironment, logger *logr.Logger) (*util.Environment, error) {
	return newConfluentEnvironment(ccloudEnvironment.Spec.Name, logger)
}
