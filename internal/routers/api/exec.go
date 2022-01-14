package api

import (
    "github.com/gin-gonic/gin"
    "github.com/spark8899/deploy-agent/global"
    "github.com/spark8899/deploy-agent/pkg/app"
    "github.com/spark8899/deploy-agent/pkg/convert"
    "github.com/spark8899/deploy-agent/pkg/errcode"
)

func Exec(c *gin.Context) {
    c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
    response.ToResponse(gin.H{})
}
