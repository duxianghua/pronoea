package handler

import (
	"context"
	"fmt"

	"github.com/duxianghua/pronoea/internal/controllers"
	"github.com/duxianghua/pronoea/internal/utils"
	"github.com/gin-gonic/gin"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ScenariosAPI struct{}

func (p *ScenariosAPI) List(c *gin.Context) {
	configMapList := coreV1.ConfigMapList{}
	labelSelector := labels.NewSelector()
	managedByRequirement, _ := labels.NewRequirement("app.kubernetes.io/managed-by", selection.Equals, []string{"pronoea"})
	componentRequirement, _ := labels.NewRequirement("app.kubernetes.io/component", selection.Equals, []string{"scenarios"})

	err := controllers.Probe.List(context.TODO(), &configMapList, &client.ListOptions{LabelSelector: labelSelector.Add(*managedByRequirement, *componentRequirement)})
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	// for _, configmap := range configMapList.Items {
	// 	configmap.ObjectMeta.ManagedFields = nil
	// }
	c.JSON(200, configMapList)
}

func (p *ScenariosAPI) Get(c *gin.Context) {
	name := c.Param("name")
	namespace, ok := c.GetQuery("namespace")
	if !ok {
		namespace = utils.GetCurrentNamespace()
	}
	configmap := coreV1.ConfigMap{}
	err := controllers.Probe.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, &configmap)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	configmap.ObjectMeta.ManagedFields = nil
	c.JSON(200, configmap)
}

func (p *ScenariosAPI) Create(c *gin.Context) {
	configmap := coreV1.ConfigMap{}
	name := c.Param("name")
	namespace, ok := c.GetQuery("namespace")
	if !ok {
		namespace = utils.GetCurrentNamespace()
	}
	err := c.ShouldBindJSON(&configmap)
	if err != nil {
		c.Status(400)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	configmap.ObjectMeta.Name = name
	configmap.ObjectMeta.Namespace = namespace
	keys := make([]string, 0, len(configmap.Data))
	for k := range configmap.Data {
		keys = append(keys, k)
	}
	if !contains(keys, "interval") {
		configmap.Data["interval"] = "60"
	}

	err = controllers.Probe.Create(context.TODO(), &configmap, &client.CreateOptions{})
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	configmap.ObjectMeta.ManagedFields = nil
	c.JSON(200, configmap)
}

func (p *ScenariosAPI) Delete(c *gin.Context) {
	name := c.Param("name")
	namespace, ok := c.GetQuery("namespace")
	if !ok {
		namespace = utils.GetCurrentNamespace()
	}
	configmap := coreV1.ConfigMap{}
	configmap.ObjectMeta.Namespace = namespace
	configmap.ObjectMeta.Name = name
	err := controllers.Probe.Delete(c.Request.Context(), &configmap, &client.DeleteAllOfOptions{})
	if err != nil {
		c.Status(404)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	c.JSON(200, gin.H{})
}

func (p *ScenariosAPI) Update(c *gin.Context) {
	configmap := coreV1.ConfigMap{}
	name := c.Param("name")
	err := c.ShouldBindJSON(&configmap)
	if err != nil {
		c.Status(400)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	configmap.ObjectMeta.Name = name
	fmt.Println(configmap)
	err = controllers.Probe.Update(c.Request.Context(), &configmap, &client.UpdateOptions{})
	if errors.IsNotFound(err) {
		c.Status(404)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	} else if err != nil {
		c.Status(500)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	configmap.ObjectMeta.ManagedFields = nil
	c.JSON(200, configmap)
}

func (p *ScenariosAPI) Status(c *gin.Context) {
	name := c.Param("name")
	namespace, ok := c.GetQuery("namespace")
	if !ok {
		namespace = utils.GetCurrentNamespace()
	}
	configmap := coreV1.ConfigMap{}
	err := controllers.Probe.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, &configmap)
	if errors.IsNotFound(err) {
		c.JSON(404, err)
	} else if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	monitor, ok := c.GetQuery("monitoring")
	if ok {
		if monitor == "true" {
			err = createOrUpdateK6Deployment(&configmap)
			if err != nil {
				c.Error(err).SetType(gin.ErrorTypeBind)
				return
			}
		}
	}
	c.JSON(200, configmap)
}

func (p *ScenariosAPI) Patch(c *gin.Context) {
	var err error
	configmap := coreV1.ConfigMap{}
	name := c.Param("name")
	namespace, ok := c.GetQuery("namespace")
	if !ok {
		namespace = utils.GetCurrentNamespace()
	}
	err = controllers.Probe.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, &configmap)
	if errors.IsNotFound(err) {
		c.JSON(404, err)
		return
	} else if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	enabled, ok := c.GetQuery("enabled")
	if ok {
		k6deployment := generateDeploymentObj(&configmap)
		if configmap.Annotations["pronoea.io/enabled"] == enabled {
			c.JSON(200, gin.H{})
			return
		}
		switch enabled {
		case "true":
			err = controllers.Probe.Create(context.TODO(), k6deployment, &client.CreateOptions{})
			if !errors.IsAlreadyExists(err) {
				c.JSON(500, err)
				return
			}
		case "false":
			err = controllers.Probe.Delete(context.TODO(), k6deployment, &client.DeleteOptions{})
			if !errors.IsNotFound(err) {
				c.JSON(500, err)
				return
			}
		default:
			c.JSON(200, gin.H{})
			return
		}
		if configmap.Annotations == nil {
			configmap.Annotations = map[string]string{}
		}
		configmap.Annotations["pronoea.io/enabled"] = enabled
		err := controllers.Probe.Update(context.TODO(), &configmap)
		if err != nil {
			c.JSON(500, configmap)
		}
		return
	}
	data := c.Request.URL.Query()
	c.JSON(200, data)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func generateDeploymentObj(scenarios *coreV1.ConfigMap) *appsV1.Deployment {
	blockOwnerDelete := true
	deployment := appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      scenarios.Name,
			Namespace: scenarios.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				{
					APIVersion:         scenarios.APIVersion,
					Kind:               scenarios.Kind,
					Name:               scenarios.Name,
					UID:                scenarios.UID,
					BlockOwnerDeletion: &blockOwnerDelete,
				},
			},
			Labels: scenarios.Labels,
		},
		Spec: appsV1.DeploymentSpec{
			Selector: &metaV1.LabelSelector{MatchLabels: scenarios.Labels},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: scenarios.Labels,
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Name:  "k6",
							Image: "xingba/k6:output-prometheus-betav0.0.2",
							Ports: []coreV1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: int32(9115),
									Protocol:      coreV1.ProtocolTCP,
								},
							},
							Env: []coreV1.EnvVar{
								{
									Name:  "K6_PROMETHEUS_RW_SERVER_URL",
									Value: "http://prometheus-operated:9090/api/v1/write",
								},
							},
							VolumeMounts: []coreV1.VolumeMount{
								{
									Name:      "k6-scripts",
									MountPath: "/test",
								},
							},
						},
					},
					Volumes: []coreV1.Volume{
						{
							Name: "k6-scripts",
							VolumeSource: coreV1.VolumeSource{
								ConfigMap: &coreV1.ConfigMapVolumeSource{
									LocalObjectReference: coreV1.LocalObjectReference{
										Name: scenarios.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return &deployment
}

func createOrUpdateK6Deployment(scenarios *coreV1.ConfigMap) (err error) {
	blockOwnerDelete := true
	deployment := appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      scenarios.Name,
			Namespace: scenarios.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				{
					APIVersion:         scenarios.APIVersion,
					Kind:               scenarios.Kind,
					Name:               scenarios.Name,
					UID:                scenarios.UID,
					BlockOwnerDeletion: &blockOwnerDelete,
				},
			},
			Labels: scenarios.Labels,
		},
		Spec: appsV1.DeploymentSpec{
			Selector: &metaV1.LabelSelector{MatchLabels: scenarios.Labels},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: scenarios.Labels,
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Name:  "k6",
							Image: "xingba/k6:output-prometheus-betav0.0.2",
							Ports: []coreV1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: int32(9115),
									Protocol:      coreV1.ProtocolTCP,
								},
							},
							VolumeMounts: []coreV1.VolumeMount{
								{
									Name:      "k6-scripts",
									MountPath: "/test",
								},
							},
						},
					},
					Volumes: []coreV1.Volume{
						{
							Name: "k6-scripts",
							VolumeSource: coreV1.VolumeSource{
								ConfigMap: &coreV1.ConfigMapVolumeSource{
									LocalObjectReference: coreV1.LocalObjectReference{
										Name: scenarios.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return controllers.Probe.Create(context.TODO(), &deployment, &client.CreateOptions{})
}
