package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"webook/internal/domain"
	"webook/pkg/logger"
)

type OAuth2Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

var redirectURL = url.PathEscape(`https://passport.yhd.com/wechat/callback.do`)

type OAuth2WechatService struct {
	appID     string
	appSecret string
	client    *http.Client
	logger    logger.LoggerV1
}

type Result struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`

	//错误返回

	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func NewOAuth2WechatService(appID string, appSecret string, logger logger.LoggerV1) OAuth2Service {
	return &OAuth2WechatService{
		appID:     appID,
		appSecret: appSecret,
		client:    http.DefaultClient,
		logger:    logger,
	}
}

func (oauths *OAuth2WechatService) AuthURL(ctx context.Context, state string) (string, error) {
	const authURLPattern = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`

	return fmt.Sprintf(authURLPattern, oauths.appID, redirectURL, state), nil

}
func (oauths *OAuth2WechatService) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	accessTokenUrl := fmt.Sprintf(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`, oauths.appID, oauths.appSecret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	httpResp, err := oauths.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	var res Result
	err = json.NewDecoder(httpResp.Body).Decode(&res)
	if err != nil {

		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("调用微信接口失败 errcode:%s,errmsg:%s", res.ErrCode, res.ErrMsg)
	}
	return domain.WechatInfo{
		UnionId: res.Unionid,
		OpenId:  res.OpenId,
	}, nil
}
