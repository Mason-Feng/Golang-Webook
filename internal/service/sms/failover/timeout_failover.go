package failover

import (
	"context"
	"sync/atomic"
	"webook/internal/service/sms"
)

type TimeoutFailoverSMSService struct {
	svcs []sms.SMSService
	//当前正在使用节点
	idx int32
	//连续几个超时了
	cnt int32
	//切换的阈值，只读的
	threshold int32
}

func NewTimeoutFailoverSMSService(svcs []sms.SMSService, threshold int32) *TimeoutFailoverSMSService {
	return &TimeoutFailoverSMSService{
		svcs:      svcs,
		threshold: threshold,
	}
}
func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	//超过阈值，执行切换
	if cnt >= t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			//重置cnt计数

			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = newIdx
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tplId, args, numbers...)
	switch err {
	case nil:
		//连续超时，所以不超时的时候要重置到0

		atomic.StoreInt32(&t.cnt, 0)
		return nil
	case context.DeadlineExceeded:
		//超时
		atomic.AddInt32(&t.cnt, 1)
	default:
		//遇到了错误，但是又不是超时错误，
		//可以增加，也可以不增加
		//如果抢到一定是超时，可以不增加
		//如果是EOF之类错误，可以直接切换
	}
	return err
}
