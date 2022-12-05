package services

import (
	util "github.com/kubbee/ccloud-cluster-operator/internal"
)

func CreateServiceAccount(sa string, description string) (*util.ServiceAccount, error) {
	return createServiceAccount(sa, description)
}
