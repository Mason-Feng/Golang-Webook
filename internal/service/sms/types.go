package sms

import "context"

// Service发送短信的抽象
// 这是一个为了适配不同的短信供应商的抽象
type Service interface {
	Send(ctx context.Context, number string, tplId string, args []string, numbers ...string) error
}
