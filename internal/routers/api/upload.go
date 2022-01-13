package api

import (
    "github.com/gin-gonic/gin"
    "github.com/spark8899/deploy-agent/global"
    "github.com/spark8899/deploy-agent/internal/service"
    "github.com/spark8899/deploy-agent/pkg/app"
    "github.com/spark8899/deploy-agent/pkg/convert"
    "github.com/spark8899/deploy-agent/pkg/errcode"
    "github.com/spark8899/deploy-agent/pkg/upload"
)

type Upload struct{}

func NewUpload() Upload {
    return Upload{}
}

func (u Upload) UploadFile(c *gin.Context) {
    response := app.NewResponse(c)
    file, fileHeader, err := c.Request.FormFile("file")
    if err != nil {
        response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
        return
    }

    fileType := convert.StrTo(c.PostForm("type")).MustInt()
    if fileHeader == nil || fileType <= 0 {
        response.ToErrorResponse(errcode.InvalidParams)
        return
    }

    svc := service.New(c.Request.Context())
    fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
    if err != nil {
        global.Logger.Errorf(c, "svc.UploadFile err: %v", err)
        response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
        return
    }

    response.ToResponse(gin.H{
        "file_access_url": fileInfo.AccessUrl,
    })
}
