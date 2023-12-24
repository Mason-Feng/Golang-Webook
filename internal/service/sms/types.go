package sms

import "context"

//go:generate mockgen -source./types.go -package=smsmocks -destination=./mocks/sms.mock.go Service
type SMSService interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}

type Req struct {
}
