package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/internal/domain"
	"webook/internal/service"
	"webook/internal/web/jwt"
	"webook/pkg/logger"
)

type ArticleHandler struct {
	svc    service.ArticleService
	logger logger.LoggerV1
}

func NewArticleHandler(logger logger.LoggerV1, svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		logger: logger,
		svc:    svc,
	}
}
func (ah *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("edit", ah.Edit)
	g.POST("publish", ah.Publish)
}

// 接收Article输入，返回一个文章的ID
func (ah *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	uc := ctx.MustGet("user").(jwt.UsersClaims)
	id, err := ah.svc.Save(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		ah.logger.Error("保存文章数据失败",
			logger.Int64("uid", uc.Uid),
			logger.Error(err))
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})

}

func (ah *ArticleHandler) Publish(ctx *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	uc := ctx.MustGet("user").(jwt.UsersClaims)
	id, err := ah.svc.Publish(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		ah.logger.Error("发表文章数据失败",
			logger.Int64("uid", uc.Uid),
			logger.Error(err))
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})

}
