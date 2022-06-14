package auth

import (
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/valyala/fasthttp"
	"golang.org/x/oauth2"
)

type OAuth2 struct {
	Host   string
	Path   string
	Config config.AuthConfig
}

func (oauth *OAuth2) Handler(ctx *fasthttp.RequestCtx) bool {
	conf := oauth.config()

	if string(ctx.Host()) == oauth.Host && string(ctx.Path()) == oauth.Path && ctx.QueryArgs().Has("code") {
		token, err := conf.Exchange(ctx, string(ctx.QueryArgs().Peek("code")))
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return true
		}
	}

	ctx.Redirect(conf.AuthCodeURL("state", oauth2.SetAuthURLParam("app_id", oauth.Config.AppId)), fasthttp.StatusFound)
	return false
}

func (oauth *OAuth2) config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     oauth.Config.ClientId,
		ClientSecret: oauth.Config.AppSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauth.Config.AuthorizeUrl,
			TokenURL: oauth.Config.AccessTokenUrl,
		},
		RedirectURL: oauth.Config.RedirectUrl,
		Scopes:      oauth.Config.Scopes,
	}
}
