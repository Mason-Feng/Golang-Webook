package failover

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"webook/internal/service/sms"
)

type FailOverSMSService struct {
	svcs []sms.SMSService

	//v1字段
	//当前服务商下标
	idx uint64
}

func NewFailOverSMSService(svcs []sms.SMSService) *FailOverSMSService {
	return &FailOverSMSService{
		svcs: svcs,
	}
}
func (f *FailOverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		log.Println(err)

	}
	return errors.New("轮询了所有的服务商，但是发送都失败了")
}

// 该方法特点是起始svc是动态计算的
// 起始下标轮询
// 并且出错也轮询
func (f *FailOverSMSService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	//迭代length次
	for i := idx; i < idx+length; i++ {
		//取余数计算下标
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplId, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			//前者是被取消，后者是超时
			return err

		}
		log.Println(err)

	}
	return errors.New("轮询了所有的服务商，但是发送都失败了")
}
