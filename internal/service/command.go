package service

import (
    "fmt"
    "errors"

    "github.com/spark8899/deploy-agent/pkg/command"
)

type ExecCommandRequest struct {
    Command       string `form:"command" binding:"required,max=100"`
}

func (svc *Service) ExecCommand(param *ExecCommandRequest) (string, error) {
    if !command.CheckCommand(param.Command) {
        return "Command is not allow", errors.New(fmt.Sprintf("command: `%s` is not allow.", param.Command))
    }
    return command.RunCommand(param.Command)
}
