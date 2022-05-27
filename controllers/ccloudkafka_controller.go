/*
Copyright 2022 Kubbee Tech.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	messagesv1alpha1 "github.com/kubbee/ccloud-cluster-operator/api/v1alpha1"
)

// CCloudKafkaReconciler reconciles a CCloudKafka object
type CCloudKafkaReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Log      logr.Logger
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=messages.kubbee.tech,resources=ccloudkafkas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=messages.kubbee.tech,resources=ccloudkafkas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=messages.kubbee.tech,resources=ccloudkafkas/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CCloudKafka object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *CCloudKafkaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)

	ccloudKafka := &messagesv1alpha1.CCloudKafka{}

	if err := r.Get(ctx, req.NamespacedName, ccloudKafka); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	if !ccloudKafka.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("Deleting")
		return reconcile.Result{}, nil // implementing the nil in the future
	}

	if ccloudKafka.Spec.Availability != "" {
		fmt.Println("hello")
	}

	return r.declareKafka(ctx, req, ccloudKafka)
}

func (r *CCloudKafkaReconciler) declareKafka(ctx context.Context, req ctrl.Request, ccloudKafka *messagesv1alpha1.CCloudKafka) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CCloudKafkaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&messagesv1alpha1.CCloudKafka{}).
		Complete(r)
}
