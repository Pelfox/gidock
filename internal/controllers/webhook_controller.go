package controllers

import "github.com/gin-gonic/gin"

type WebhookController struct{}

func NewWebhookController() *WebhookController {
	return &WebhookController{}
}

func (c *WebhookController) HandleWebhook(ctx *gin.Context) {
}
