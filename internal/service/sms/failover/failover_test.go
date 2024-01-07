package failover

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/internal/service/sms"
	smsmocks "webook/internal/service/sms/mocks"
)

func TestFailOverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) []sms.SMSService
		wantErr error
	}{
		{
			name: "一次发送成功",
			mock: func(ctrl *gomock.Controller) []sms.SMSService {
				svc0 := smsmocks.NewMockSMSService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.SMSService{svc0}

			},
		},
		{
			name: "第二次发送成功",
			mock: func(ctrl *gomock.Controller) []sms.SMSService {
				svc0 := smsmocks.NewMockSMSService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				svc1 := smsmocks.NewMockSMSService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.SMSService{svc0, svc1}

			},
		},
		{
			name: "全部失败",
			mock: func(ctrl *gomock.Controller) []sms.SMSService {
				svc0 := smsmocks.NewMockSMSService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				svc1 := smsmocks.NewMockSMSService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				return []sms.SMSService{svc0, svc1}

			},
			wantErr: errors.New("轮询了所有的服务商，但是发送都失败了"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewFailOverSMSService(tc.mock(ctrl))
			err := svc.Send(context.Background(), "123", []string{"1234"}, "1234")
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
