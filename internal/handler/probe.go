package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	v1 "github.com/duxianghua/pronoea/internal/api/v1"
	"github.com/duxianghua/pronoea/internal/controllers"
	"github.com/duxianghua/pronoea/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ProbeAPI struct{}

func (p *ProbeAPI) List(c *gin.Context) {
	probeList := v1.ProbeList{}
	err := controllers.Probe.List(context.TODO(), &probeList, &client.ListOptions{})
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(200, probeList)
}

func (p *ProbeAPI) Get(c *gin.Context) {
	name := c.Param("name")
	namespace, ok := c.GetQuery("namespace")
	if !ok {
		namespace = utils.GetCurrentNamespace()
	}
	probe := v1.Probe{}
	err := controllers.Probe.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, &probe)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	probe.ObjectMeta.ManagedFields = nil
	c.JSON(200, probe)
}

func (p *ProbeAPI) Create(c *gin.Context) {
	probe := v1.Probe{}
	name := c.Param("name")
	namespace, ok := c.GetQuery("namespace")
	if !ok {
		namespace = utils.GetCurrentNamespace()
	}
	err := c.ShouldBindJSON(&probe)
	if err != nil {
		c.Status(400)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	probe.ObjectMeta.Name = name
	probe.ObjectMeta.Namespace = namespace
	if !probe.Spec.Pause {
		probe.Spec.Pause = false
	}
	err = controllers.Probe.Create(context.TODO(), &probe, &client.CreateOptions{})
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	probe.ObjectMeta.ManagedFields = nil
	c.JSON(200, probe)
}

func (p *ProbeAPI) Delete(c *gin.Context) {
	name := c.Param("name")
	namespace, ok := c.GetQuery("namespace")
	if !ok {
		namespace = utils.GetCurrentNamespace()
	}
	probe := v1.Probe{}
	probe.ObjectMeta.Namespace = namespace
	probe.ObjectMeta.Name = name
	err := controllers.Probe.Delete(c.Request.Context(), &probe, &client.DeleteAllOfOptions{})
	if err != nil {
		c.Status(404)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(200, gin.H{})
}

func (p *ProbeAPI) Update(c *gin.Context) {
	name := c.Param("name")
	// namespace, ok := c.GetQuery("namespace")
	// if !ok {
	// 	namespace = utils.GetCurrentNamespace()
	// }
	probe := v1.Probe{}
	err := c.ShouldBindJSON(&probe)
	if err != nil {
		c.Status(400)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	// probe.ObjectMeta.Namespace = namespace
	probe.ObjectMeta.Name = name
	fmt.Printf("versionID1: %s \n", probe.ObjectMeta.ResourceVersion)
	err = controllers.Probe.Update(c.Request.Context(), &probe, &client.UpdateOptions{})
	fmt.Printf("versionID2: %s \n", probe.ObjectMeta.ResourceVersion)
	if errors.IsNotFound(err) {
		c.Status(404)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	} else if err != nil {
		c.Status(500)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	// controllers.Probe.Get(c.Request.Context(), types.NamespacedName{
	// 	Namespace: namespace,
	// 	Name:      name}, &probe)
	fmt.Printf("versionID3: %s \n", probe.ObjectMeta.ResourceVersion)
	probe.ObjectMeta.ManagedFields = nil
	c.JSON(200, probe)
}

func (p *ProbeAPI) Status(c *gin.Context) {
	name := c.Param("name")
	namespace, ok := c.GetQuery("namespace")
	if !ok {
		namespace = utils.GetCurrentNamespace()
	}
	probe := v1.Probe{}
	err := controllers.Probe.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, &probe)
	if errors.IsNotFound(err) {
		c.JSON(404, err)
	} else if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	data := gin.H{}
	for _, host := range probe.Spec.Targets {
		url := fmt.Sprintf("http://%s.%s.svc.cluster.local:9115/probe?target=%s&module=%s&debug=true", controllers.BlackboxServiceName, namespace, host, probe.ObjectMeta.Name)
		//url := fmt.Sprintf("http://localhost:59906/probe?target=%s&module=%s", host, probe.ObjectMeta.Name)
		response, err := http.Get(url)
		if err != nil {
			log.Error().Err(err).Msg("access blackbox error")
			data[host] = err.Error()
			continue
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			data[host] = err.Error()
		} else {
			data[host] = string(body)
		}
	}
	c.JSON(200, data)
}
