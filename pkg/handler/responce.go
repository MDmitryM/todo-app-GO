package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorRsponce struct {
	Message string `json:"message"`
}

type statusResponce struct {
	Status string `json:"status"`
}

func newErrorResponce(c *gin.Context, statusCode int, message string) {
	logrus.Errorf(message)
	c.AbortWithStatusJSON(statusCode, errorRsponce{message})
}
