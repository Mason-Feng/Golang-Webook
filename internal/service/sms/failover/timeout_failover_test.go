package failover

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/internal/service/sms"
	smsmocks "webook/internal/service/sms/mocks"
)

func TestTimeoutFailoverSMSService_Send(t1 *testing.T) {
	testCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) []sms.SMSService
		threshold int32
		idx       int32
		cnt       int32
		wantErr   error
		wantCnt   int32
		wantIdx   int32
	}{
		{
			name: "没有触发切换",
			mock: func(ctrl *gomock.Controller) []sms.SMSService {
				svc0 := smsmocks.NewMockSMSService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.SMSService{svc0}

			},
			idx:       0,
			cnt:       12,
			threshold: 15,
			wantIdx:   0,
			wantCnt:   0,
			wantErr:   nil,
		},
		{
			name: "触发切换,成功",
			mock: func(ctrl *gomock.Controller) []sms.SMSService {
				svc0 := smsmocks.NewMockSMSService(ctrl)

				svc1 := smsmocks.NewMockSMSService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.SMSService{svc0, svc1}

			},
			idx:       0,
			cnt:       15,
			threshold: 15,
			wantIdx:   1,
			wantCnt:   0,
			wantErr:   nil,
		},
		{
			name: "触发切换,失败",
			mock: func(ctrl *gomock.Controller) []sms.SMSService {
				svc0 := smsmocks.NewMockSMSService(ctrl)

				svc1 := smsmocks.NewMockSMSService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.SMSService{svc0, svc1}

			},
			idx:       0,
			cnt:       15,
			threshold: 15,

			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t1.Run(tc.name, func(t1 *testing.T) {
			ctrl := gomock.NewController(t1)
			defer ctrl.Finish()
			svc := NewTimeoutFailoverSMSService(tc.mock(ctrl), tc.threshold)
			svc.cnt = tc.cnt
			svc.idx = tc.idx
			err := svc.Send(context.Background(), "1234", []string{"12", "34"}, "12345678912")
			assert.Equal(t1, tc.wantErr, err)
			assert.Equal(t1, tc.wantCnt, svc.cnt)
			assert.Equal(t1, tc.wantIdx, svc.idx)
		})
	}
}
