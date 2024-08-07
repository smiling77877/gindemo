package web

import (
	"bytes"
	"encoding/json"
	"errors"
	service2 "gindemo/webook/interactive/service"
	"gindemo/webook/internal/domain"
	"gindemo/webook/internal/service"
	svcmocks "gindemo/webook/internal/service/mocks"
	ijwt "gindemo/webook/internal/web/jwt"
	"gindemo/webook/pkg/ginx"
	"gindemo/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (service.ArticleService, service2.InteractiveService)

		reqBody  string
		wantCode int
		wantRes  ginx.Result
	}{
		{
			name: "新建并且发表成功",
			mock: func(ctrl *gomock.Controller) (service.ArticleService, service2.InteractiveService) {
				svc, intr := svcmocks.NewMockArticleService(ctrl), svcmocks.NewMockInteractiveService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc, intr
			},
			reqBody: `
{
"title": "我的标题",
"content": "我的内容"
}
`,
			wantCode: http.StatusOK,
			wantRes: ginx.Result{
				// 原本是 int64的，但是因为 Data 是any，所以在反序列化的时候用的 float64
				Data: float64(1),
			},
		},
		{
			name: "修改并且发表成功",
			mock: func(ctrl *gomock.Controller) (service.ArticleService, service2.InteractiveService) {
				svc, intr := svcmocks.NewMockArticleService(ctrl), svcmocks.NewMockInteractiveService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "新的标题",
					Content: "新的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc, intr
			},
			reqBody: `
{
"id": 1,
"title": "新的标题",
"content": "新的内容"
}
`,
			wantCode: http.StatusOK,
			wantRes: ginx.Result{
				Data: float64(1),
			},
		},
		{
			name: "输入有误",
			mock: func(ctrl *gomock.Controller) (service.ArticleService, service2.InteractiveService) {
				svc, intr := svcmocks.NewMockArticleService(ctrl), svcmocks.NewMockInteractiveService(ctrl)
				return svc, intr
			},
			reqBody: `
{
"id": 1,
"title": "新的标题",
"content": "新的内容",,,,
}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "publish错误",
			mock: func(ctrl *gomock.Controller) (service.ArticleService, service2.InteractiveService) {
				svc, intr := svcmocks.NewMockArticleService(ctrl), svcmocks.NewMockInteractiveService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "新的标题",
					Content: "新的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("mock error"))
				return svc, intr
			},
			reqBody: `
{
"id": 1,
"title": "新的标题",
"content": "新的内容"
}
`,
			wantCode: http.StatusOK,
			wantRes: ginx.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 构造Handler
			svc, interact := tc.mock(ctrl)
			hdl := NewArticleHandler(logger.NewNopLogger(), svc, interact)

			// 准备服务器，注册路由
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user", ijwt.UserClaims{
					Uid: 123,
				})
			})
			hdl.RegisterRoutes(server)

			// 准备Req和记录的 recorder
			req, err := http.NewRequest(http.MethodPost, "/articles/publish",
				bytes.NewReader([]byte(tc.reqBody)))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			// 执行
			server.ServeHTTP(recorder, req)
			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			if recorder.Code != http.StatusOK {
				return
			}
			var res ginx.Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
