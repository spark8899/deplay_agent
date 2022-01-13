package service

import (
    "context"

    "github.com/spark8899/deploy-agent/global"
)

type Service struct {
    ctx context.Context
}

func New(ctx context.Context) Service {
    svc := Service{ctx: ctx}
    return svc
}
