package integration

import (
	"gindemo/webbook/internal/integration/startup"
	"gindemo/webbook/internal/service/sms"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
)

type AsyncSMSTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *AsyncSMSTestSuite) SetupSuite() {
	s.db = startup.InitDB()
}

func (s *AsyncSMSTestSuite) TearDownSuite() {
	s.db.Exec("TRUNCATE table `async_sms`")
}

func (s *AsyncSMSTestSuite) TestSend() {
	t := s.T()
	testCases := []struct {
		name string

		// 虽然是集成测试，但是我们也不想真的发短信，所以用 mock
		mock func(ctrl *gomock.Controller) sms.Service

		tplId   string
		args    []string
		numbers []string

		wantErr error
	}{
		{},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

		})
	}
}
