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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	messagesv1alpha1 "github.com/kubbee/ccloud-cluster-operator/api/v1alpha1"
	util "github.com/kubbee/ccloud-cluster-operator/internal"
	services "github.com/kubbee/ccloud-cluster-operator/services"

	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
)

// CCloudEnvironmentReconciler reconciles a CCloudEnvironment object
type CCloudEnvironmentReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Log      logr.Logger
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=messages.kubbee.tech,resources=ccloudenvironments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=messages.kubbee.tech,resources=ccloudenvironments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=messages.kubbee.tech,resources=ccloudenvironments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CCloudEnvironment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *CCloudEnvironmentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)

	ccloudEnvironment := &messagesv1alpha1.CCloudEnvironment{}

	if err := r.Get(ctx, req.NamespacedName, ccloudEnvironment); err != nil {
		if k8sErrors.IsNotFound(err) {
			logger.Info("KafkaResource Not Found.")

			if !ccloudEnvironment.ObjectMeta.DeletionTimestamp.IsZero() {
				logger.Info("Was marked for deletion.")
				return reconcile.Result{}, nil // implementing the nil in the future
			}
		}
		return reconcile.Result{}, nil
	}

	//definition of secret name
	secretName := ccloudEnvironment.Spec.Name

	//type definition
	foundSecret := &corev1.Secret{}

	// try find the secret, and if not find create call the process to declare the environment
	err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: req.Namespace}, foundSecret)

	if err != nil && k8sErrors.IsNotFound(err) {

		logger.Info("call method to create the environment.")
		return r.declareEnvironment(ctx, req, ccloudEnvironment)
	}

	return reconcile.Result{}, nil
}

/**
 * This function create an environment on the confluent cloud
 */
func (r *CCloudEnvironmentReconciler) declareEnvironment(ctx context.Context, req ctrl.Request, environment *messagesv1alpha1.CCloudEnvironment) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)

	if env, err := services.BuildEnvironment(environment, &logger); err != nil {
		logger.Error(err, "Error to create environment")
		return reconcile.Result{}, err
	} else {

		sa, err := r.declareServiceAccountForEnvironment(ctx, req, environment)

		if err != nil {
			logger.Error(err, "was not possible create service account")
		}

		secret := r.declareEnvironmentSecret(req, env, sa)

		if err := r.Create(ctx, secret); err != nil {
			logger.Error(err, "error to create environment secret")
			return reconcile.Result{}, nil
		} else {
			return reconcile.Result{}, nil
		}
	}
}

/**
 * This function creates an secret environment on the namespace
 */
func (r *CCloudEnvironmentReconciler) declareEnvironmentSecret(req ctrl.Request, environment *util.Environment, sa *util.ServiceAccount) *corev1.Secret {

	var labels = make(map[string]string)
	labels["name"] = environment.Name
	labels["owner"] = "ccloud-messaging-topology-operator"
	labels["controller"] = "ccloudenvironment_controller"

	var immutable bool = true

	// create and return secret object.
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      environment.Name,
			Namespace: req.Namespace,
			Labels:    labels,
		},
		Type:      "kubbee.tech/secret",
		Data:      map[string][]byte{"environmentName": []byte(environment.Name), "environmentId": []byte(environment.Id), "serviceAccount": []byte(sa.Id)},
		Immutable: &immutable,
	}
}

/**
 * This function declare a Service Account with the same name of environment
 */
func (r *CCloudEnvironmentReconciler) declareServiceAccountForEnvironment(ctx context.Context, req ctrl.Request, environment *messagesv1alpha1.CCloudEnvironment) (*util.ServiceAccount, error) {
	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Declaring Service Account for environment " + environment.Name)

	return services.CreateServiceAccount(environment.Name, environment.Namespace+"-"+environment.Name)
}

// SetupWithManager sets up the controller with the Manager.
func (r *CCloudEnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&messagesv1alpha1.CCloudEnvironment{}).
		Complete(r)
}
