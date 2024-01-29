package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/internal/domain"
	"webook/internal/service"
	svcmocks "webook/internal/service/mocks"
	ijwt "webook/internal/web/jwt"
	"webook/pkg/logger"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.ArticleService
		reqBody  string
		wantCode int
		wantRes  Result
	}{
		{
			name: "新建并发表成功",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `
{
	"title":"我的标题",
	"content":"我的内容"
}`,
			wantCode: 200,
			wantRes: Result{
				Data: float64(1),
			},
		},
		{
			name: "已有帖子发表成功",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      123,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(123), nil)
				return svc
			},
			reqBody: `
{
	"id":123,
	"title":"我的标题",
	"content":"我的内容"
}`,
			wantCode: 200,
			wantRes: Result{
				Data: float64(123),
			},
		},
		{
			name: "发表失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      123,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("mock error"))
				return svc
			},
			reqBody: `
{
	"id":123,
	"title":"我的标题",
	"content":"我的内容"
}`,
			wantCode: 200,
			wantRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "Bind错误",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				return svc
			},
			reqBody: `
{
	"id":123,
	"title":"我的标题",
	"content":"我的内容"dsd
}`,
			wantCode: 400,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := tc.mock(ctrl)
			hdl := NewArticleHandler(logger.NewNopLogger(), svc)
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user", ijwt.UsersClaims{
					Uid: 123,
				})
			})
			hdl.RegisterRoutes(server)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/publish",
				bytes.NewBufferString(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)
			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}
			var res Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
