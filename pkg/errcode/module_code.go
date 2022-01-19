package errcode

var (
    ErrorUploadFileFail = NewError(20030001, "上传文件失败")
    ErrorCommandFail = NewError(2003002, "命令执行失败")
)
