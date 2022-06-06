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
	"time"

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
		if k8sErrors.IsNotFound(err) {
			logger.Info("KafkaResource Not Found.")

			if !ccloudKafka.ObjectMeta.DeletionTimestamp.IsZero() {
				logger.Info("Was marked for deletion.")
				return reconcile.Result{}, nil // implementing the nil in the future
			}
		}
		return reconcile.Result{}, nil
	}

	secretName := "kafka-" + ccloudKafka.Spec.Environment

	foundSecret := &corev1.Secret{}

	err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: req.Namespace}, foundSecret)

	/*
	 * if the secret not exist on the namespace we call
	 */
	if err != nil && k8sErrors.IsNotFound(err) {

		logger.Info("call method to create the environment.")
		return r.declareKafka(ctx, req, ccloudKafka)

	}

	return reconcile.Result{}, nil
}

/**
 * This function Get the Kafka Cluster Settings and verify if the cluster was provisioned, if not wait 5 secods and try again;
 */
func kafkaCreationStatus(id string, environmentId string, logger *logr.Logger) bool {
	logger.Info("Start::kafkaCreationStatus")
	logger.Info("ClusterId >>>> " + id)

	if kafkaClusterSettings, ee := services.GetKafkaClusterSettings(id, environmentId, logger); ee == nil {

		if kafkaClusterSettings.Status == "PROVISIONING" {
			logger.Info("The status is PROVISIONING")

			time.Sleep(5 * time.Second)

			return false
		} else {

			logger.Info("The Kafka Cluster was PROVISIONED")
			return true
		}

	} else {
		logger.Info("Not fount the cluster yet.")
		return false
	}
}

/**
 * This function creates the Kafka Cluster on the Confluent Cloud, and it is responsible do call another 2 important functions
 *
 * <declareKafkaApiKey>
 * <declareKafkaSecret>
 *
 */
func (r *CCloudKafkaReconciler) declareKafka(ctx context.Context, req ctrl.Request, ccloudKafka *messagesv1alpha1.CCloudKafka) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Start::declareKafka")

	connectionCreds, cCerr := r.readCredentials(ctx, req.Namespace, ccloudKafka.Spec.Environment)

	if cCerr != nil {
		logger.Error(cCerr, "Error to get environment secret from cluster")
	}

	if envId, isEnvIdOk := connectionCreds.Data("environmentId"); isEnvIdOk {

		kafka, bErr := services.BuildKafka(ccloudKafka, string(envId), &logger)

		if bErr != nil {
			logger.Error(bErr, "Error to Build Cluster Kafka")
			return reconcile.Result{}, bErr
		}

		for {
			logger.Info("Check if kafka cluster are ready.")

			if kafkaCreationStatus(kafka.Id, string(envId), &logger) {

				logger.Info("Prepare to create Kafka Api-Key")
				logger.Info("kafkaId >>>>>>>>>>> " + kafka.Id)

				if kafkaApiKey, e := r.declareKafkaApiKey(ctx, kafka, ccloudKafka); e == nil {
					if kafkaClusterSettings, ee := services.GetKafkaClusterSettings(kafka.Id, string(envId), &logger); ee == nil {

						if secret, eee := r.declareKafkaSecret(ctx, req, ccloudKafka.Spec.Environment, kafkaClusterSettings, kafkaApiKey); eee == nil {

							logger.Info("Run command to create Secret")

							if eeee := r.Create(ctx, secret); eeee != nil {
								return reconcile.Result{}, eeee
							} else {
								return reconcile.Result{}, nil
							}

						} else {
							return reconcile.Result{}, eee
						}

					} else {
						return reconcile.Result{}, ee
					}

				} else {
					return reconcile.Result{}, e
				}
			}
		}
	}

	return ctrl.Result{}, nil
}

/**
 * This function creates an Api-Key for the Kafka Cluster
 */
func (r *CCloudKafkaReconciler) declareKafkaApiKey(ctx context.Context, kafka *util.ClusterKafka,
	ccloudKafka *messagesv1alpha1.CCloudKafka) (*util.ApiKey, error) {

	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Start::declareKafkaApiKey")

	kafkaApiKey, cErr := services.CreateKafkaApiKey(kafka.Id, ccloudKafka.Spec.ApiKeyName, &logger)

	if cErr != nil {
		logger.Error(cErr, "Error to get Kafka Cluster Settings")
		return &util.ApiKey{}, cErr
	}

	return kafkaApiKey, nil
}

/**
 * This function creates a secret with Api-Key and Endpoint Kafka
 */
func (r *CCloudKafkaReconciler) declareKafkaSecret(ctx context.Context, req ctrl.Request, environment string,
	kafkaClusterSettings *util.ClusterKafka, apiKey *util.ApiKey) (*corev1.Secret, error) {

	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Start::declareKafkaSecret")

	var labels = make(map[string]string)
	labels["name"] = kafkaClusterSettings.Id
	labels["owner"] = "ccloud-cluster-operator"
	labels["controller"] = "ccloudkafka_controller"

	var secretName string = "kafka-" + environment

	var immutable bool = true

	// create and return secret object.
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: req.Namespace,
			Labels:    labels,
		},
		Type:      "kubbee.tech/secret",
		Data:      map[string][]byte{"KafkaEndpoint": []byte(kafkaClusterSettings.Endpoint), "ApiKey": []byte(apiKey.Api), "ApiSecret": []byte(apiKey.Secret)},
		Immutable: &immutable,
	}, nil
}

/**
 * This function is responsible to get the Environment Secret on the namespace;
 */
func (r *CCloudKafkaReconciler) readCredentials(ctx context.Context, requestNamespace string, secretName string) (util.ConnectionCredentials, error) {
	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Read credentials from cluster")

	secret := &corev1.Secret{}

	if err := r.Get(ctx, types.NamespacedName{Namespace: requestNamespace, Name: secretName}, secret); err != nil {
		return nil, err
	}

	return r.readCredentialsFromKubernetesSecret(secret), nil
}

/**
 * This function is responsible to ready the content of the secret and return for execute operations with the data
 */
func (r *CCloudKafkaReconciler) readCredentialsFromKubernetesSecret(secret *corev1.Secret) *util.ClusterCredentials {
	return &util.ClusterCredentials{
		DataContent: map[string][]byte{
			"environmentName": secret.Data["environmentName"],
			"environmentId":   secret.Data["environmentId"],
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *CCloudKafkaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&messagesv1alpha1.CCloudKafka{}).
		Complete(r)
}
