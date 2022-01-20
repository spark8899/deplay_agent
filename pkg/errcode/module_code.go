package errcode

var (
    ErrorUploadFileFail = NewError(20030001, "上传文件失败")
    ErrorCommandFail = NewError(2003002, "命令执行失败")
    ErrorCommandNotAllow = NewError(2003003, "命令不允许执行")
    ErrorCommandPath = NewError(2003004, "命令路径找不到")
)
