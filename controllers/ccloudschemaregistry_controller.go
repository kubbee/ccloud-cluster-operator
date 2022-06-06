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
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	messagesv1alpha1 "github.com/kubbee/ccloud-cluster-operator/api/v1alpha1"
	util "github.com/kubbee/ccloud-cluster-operator/internal"
	"github.com/kubbee/ccloud-cluster-operator/services"

	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CCloudSchemaRegistryReconciler reconciles a CCloudSchemaRegistry object
type CCloudSchemaRegistryReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Log      logr.Logger
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=messages.kubbee.tech,resources=ccloudschemaregistries,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=messages.kubbee.tech,resources=ccloudschemaregistries/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=messages.kubbee.tech,resources=ccloudschemaregistries/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CCloudSchemaRegistry object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *CCloudSchemaRegistryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)

	ccloudSchemaRegistry := &messagesv1alpha1.CCloudSchemaRegistry{}

	if err := r.Get(ctx, req.NamespacedName, ccloudSchemaRegistry); err != nil {
		if k8sErrors.IsNotFound(err) {
			logger.Info("ScheRegistryResource Not Found.")

			if !ccloudSchemaRegistry.ObjectMeta.DeletionTimestamp.IsZero() {
				logger.Info("Was marked for deletion.")
				return reconcile.Result{}, nil // implementing the nil in the future
			}
		}
		return reconcile.Result{}, err
	}

	//Defines the name of secret
	secretName := "schemaregistry-" + ccloudSchemaRegistry.Spec.Environment

	foundSecret := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: req.Namespace}, foundSecret)

	/*
	 * if the secret not exist on the namespace we call
	 */
	if err != nil && k8sErrors.IsNotFound(err) {
		return r.declareSchemaRegistry(ctx, req, ccloudSchemaRegistry)
	}

	return reconcile.Result{}, nil
}

//
func (r *CCloudSchemaRegistryReconciler) declareSchemaRegistry(ctx context.Context, req ctrl.Request,
	ccloudSR *messagesv1alpha1.CCloudSchemaRegistry) (ctrl.Result, error) {

	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Start::declareSchemaRegistry")

	connectionCreds, cCerr := r.readCredentials(ctx, req.Namespace, ccloudSR.Spec.Environment)

	if cCerr != nil {
		logger.Error(cCerr, "Error to get environment secret from cluster")
	}

	if envId, isEnvIdOk := connectionCreds.Data("environmentId"); isEnvIdOk {

		sr, err := services.BuildSechameRegistry(ccloudSR, string(envId), &logger)

		if err != nil {
			logger.Error(err, "Error to create Schema Registry")
			return reconcile.Result{}, nil
		}

		if apiKey, ee := r.declareSRApiKey(ctx, req, ccloudSR.Spec.ApiKeyName, sr); ee == nil {
			if secret, eee := r.delcareSRSecret(ctx, req, ccloudSR.Spec.Environment, sr, apiKey); eee == nil {
				r.Create(ctx, secret)
			} else {
				return reconcile.Result{}, eee
			}
		} else {
			return reconcile.Result{}, ee
		}
	}

	return reconcile.Result{}, nil
}

func (r *CCloudSchemaRegistryReconciler) declareSRApiKey(ctx context.Context, req ctrl.Request, apiKeyName string, schemaRegistry *util.SchemaRegistry) (*util.ApiKey, error) {
	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Start::declareSchemaRegistry")

	if schemaRegistry.Id == "" {
		return &util.ApiKey{}, errors.New("SchemaRegistry no content")
	}

	apiKey, err := services.CreateSRApiKey(schemaRegistry, apiKeyName, &logger)

	if err != nil {
		logger.Error(err, "Error to create SchemaRegistry ApiKey")
		return &util.ApiKey{}, nil
	}

	return apiKey, nil
}

func (r *CCloudSchemaRegistryReconciler) delcareSRSecret(ctx context.Context, req ctrl.Request, environment string, sr *util.SchemaRegistry, apiKey *util.ApiKey) (*corev1.Secret, error) {
	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Start::declareSchemaRegistrySecret")

	if (apiKey.Api == "" || apiKey.Secret == "") || (sr.Id == "" || sr.EnpointURL == "") {
		return &corev1.Secret{}, errors.New("ApiKey Or SchemaRegistry no content")
	}

	var labels = make(map[string]string)
	labels["name"] = sr.Id
	labels["owner"] = "ccloud-cluster-operator"
	labels["controller"] = "ccloudschemaregistry_controller"

	var secretName string = "schemaregistry-" + environment

	var immutable bool = true

	// create and return secret object.
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: req.Namespace,
			Labels:    labels,
		},
		Type:      "kubbee.tech/secret",
		Data:      map[string][]byte{"SREndpoint": []byte(sr.EnpointURL), "ApiKey": []byte(apiKey.Api), "ApiSecret": []byte(apiKey.Secret)},
		Immutable: &immutable,
	}, nil
}

func (r *CCloudSchemaRegistryReconciler) readCredentials(ctx context.Context, requestNamespace string, secretName string) (util.ConnectionCredentials, error) {
	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Read credentials from cluster")

	secret := &corev1.Secret{}

	if err := r.Get(ctx, types.NamespacedName{Namespace: requestNamespace, Name: secretName}, secret); err != nil {
		return nil, err
	}

	return r.readCredentialsFromKubernetesSecret(secret), nil
}

func (r *CCloudSchemaRegistryReconciler) readCredentialsFromKubernetesSecret(secret *corev1.Secret) *util.ClusterCredentials {
	return &util.ClusterCredentials{
		DataContent: map[string][]byte{
			"environmentName": secret.Data["environmentName"],
			"environmentId":   secret.Data["environmentId"],
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *CCloudSchemaRegistryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&messagesv1alpha1.CCloudSchemaRegistry{}).
		Complete(r)
}
