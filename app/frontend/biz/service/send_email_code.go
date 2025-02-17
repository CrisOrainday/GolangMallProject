package service

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	auth "github.com/cloudwego/biz-demo/gomall/app/frontend/hertz_gen/frontend/auth"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/sessions"
)

type SendEmailCodeService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSendEmailCodeService(Context context.Context, RequestContext *app.RequestContext) *SendEmailCodeService {
	return &SendEmailCodeService{RequestContext: RequestContext, Context: Context}
}

func (h *SendEmailCodeService) Run(req *auth.SendEmailCodeReq) (resp *auth.SendEmailCodeResp, err error) {
	// 1. 验证邮箱格式
	if !isValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email address format")
	}

	// 2. 生成验证码（6位数）
	code := generateCode()

	// 3. 获取 session 并保存验证码到 session
	session := sessions.Default(h.RequestContext)
	session.Set("email_code_"+req.Email, code) // 存储邮箱与验证码的映射
	err = session.Save()
	if err != nil {
		hlog.Errorf("failed to save session: %v", err)
		return nil, err
	}

	// 4. 模拟发送邮件验证码（实际项目中可以调用邮件平台API）
	// 此处仅打印出验证码来模拟发送邮件
	hlog.Infof("Sending Email code %s to email %s", code, req.Email)

	// 5. 返回成功响应
	resp = &auth.SendEmailCodeResp{
		Success: true,
		Message: fmt.Sprintf("Email code sent successfully to %s", req.Email),
	}
	return resp, nil
}

// isValidEmail 验证邮箱格式
func isValidEmail(email string) bool {
	// 使用正则表达式验证邮箱格式
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// generateCode 生成一个6位的随机验证码
func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000 // 生成100000到999999之间的随机数
	return strconv.Itoa(code)
}