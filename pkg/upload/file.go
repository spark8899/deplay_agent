package upload

import (
    "io"
    "io/ioutil"
    "mime/multipart"
    "os"
    "path"
    "strings"

    "github.com/spark8899/deploy-agent/global"
    "github.com/spark8899/deploy-agent/pkg/util"
)

type FileType int

const TypeFileGroup FileType = iota + 1

func GetFileName(name string) string {
    ext := GetFileExt(name)
    fileName := strings.ToLower(strings.TrimSuffix(name, ext))
    //fileName = util.EncodeMD5(fileName)

    return fileName + ext
}

func GetFileExt(name string) string {
    return strings.ToLower(path.Ext(name))
}

//func GetSavePath() string {
    //return global.AppSetting.UploadSavePath
//    return global.AppSetting.DeployPath
//}

//func GetServerUrl() string {
//    return global.AppSetting.UploadServerUrl
//}

func CheckSavePath(dst string) bool {
    _, err := os.Stat(dst)

    return os.IsNotExist(err)
}

func CheckUploadPath(dst string) bool {
    checkPath := strings.ToUpper(dst)
    for _, allowSavePath := range global.AppSetting.DeployPath {
        if strings.ToUpper(allowSavePath) == checkPath {
            return true
        }
    }

    return false
}

func CheckContainExt(t FileType, name string) bool {
    ext := GetFileExt(name)
    ext = strings.ToUpper(ext)
    switch t {
    case TypeFileGroup:
        for _, allowExt := range global.AppSetting.UploadAllowExts {
            if strings.ToUpper(allowExt) == ext {
                return true
            }
        }

    }

    return false
}

func CheckFileName(FileName string) bool {
    checkName := strings.ToUpper(FileName)
    for _, allowFileName := range global.AppSetting.DeployFiles {
        if strings.ToUpper(allowFileName) == checkName {
            return true
        }
    }

    return false
}

func CheckMaxSize(t FileType, f multipart.File) bool {
    content, _ := ioutil.ReadAll(f)
    size := len(content)
    switch t {
    case TypeFileGroup:
        if size >= global.AppSetting.UploadMaxSize*1024*1024 {
            return true
        }
    }

    return false
}

func GetFileMD5(filePath string) string {
    return util.FileMD5(filePath)
}

func CheckPermission(dst string) bool {
    _, err := os.Stat(dst)

    return os.IsPermission(err)
}

func CreateSavePath(dst string, perm os.FileMode) error {
    err := os.MkdirAll(dst, perm)
    if err != nil {
        return err
    }

    return nil
}

func SaveFile(file *multipart.FileHeader, dst string) error {
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, src)
    return err
}
