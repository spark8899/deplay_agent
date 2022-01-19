package command

import (
    "bytes"
    "strings"
    "fmt"
    "os/exec"

    "github.com/spark8899/deploy-agent/global"
)

func RunCommand(command string) (string, error) {
    // get app, args
    parts := strings.Split(command, " ")
    app, args := parts[0], parts[1:]
    // check app path
    _, err1 := exec.LookPath(app)
    if err1 !=nil {
        return "path error\n", err1
    }
    cmd := exec.Command(app, args...)
    var stdout, stderr bytes.Buffer
    cmd.Stdout, cmd.Stderr = &stdout, &stderr
    err2 := cmd.Run()
    cmd.Process.Wait()
    // get outStr, errStr
    outStr, errStr := stdout.String(), stderr.String()
    if err2 != nil {
        return fmt.Sprintf("out:%s,err:%s", outStr, errStr),  err2
    }
    info := fmt.Sprintf("out:\n%s\nerr:\n%s\n", outStr, errStr)
    return info, nil
}

func CheckCommand(command  string) bool {
    for _, allowCommand := range global.AppSetting.ExecScripts {
        if allowCommand == command {
            return true
        }
    }

    return false
}
