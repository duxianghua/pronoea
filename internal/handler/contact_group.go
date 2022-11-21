package handler

import (
	"context"

	v1 "github.com/duxianghua/pronoea/internal/api/v1"
	"github.com/duxianghua/pronoea/internal/controllers"
	"github.com/duxianghua/pronoea/internal/utils"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ContactGroupAPI struct{}

func (p *ContactGroupAPI) List(c *gin.Context) {
	contactGroupList := v1.ContactGroupList{}
	err := controllers.Probe.List(context.TODO(), &contactGroupList)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(200, contactGroupList)
}

func (p *ContactGroupAPI) Get(c *gin.Context) {
	name := c.Param("name")
	namespace, ok := c.GetQuery("namespace")
	if !ok {
		namespace = utils.GetCurrentNamespace()
	}
	contactGroup := v1.ContactGroup{}
	err := controllers.Probe.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, &contactGroup)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	contactGroup.ObjectMeta.ManagedFields = nil
	c.JSON(200, contactGroup)
}

func (p *ContactGroupAPI) Create(c *gin.Context) {
	namespace, ok := c.GetQuery("namespace")
	if !ok {
		namespace = utils.GetCurrentNamespace()
	}
	contactGroup := v1.ContactGroup{}
	name := c.Param("name")
	err := c.ShouldBindJSON(&contactGroup)
	if err != nil {
		c.Status(400)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	contactGroup.ObjectMeta.Name = name
	contactGroup.ObjectMeta.Namespace = namespace
	err = controllers.Probe.Create(context.TODO(), &contactGroup, &client.CreateOptions{})
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	contactGroup.ObjectMeta.ManagedFields = nil
	c.JSON(200, contactGroup)
}

func (p *ContactGroupAPI) Delete(c *gin.Context) {
	name := c.Param("name")
	namespace, ok := c.GetQuery("namespace")
	if !ok {
		namespace = utils.GetCurrentNamespace()
	}
	contactGroup := v1.ContactGroup{}
	contactGroup.ObjectMeta.Namespace = namespace
	contactGroup.ObjectMeta.Name = name
	err := controllers.Probe.Delete(c.Request.Context(), &contactGroup, &client.DeleteAllOfOptions{})
	if err != nil {
		c.Status(404)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	c.JSON(200, gin.H{})
}

func (p *ContactGroupAPI) Update(c *gin.Context) {
	name := c.Param("name")
	namespace, ok := c.GetQuery("namespace")
	if !ok {
		namespace = utils.GetCurrentNamespace()
	}

	contactGroup := v1.ContactGroup{}
	err := c.ShouldBindJSON(&contactGroup)
	if err != nil {
		c.Status(400)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	contactGroup.ObjectMeta.Namespace = namespace
	contactGroup.ObjectMeta.Name = name
	if err := controllers.Probe.Update(c.Request.Context(), &contactGroup, &client.UpdateOptions{}); err != nil {
		c.Status(404)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	contactGroup.ObjectMeta.ManagedFields = nil
	c.JSON(200, contactGroup)
}
