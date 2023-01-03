package handler

import (
	"context"
	"strings"

	v1 "github.com/duxianghua/pronoea/internal/api/v1"
	"github.com/duxianghua/pronoea/internal/config"
	"github.com/duxianghua/pronoea/internal/controllers"
	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/alertmanager/template"
	"github.com/rs/zerolog/log"
)

type AlertmanagerWebhook struct{}

func (aw *AlertmanagerWebhook) Post(c *gin.Context) {
	result := getResult("200", "", nil)
	var message webhook.Message
	err := c.BindJSON(&message)

	if err != nil {
		result.Code = "400"
		result.Msg = err.Error()
		return
	}

	// Check whether project exist in labels
	if ok := labelExist("project", message); !ok {
		result.Code = "500"
		result.Msg = "no project found in labels"
		log.Error().Interface("labels", message.CommonLabels).Msg(result.Msg)
		c.JSON(200, result)
		return
	}

	// 获取联系人信息
	contactEmails := getContact(message.CommonLabels["project"])
	if len(contactEmails) == 0 {
		result.Code = "500"
		result.Msg = "no matching contacts"
		log.Error().Interface("labels", message.CommonLabels).Msg(result.Msg)
		c.JSON(200, result)
		return
	}
	log.Debug().Interface("labels", message.CommonLabels).Strs("mails", contactEmails).Msg("get contact email list success")

	// 渲染模版
	tmpl, err := template.FromGlobs("./templates/*.tmpl")
	if err != nil {
		result.Code = "500"
		result.Msg = err.Error()
		log.Error().Interface("labels", message.CommonLabels).Msg(result.Msg)
		c.JSON(200, result)
		return
	}

	body, err := tmpl.ExecuteHTMLString(
		config.Options.Email.Html, message.Data)
	if err != nil {
		result.Code = "500"
		result.Msg = err.Error()
		log.Error().Interface("labels", message.CommonLabels).Msg(result.Msg)
		c.JSON(200, result)
		return
	}
	subject, err := tmpl.ExecuteHTMLString(
		config.Options.Email.Subject, message.Data)

	if err != nil {
		result.Code = "500"
		result.Msg = err.Error()
		log.Error().Interface("labels", message.CommonLabels).Msg(result.Msg)
		c.JSON(200, result)
		return
	}

	// send mail
	SendEmail(subject, body, contactEmails...)
	c.JSON(200, result)
}

func SendEmail(subject string, body string, to ...string) {
	m := gomail.NewMessage()
	from := m.FormatAddress("heimdallr@homepartners.tech", "DevOps Alerts")
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	if len(strings.Split(config.Options.Email.Bcc, ",")) > 0 {
		m.SetHeader("Bcc", strings.Split(config.Options.Email.Bcc, ",")...)
	}
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(config.Options.Email.Host, config.Options.Email.Port, config.Options.Email.Username, config.Options.Email.Password)
	if err := d.DialAndSend(m); err != nil {
		log.Error().Str("From", from).Strs("To", to).Msg(err.Error())
		return
	}
	log.Info().Str("subject", subject).Str("From", from).Strs("To", to).Msg("send email success")
}

func labelExist(label string, m webhook.Message) bool {
	for _, v := range m.CommonLabels.Names() {
		if v == label {
			return true
		}
	}
	return false
}

func getContact(project string) (emails []string) {
	contactGroupList := v1.ContactGroupList{}
	err := controllers.Probe.List(context.TODO(), &contactGroupList)
	if err != nil {
		log.Error().Err(err).Msg("get contact fail")
		return emails
	}
	for _, cg := range contactGroupList.Items {
		if cg.Spec.Project == project {
			return cg.Spec.Members
		}
	}
	return emails
}
