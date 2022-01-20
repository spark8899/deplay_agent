package api

import (
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/spark8899/deploy-agent/global"
    "github.com/spark8899/deploy-agent/internal/service"
    "github.com/spark8899/deploy-agent/pkg/app"
    "github.com/spark8899/deploy-agent/pkg/errcode"
)

func PostExecCommand(c *gin.Context) {
    param := service.ExecCommandRequest{}
    response := app.NewResponse(c)
    valid, errs := app.BindAndValid(c, &param)
    if !valid {
        global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
        response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
        return
    }

    svc := service.New(c.Request.Context())
    info, err := svc.ExecCommand(&param)
    if err != nil {
        global.Logger.Errorf(c, "svc.exec command: `%v` info: `%v` err: `%v`", param.Command, info, err)
        response.ToResponse(gin.H{"code": 2003002, "msg": info})
        //response.ToErrorResponse(errcode.ErrorCommandFail)
        return
    }
    global.Logger.Infof(c, "svc.exec command: `%v`, info: `%v`", param.Command, info)

    strArr := strings.Split(info, `::`)
    outStr, errStr := strArr[0], strArr[1]
    response.ToResponse(gin.H{"out": outStr, "err": errStr})
}
