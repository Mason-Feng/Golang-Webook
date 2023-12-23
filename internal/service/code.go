package service

import "context"

type CodeService struct {
}

func (c *CodeService) Send(ctx context.Context, biz, phone string) error {
	return nil
}

func (c *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return false, nil
}
