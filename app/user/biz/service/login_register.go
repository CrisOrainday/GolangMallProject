package service

import (
	"context"

	"github.com/cloudwego/biz-demo/gomall/app/user/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/user/biz/model"
	user "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/crypto/bcrypt"
)

type LoginRegisterService struct {
	ctx context.Context
} // NewLoginRegisterService new LoginRegisterService
func NewLoginRegisterService(ctx context.Context) *LoginRegisterService {
	return &LoginRegisterService{ctx: ctx}
}

// Run create note info
func (s *LoginRegisterService) Run(req *user.LoginReq) (resp *user.LoginResp, err error) {
	// Finish your business logic.
	klog.Infof("Req email: %s", req.Email)
	userRow, err := model.GetByEmail(mysql.DB, s.ctx, req.Email)
	// 4. 如果用户不存在，进行注册
	if err != nil {
		klog.Info("User is not existed")
		hashedPassword, err1 := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err1 != nil {
			return
		}
		newUser := &model.User{
			Email:          req.Email,
			PasswordHashed: string(hashedPassword),
		}
		if err1 = model.Create(mysql.DB, s.ctx, newUser); err1 != nil {
			return
		}

		return &user.LoginResp{UserId: int32(newUser.ID)}, nil
	}
	klog.Infof("Found email: %s", userRow.Email)
	// 5. 如果用户已存在，进行登录
	klog.Info("User is existed")
	return &user.LoginResp{UserId: int32(userRow.ID)}, nil
}
