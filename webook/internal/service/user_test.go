package service

import (
	"context"
	"errors"
	"gindemo/webook/internal/domain"
	"gindemo/webook/internal/repository"
	repomocks "gindemo/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("123456#hello")
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, []byte("123456#hello"))
	assert.NoError(t, err)
}

func TestUserService_Login(t *testing.T) {
	testcases := []struct {
		name string

		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 预期收入
		ctx      context.Context
		email    string
		password string

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$03WK4TVHs5f5Nl5ymE1m5uSrYYIIkCipeplmGQ7iCGG6CgkDZH2LW",
						Phone:    "15212345678",
					}, nil)
				return repo
			},
			email: "123@qq.com",
			// 用户输入的。没有加密的
			password: "123456#hello",

			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$03WK4TVHs5f5Nl5ymE1m5uSrYYIIkCipeplmGQ7iCGG6CgkDZH2LW",
				Phone:    "15212345678",
			},
		},
		{
			name: "用户未找到",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123@qq.com",
			password: "123456#hello",
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("db错误"))
				return repo
			},
			email:    "123@qq.com",
			password: "123456#hello",
			wantErr:  errors.New("db错误"),
		},
		{
			name: "密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "123456#hello$2a$10$03WK4TVHs5f5Nl5ymE1m5uSrYYIIkCipeplmGQ7iCGG6CgkDZH2LW",
						Phone:    "15212345678",
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "123456#he",
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mock(ctrl)
			svc := NewUserService(repo)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantUser, user)
			assert.Equal(t, tc.wantErr, err)
		})

	}
}
