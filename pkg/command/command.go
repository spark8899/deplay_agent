package command

import (
    "bytes"
    "strings"
    "fmt"
    "os/exec"

    "github.com/spark8899/deploy-agent/global"
)

func RunCommand(dst string) (string, error) {
    // get app, args
    parts := strings.Split(dst, " ")
    app, args := parts[0], parts[1:]
    // check app path
    _, err1 := exec.LookPath(app)
    if err1 !=nil {
        return "Command path error", err1
    }
    cmd := exec.Command(app, args...)
    var stdout, stderr bytes.Buffer
    cmd.Stdout, cmd.Stderr = &stdout, &stderr
    err2 := cmd.Run()
    cmd.Process.Wait()
    // get outStr, errStr
    outStr, errStr := stdout.String(), stderr.String()
    info := fmt.Sprintf("%s::%s", outStr, errStr)
    if err2 != nil {
        return info,  err2
    }
    return info, nil
}

func CheckCommand(command string) bool {
    for _, allowCommand := range global.AppSetting.ExecScripts {
        if allowCommand == command {
            return true
        }
    }

    return false
}

func CheckPath(path string) bool {
    for _, allowPath := range global.AppSetting.DeployPath {
        if allowPath == path {
            return true
        }
    }

    return false
}
