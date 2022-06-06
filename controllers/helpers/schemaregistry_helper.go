package helpers

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/go-logr/logr"
	"github.com/kubbee/ccloud-cluster-operator/controllers"
)

/*
 *
 */
type SchemaRegistryHelper struct {
	Ctx        context.Context
	Req        ctrl.Request
	Logger     logr.Logger
	Reconciler controllers.CCloudSchemaRegistryReconciler
}

/*
 *
 */
func NewSchemaRegistryHelper(ctx context.Context, req ctrl.Request, logger logr.Logger, r controllers.CCloudSchemaRegistryReconciler) *SchemaRegistryHelper {
	return &SchemaRegistryHelper{
		Ctx:        ctx,
		Req:        req,
		Logger:     logger,
		Reconciler: r,
	}
}
