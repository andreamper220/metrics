package handlers_test

import (
	"github.com/andreamper220/metrics.git/internal/server/application"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	Server *httptest.Server
}

func (s *HandlerTestSuite) SetupTest() {
	s.Require().NoError(os.Setenv("KEY", "test_key"))
	s.Require().NoError(application.ParseFlags())

	s.Require().NoError(application.Run(true))

	s.Server = httptest.NewServer(application.MakeRouter())
}

func TestHandlersSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
