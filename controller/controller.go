package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"net/http"
	"prometheus-webhook-wechat/models"
	"prometheus-webhook-wechat/notifier"
	"prometheus-webhook-wechat/template"
)

func CallWechatController(c *models.NewGinContext) {
	var req models.Data
	callID := uuid.New().String()
	err := c.Bind(&req)
	logger := c.Logger
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Bad Request": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
	alertContent, err := template.TransferContent(req, c.TemplateFile)
	if err != nil {
		level.Error(c.Logger).Log("msg", "Generate alert content failed", "error", err)
		return
	}
	level.Info(logger).Log("TraceID", callID, "ClientIP", c.ClientIP(), "UserAgent", c.Request.UserAgent())
	notifier.SendNotification(c.NotifyTargets, alertContent, logger, callID)
}
