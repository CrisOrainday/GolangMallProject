package service

import (
	"context"
	"fmt"

	auth "github.com/cloudwego/biz-demo/gomall/app/frontend/hertz_gen/frontend/auth"
	common "github.com/cloudwego/biz-demo/gomall/app/frontend/hertz_gen/frontend/common"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/infra/rpc"
	rpcuser "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/sessions"
)

type LoginRegisterWithEmailCodeService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLoginRegisterWithEmailCodeService(Context context.Context, RequestContext *app.RequestContext) *LoginRegisterWithEmailCodeService {
	return &LoginRegisterWithEmailCodeService{RequestContext: RequestContext, Context: Context}
}

func (h *LoginRegisterWithEmailCodeService) Run(req *auth.LoginRegisterWithEmailCodeReq) (resp *common.Empty, err error) {
	// 1. 从 session 获取用户提交的邮箱验证码
	session := sessions.Default(h.RequestContext)
	code := session.Get("email_code_" + req.Email)
	// 2. 校验验证码
	if code == nil || code != req.VerifyCode {
		// 如果验证码不一致，返回错误
		return nil, fmt.Errorf("invalid verification code")
	}
	// 此时，验证码通过验证
	hlog.Info("Pass verification")
	res, err := rpc.UserClient.LoginRegister(h.Context, &rpcuser.LoginReq{Email: req.Email, Password: req.Password})
	if err != nil {
		return
	}
	// 保存用户到 sessions
	session.Set("user_id", res.UserId)
	err = session.Save()
	if err != nil {
		return
	}
	return

	// 这里可能要像 /root/project/MallPro/app/frontend/biz/service/login.go 里面一样重定向到跟页面，我这里先不写了
}
