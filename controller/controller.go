package controller

import (
	"net/http"
	"prometheus-webhook-wechat/models"
	"prometheus-webhook-wechat/notifier"
	"prometheus-webhook-wechat/template"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
)

func CallWechatController(c *models.NewGinContext) {
	var req models.Data
	callID := uuid.New().String()
	err := c.BindJSON(&req)
	logger := c.Logger
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Bad Request": err.Error(),
		})
		return
	}
	alertContent, err := template.TransferContent(req, c.TemplateFile)
	if err != nil {
		level.Error(c.Logger).Log("TraceID", callID, "msg", "Generate alert content failed", "error", err)
		return
	}
	alertContentLen := utf8.RuneCountInString(alertContent)
	if alertContentLen > c.MaxContentLength {
		level.Error(c.Logger).Log("TraceID", callID, "msg", "Too large content")
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
	level.Info(logger).Log("TraceID", callID, "ClientIP", c.ClientIP(), "UserAgent", c.Request.UserAgent())
	notifier.SendNotification(c.NotifyTargets, alertContent, logger, callID)
}
