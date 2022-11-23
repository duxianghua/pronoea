/*
Copyright 2022.

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
	"reflect"
	"time"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	v1 "github.com/duxianghua/pronoea/internal/api/v1"
	operatorV1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/rs/zerolog/log"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	Probe       *ProbeReconciler
	someIndexer client.FieldIndexer
)

const (
	BlackboxConfigMapName  = "pronoea-blackbox-configmap"
	BlackboxDeploymentName = "pronoea-blackbox"
	BlackboxServiceName    = "pronoea-blackbox"
)

// ProbeReconciler reconciles a Probe object
type ProbeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	cfg    Config
}

type Config struct {
	Modules map[string]v1.Module `yaml:"modules" json:"modules"`
}

//+kubebuilder:rbac:groups=pronoea.io,resources=probes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pronoea.io,resources=probes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pronoea.io,resources=probes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Probe object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ProbeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Info().Msg("ProbeReconciler")
	probe := v1.Probe{}
	err := r.Get(ctx, req.NamespacedName, &probe)
	if errors.IsNotFound(err) || !probe.Spec.Pause {
		r.delPromProbe(ctx, req.Namespace, req.Name)
	} else if err != nil {
		log.Error().Err(err).Msg("get probe fail")
		return ctrl.Result{}, err
	} else {
		r.syncPromProbe(ctx, probe)
	}

	//
	probeList := v1.ProbeList{}
	err = r.List(ctx, &probeList, &client.ListOptions{})
	if err != nil {
		return ctrl.Result{}, nil
	}

	modules := make(map[string]v1.Module)
	for _, probe := range probeList.Items {
		modules[probe.ObjectMeta.Name] = probe.Spec.Module
	}
	cfg := Config{Modules: modules}
	// 如果cfg没有改变则退出
	if !reflect.DeepEqual(cfg, r.cfg) {
		yamlByte, err := yaml.Marshal(cfg)
		if err != nil {
			log.Error().Msg(err.Error())
			return ctrl.Result{}, nil
		}
		if err = r.createOrUpdateBlackboxConfigmap(ctx, req.Namespace, req.Name, map[string]string{"blackbox.yml": string(yamlByte)}); err != nil {
			log.Error().Msg(err.Error())
			return ctrl.Result{}, nil
		}
		log.Info().Msg("update configmap")
		err = r.createBlackbox(ctx, req.Namespace)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		// if _, err := http.Get(fmt.Sprintf("http://%s.%s.svc.cluster.local:9115/-/reload", BlackboxServiceName, probe.ObjectMeta.Namespace)); err != nil {
		// 	log.Error().Msg(err.Error())
		// }
		r.reloadBlackbox(ctx, probe.ObjectMeta.Namespace)
	}
	r.cfg = cfg
	// TODO(user): your logic here
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProbeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	Probe = r
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Probe{}).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(de event.DeleteEvent) bool {
				ctrl.Log.Info("deleteEvent", de.Object.GetName())
				return true
			},
			CreateFunc: func(ce event.CreateEvent) bool {
				ctrl.Log.Info("createEvent", ce.Object.GetName())
				return true
			},
		}).
		Complete(r)
}

func (r *ProbeReconciler) createOrUpdateBlackboxConfigmap(ctx context.Context, namespace string, name string, data map[string]string) error {
	configmap := coreV1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      BlackboxConfigMapName,
		}}

	err := r.Get(ctx, client.ObjectKey{Namespace: namespace, Name: BlackboxConfigMapName}, &configmap)
	if errors.IsNotFound(err) {
		configmap.Data = data
		err := r.Create(ctx, &configmap, &client.CreateOptions{})
		if err != nil {
			log.Debug().Err(err).Msg("Create blackbox configmap fail")
		}
		return err
	} else if err != nil {
		log.Error().Err(err).Msg("Get blackbox configmap fail")
		return err
	} else {
		configmap.Data = data
		err = r.Update(ctx, &configmap, &client.UpdateOptions{})
		if err != nil {
			log.Error().Err(err).Msg("Update blackbox configmap fail")
		}
		return err
	}

}

func (r *ProbeReconciler) reloadBlackbox(ctx context.Context, namespace string) {
	deployment := appsV1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: BlackboxDeploymentName}, &deployment)
	if err != nil {
		log.Error().Err(err).Msg("reloadblackbox: get deployment fail")
		return
	}

	if len(deployment.Spec.Template.GetAnnotations()) != 0 {
		deployment.Spec.Template.ObjectMeta.Annotations["pronoea.io/restartedAt"] = time.Now().String()
	} else {
		deployment.Spec.Template.ObjectMeta.Annotations = map[string]string{"pronoea.io/restartedAt": time.Now().String()}
	}
	err = r.Update(ctx, &deployment)
	if err != nil {
		log.Error().Err(err).Msg("reloadblackbox: update deployment fail")
	}
}

func (r *ProbeReconciler) createBlackbox(ctx context.Context, namespace string) error {
	podLabels := map[string]string{
		"app": "blackbox",
	}

	blackboxDeployment := appsV1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      BlackboxDeploymentName,
		},
		Spec: appsV1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: podLabels},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: podLabels,
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Name:  "blackbox",
							Image: "bitnami/blackbox-exporter:latest",
							Ports: []coreV1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: int32(9115),
									Protocol:      coreV1.ProtocolTCP,
								},
							},
							VolumeMounts: []coreV1.VolumeMount{
								{
									Name:      BlackboxConfigMapName,
									MountPath: "/opt/bitnami/blackbox-exporter/blackbox.yml",
									SubPath:   "blackbox.yml",
								},
							},
						},
					},
					Volumes: []coreV1.Volume{
						{
							Name: BlackboxConfigMapName,
							VolumeSource: coreV1.VolumeSource{
								ConfigMap: &coreV1.ConfigMapVolumeSource{
									LocalObjectReference: coreV1.LocalObjectReference{
										Name: BlackboxConfigMapName,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	blackboxSVC := coreV1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      BlackboxServiceName,
			Namespace: namespace,
			Labels:    podLabels,
		},
		Spec: coreV1.ServiceSpec{
			Selector: podLabels,
			Ports: []coreV1.ServicePort{
				{
					Name:       "http",
					Port:       int32(9115),
					Protocol:   coreV1.ProtocolTCP,
					TargetPort: intstr.IntOrString{IntVal: int32(9115)},
				},
			},
		},
	}
	err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: BlackboxDeploymentName}, &appsV1.Deployment{})
	if errors.IsNotFound(err) {
		r.Create(ctx, &blackboxDeployment, &client.CreateOptions{})
	} else if err != nil {
		log.Error().Err(err).Msg("Get blackbox deployment fail")
		return err
	}

	err = r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: BlackboxServiceName}, &coreV1.Service{})
	if errors.IsNotFound(err) {
		err = r.Create(ctx, &blackboxSVC)
		if err != nil {
			log.Error().Err(err).Msg("Create blackbox SVC fail")
		}
	} else if err != nil {
		log.Error().Err(err).Msg("Get blackbox SVC fail")
		return err
	}
	return nil
}

func (r *ProbeReconciler) syncPromProbe(ctx context.Context, probe v1.Probe) {
	promProbe := operatorV1.Probe{
		ObjectMeta: metav1.ObjectMeta{
			Name:      probe.ObjectMeta.Name,
			Namespace: probe.ObjectMeta.Namespace,
			Labels: map[string]string{
				"release": "prometheus-operator",
			},
		},
		Spec: operatorV1.ProbeSpec{
			JobName: probe.ObjectMeta.Name,
			Module:  probe.ObjectMeta.Name,
			Targets: operatorV1.ProbeTargets{
				StaticConfig: &operatorV1.ProbeTargetStaticConfig{
					Targets: probe.Spec.Targets,
					Labels:  probe.Labels,
				},
			},
			ProberSpec: operatorV1.ProberSpec{
				Scheme: "http",
				URL:    fmt.Sprintf("%s.%s.svc.cluster.local:9115", BlackboxServiceName, probe.Namespace),
			},
			Interval:      operatorV1.Duration("1m"),
			ScrapeTimeout: operatorV1.Duration("30s"),
		},
	}
	found := operatorV1.Probe{}
	err := r.Get(ctx, client.ObjectKey{Namespace: probe.Namespace, Name: probe.Name}, &found)
	if err != nil && errors.IsNotFound(err) {
		if err = r.Create(ctx, &promProbe, &client.CreateOptions{}); err != nil {
			log.Error().Interface("err", err).Msg("Failed to create the probe")
		}
		log.Info().Msg("success to Create new probe")
		return
	} else if err != nil {
		log.Error().Err(err).Msg("Failed to Update the probe")
		return
	} else {
		promProbe.ResourceVersion = found.ResourceVersion
		if err = r.Update(ctx, &promProbe, &client.UpdateOptions{}); err != nil {
			log.Error().Err(err).Msg("Failed to Update the prometheus probe")
		}
	}
	err = r.Update(ctx, &probe)
	if err != nil {
		log.Error().Err(err).Msg("Failed to Update the probe status")
	}

}

func (r *ProbeReconciler) delPromProbe(ctx context.Context, namespace string, name string) {
	promProbe := operatorV1.Probe{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	err := r.Delete(ctx, &promProbe)
	if err != nil && !errors.IsNotFound(err) {
		log.Error().Err(err).Msg("Failed to Update the probe")
	}
}
