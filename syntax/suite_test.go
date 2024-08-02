package syntax

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type DemoSuite struct {
	suite.Suite
}

func (s *DemoSuite) SetupSubTest() {
	s.T().Log("这是SetupSubTest")
}

func (s *DemoSuite) TearDownSubTest() {
	s.T().Log("这是TearDownSubTest")
}

func (s *DemoSuite) BeforeTest(suiteName, testName string) {
	s.T().Log("这是BeforeTest")
}

func (s *DemoSuite) AfterTest(suiteName, testName string) {
	s.T().Log("这是AfterTest")
}

func (s *DemoSuite) TearDownTest() {
	s.T().Log("这是TearDownTest")
}

func (s *DemoSuite) SetupTest() {
	s.T().Log("这是SetupTest")
}

func (s *DemoSuite) SetupSuite() {
	s.T().Log("这是SetupSuite")
}

func (s *DemoSuite) TearDownSuite() {
	s.T().Log("这是TearDownSuite")
}

func (s *DemoSuite) TestCases() {
	testCases := []struct {
		name string
	}{
		{
			name: "测试1",
		},
		{
			name: "测试2",
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			t.Log("这是测试TestCases", tc.name)
		})
	}
}

func (s *DemoSuite) TestSub() {
	s.T().Run("子测试", func(t *testing.T) {
		t.Log("子测试1")
	})

	s.T().Run("子测试", func(t *testing.T) {
		t.Log("子测试2")
	})
}

func TestDemo(t *testing.T) {
	suite.Run(t, new(DemoSuite))
}
