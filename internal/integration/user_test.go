package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webook/internal/integration/startup"
	"webook/internal/web"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}
func TestUserHandler_SendSMSCode(t *testing.T) {
	rdb := startup.InitRedis()
	server := startup.InitWebServer()
	testCases := []struct {
		name string
		//BEFORE
		before   func(t *testing.T)
		after    func(t *testing.T)
		phone    string
		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功的用例",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:12345678912"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)
				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9+time.Second+50)
				rdb.Del(ctx, key).Err()

			},
			phone:    "12345678912",
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "未输入手机号",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},

			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 4,
				Msg:  "请输入手机号",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:12345678912"
				err := rdb.Set(ctx, key, "123456", time.Minute*10).Err()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:12345678912"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)

			},
			phone:    "12345678912",
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "短信发送太频繁，请稍候再试",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			//server := startup.InitWebServer()
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewReader([]byte(fmt.Sprintf(`{"phone":"%s"}`, tc.phone))))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			//执行
			server.ServeHTTP(recorder, req)
			//断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			var res web.Result
			json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)
		})
	}
}
