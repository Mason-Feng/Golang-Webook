package service

import (
	"context"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

// var ErrCodeSendTooMany = repository.ErrCodeVerifyTooMany
var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}
type codeService struct {
	repo repository.CodeRepository
	sms  sms.SMSService
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.SMSService) CodeService {
	return &codeService{
		repo: repo,
		sms:  smsSvc,
	}

}

func (svc *codeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	const codeTplId = "1877556"

	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)

}

func (svc *codeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if err == repository.ErrCodeVerifyTooMany {
		//相当于，对外屏蔽了验证次数过多的错误，告诉调用者这个不对
		return false, nil
	}
	return ok, err
}
func (svc *codeService) generate() string {
	//生成范围为0-999999的随机数
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
