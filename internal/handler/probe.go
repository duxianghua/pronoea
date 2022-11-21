package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	v1 "github.com/duxianghua/pronoea/internal/api/v1"
	"github.com/duxianghua/pronoea/internal/controllers"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
	log.Debug().Msg(fmt.Sprintf("name: %s", name))
	probe := v1.Probe{}
	err := controllers.Probe.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: name}, &probe)
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
	err := c.ShouldBindJSON(&probe)
	if err != nil {
		c.Status(400)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	probe.ObjectMeta.Name = name
	probe.ObjectMeta.Namespace = "default"
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
	probe := v1.Probe{}
	probe.ObjectMeta.Namespace = c.Param("namespace")
	probe.ObjectMeta.Name = c.Param("name")
	err := controllers.Probe.Delete(c.Request.Context(), &probe, &client.DeleteAllOfOptions{})
	if err != nil {
		c.Status(404)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(200, gin.H{})
}

func (p *ProbeAPI) Update(c *gin.Context) {
	probe := v1.Probe{}
	err := c.ShouldBindJSON(&probe)
	if err != nil {
		c.Status(400)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	probe.ObjectMeta.Namespace = c.Param("namespace")
	probe.ObjectMeta.Name = c.Param("name")
	if err := controllers.Probe.Update(c.Request.Context(), &probe, &client.UpdateOptions{}); err != nil {
		c.Status(404)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	probe.ObjectMeta.ManagedFields = nil
	c.JSON(200, probe)
}

func (p *ProbeAPI) Status(c *gin.Context) {
	name := c.Param("name")
	log.Debug().Msg(fmt.Sprintf("name: %s", name))
	probe := v1.Probe{}
	err := controllers.Probe.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: name}, &probe)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	data := gin.H{}
	for _, host := range probe.Spec.Targets {
		url := fmt.Sprintf("http://%s:9115/probe?target=%s&module=%s", controllers.BlackboxServiceName, host, probe.ObjectMeta.Name)
		//url := fmt.Sprintf("http://localhost:63299/probe?target=%s&module=%s", host, probe.ObjectMeta.Name)
		response, err := http.Get(url)
		if err != nil {
			data[host] = err.Error()
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
