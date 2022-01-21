package service

import (
    "errors"
    "mime/multipart"
    "os"

    //"github.com/spark8899/deploy-agent/global"
    "github.com/spark8899/deploy-agent/pkg/upload"
)

type FileInfo struct {
    Name      string
    FileMD5   string
}

func (svc *Service) UploadFile(fileType upload.FileType, file multipart.File, fileHeader *multipart.FileHeader, uploadSavePath string) (*FileInfo, error) {
    fileName := upload.GetFileName(fileHeader.Filename)
    if !upload.CheckUploadPath(uploadSavePath) {
        return nil, errors.New("file upload path is not supported.")
    }
    if !upload.CheckContainExt(fileType, fileName) {
        return nil, errors.New("file suffix is not supported.")
    }
    if !upload.CheckFileName(fileName) {
        return nil, errors.New("file name is not supported.")
    }
    if upload.CheckMaxSize(fileType, file) {
        return nil, errors.New("exceeded maximum file limit.")
    }

    //uploadSavePath := upload.GetSavePath()
    if upload.CheckSavePath(uploadSavePath) {
        if err := upload.CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
            return nil, errors.New("failed to create save directory.")
        }
    }
    if upload.CheckPermission(uploadSavePath) {
        return nil, errors.New("insufficient file permissions.")
    }

    dst := uploadSavePath + "/" + fileName
    if err := upload.SaveFile(fileHeader, dst); err != nil {
        return nil, err
    }

    //accessUrl := global.AppSetting.UploadServerUrl + "/" + fileName
    md5Info := upload.GetFileMD5(dst)
    return &FileInfo{Name: fileName, FileMD5: md5Info}, nil
}
