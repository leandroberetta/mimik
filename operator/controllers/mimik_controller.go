/*
Copyright 2021.

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
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	mimikv1alpha1 "github.com/leandroberetta/mimik-operator/api/v1alpha1"
	"github.com/prometheus/common/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MimikReconciler reconciles a Mimik object
type MimikReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=mimik.veicot.io,resources=mimiks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mimik.veicot.io,resources=mimiks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mimik.veicot.io,resources=mimiks/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Mimik object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *MimikReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("mimik", req.NamespacedName)

	mimik := &mimikv1alpha1.Mimik{}
	err := r.Get(ctx, req.NamespacedName, mimik)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Mimik resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get Mimik")
		return ctrl.Result{}, err
	}

	cm := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Name: getNameWithVersion(mimik.Spec.Service, mimik.Spec.Version), Namespace: mimik.Namespace}, cm)
	if err != nil && errors.IsNotFound(err) {
		cm := r.configMapForMimik(mimik)
		log.Info("Creating a new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
		err = r.Create(ctx, cm)
		if err != nil {
			log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get ConfigMap")
		return ctrl.Result{}, err
	}

	svc := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: mimik.Spec.Service, Namespace: mimik.Namespace}, svc)
	if err != nil && errors.IsNotFound(err) {
		svc := r.serviceForMimik(mimik)
		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.Create(ctx, svc)
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	deploy := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: getNameWithVersion(mimik.Spec.Service, mimik.Spec.Version), Namespace: mimik.Namespace}, deploy)
	if err != nil && errors.IsNotFound(err) {
		deploy := r.deploymentForMimik(mimik)
		log.Info("Creating a new Deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
		err = r.Create(ctx, deploy)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *MimikReconciler) deploymentForMimik(mimik *mimikv1alpha1.Mimik) *appsv1.Deployment {
	labels := labelsForMimik(mimik.Spec.Service, mimik.Spec.Version)
	replicas := int32(1)
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getNameWithVersion(mimik.Spec.Service, mimik.Spec.Version),
			Namespace: mimik.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: map[string]string{"sidecar.istio.io/inject": "true"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "quay.io/leandroberetta/mimik:latest",
						Name:  "mimik",
						Env: []corev1.EnvVar{
							{
								Name:  "MIMIK_SERVICE_NAME",
								Value: mimik.Spec.Service,
							},
							{
								Name:  "MIMIK_SERVICE_PORT",
								Value: "8080",
							},
							{
								Name:  "MIMIK_ENDPOINTS_FILE",
								Value: fmt.Sprintf("/data/%s.json", getNameWithVersion(mimik.Spec.Service, mimik.Spec.Version)),
							},
							{
								Name:  "MIMIK_LABELS_FILE",
								Value: "/tmp/etc/pod_labels",
							},
						},
						ImagePullPolicy: corev1.PullAlways,
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "pod-info",
								MountPath: "/tmp/etc",
							},
							{
								Name:      "endpoints",
								MountPath: "/data",
							},
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "http",
						}},
					}},
					Volumes: []corev1.Volume{
						{
							Name: "pod-info",
							VolumeSource: corev1.VolumeSource{
								DownwardAPI: &corev1.DownwardAPIVolumeSource{
									Items: []corev1.DownwardAPIVolumeFile{
										{
											FieldRef: &corev1.ObjectFieldSelector{
												FieldPath: "metadata.labels",
											},
											Path: "pod_labels",
										},
									},
								},
							},
						},
						{
							Name: "endpoints",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									Items: []corev1.KeyToPath{
										{
											Key:  fmt.Sprintf("%s.json", getNameWithVersion(mimik.Spec.Service, mimik.Spec.Version)),
											Path: fmt.Sprintf("%s.json", getNameWithVersion(mimik.Spec.Service, mimik.Spec.Version)),
										},
									},
									LocalObjectReference: corev1.LocalObjectReference{
										Name: getNameWithVersion(mimik.Spec.Service, mimik.Spec.Version),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return deploy
}

func (r *MimikReconciler) serviceForMimik(mimik *mimikv1alpha1.Mimik) *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mimik.Spec.Service,
			Namespace: mimik.Namespace,
			Labels:    labelsForMimik(mimik.Spec.Service, mimik.Spec.Version),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name: "http",
				Port: 8080,
			},
			},
			Selector: map[string]string{
				"app": mimik.Spec.Service,
			},
		},
	}
	return svc
}

func (r *MimikReconciler) configMapForMimik(mimik *mimikv1alpha1.Mimik) *corev1.ConfigMap {
	jsonData, _ := json.Marshal(mimik.Spec.Endpoints)
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getNameWithVersion(mimik.Spec.Service, mimik.Spec.Version),
			Namespace: mimik.Namespace,
			Labels:    labelsForMimik(mimik.Spec.Service, mimik.Spec.Version),
		},
		Data: map[string]string{fmt.Sprintf("%s.json", getNameWithVersion(mimik.Spec.Service, mimik.Spec.Version)): string(jsonData)},
	}
	return cm
}

func getNameWithVersion(name, version string) string {
	return fmt.Sprintf("%s-%s", name, version)
}

func labelsForMimik(name, version string) map[string]string {
	return map[string]string{"app": name, "version": version}
}

// SetupWithManager sets up the controller with the Manager.
func (r *MimikReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mimikv1alpha1.Mimik{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
