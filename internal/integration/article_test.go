package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/internal/integration/startup"
	"webook/internal/repository/dao"
	ijwt "webook/internal/web/jwt"
)

type ArticleHandlerSuite struct {
	suite.Suite
	db     *gorm.DB
	server *gin.Engine
}

func TestArticleHandler(t *testing.T) {
	suite.Run(t, &ArticleHandlerSuite{})
}
func (as *ArticleHandlerSuite) SetupSuite() {
	as.db = startup.InitDB()
	hdl := startup.InitArticleHandler(dao.NewArticleGORMDAO(as.db))
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user", ijwt.UsersClaims{
			Uid: 123,
		})
	})
	hdl.RegisterRoutes(server)
	as.server = server
}

func (as *ArticleHandlerSuite) TestEdit() {
	t := as.T()

	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		//前端传过来，肯定是一个json
		art Article

		wantCode int
		wantRes  Result[int64]
	}{
		{
			name:   "新建帖子",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				//要验证保存到了数据库中
				var art dao.Article
				err := as.db.Where("author_id=?", 123).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				assert.True(t, art.Id > 0)
				assert.Equal(t, "我的标题", art.Title)
				assert.Equal(t, "我的内容", art.Content)
				assert.Equal(t, int64(123), art.AuthorId)

			},
			art: Article{

				Title:   "我的标题",
				Content: "我的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
			},
		},
		{
			name: "修改帖子",
			before: func(t *testing.T) {
				err := as.db.Create(dao.Article{
					Id:       2,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//要验证保存到了数据库中
				var art dao.Article
				err := as.db.Where("id=?", 2).First(&art).Error
				assert.NoError(t, err)

				assert.True(t, art.Utime > 789)

				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       2,
					Title:    "新的标题",
					Content:  "新的内容",
					AuthorId: 123,
					Ctime:    456,
				}, art)

				//as.db.Exec("truncate table `article`")
			},
			art: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 2,
			},
		},
		{
			name: "修改帖子-修改别人的帖子",
			before: func(t *testing.T) {
				err := as.db.Create(dao.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 234,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//要验证保存到了数据库中
				var art dao.Article
				err := as.db.Where("id=?", 3).First(&art).Error
				assert.NoError(t, err)

				assert.Equal(t, dao.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 234,
					Ctime:    456,
					Utime:    789,
				}, art)

				//as.db.Exec("truncate table `article`")
			},
			art: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg: "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)

			//server := startup.InitWebServer()
			req, err := http.NewRequest(http.MethodPost,
				"/articles/edit",
				bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			//执行
			as.server.ServeHTTP(recorder, req)
			//断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			var res Result[int64]
			err = json.NewDecoder(recorder.Body).Decode(&res)
			if err != nil {
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
func (as *ArticleHandlerSuite) TearDownTest() {
	err := as.db.Exec("truncate table `articles`").Error
	assert.NoError(as.T(), err)

}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
type Article struct {
	Id      int64
	Title   string `json:"title"`
	Content string `json:"content"`
}
