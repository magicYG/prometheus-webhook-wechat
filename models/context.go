package models

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
)

type NewGinContext struct {
	*gin.Context
	NotifyTargets    []Target `json:"notify_targets"`
	Logger           log.Logger
	TemplateFile     *string
	MaxContentLength int
}

type HandlerFunc func(ctx *NewGinContext)
