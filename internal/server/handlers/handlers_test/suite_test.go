package handlers_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/andreamper220/metrics.git/internal/server"
)

type HandlerTestSuite struct {
	suite.Suite
	Server *httptest.Server
}

func (s *HandlerTestSuite) SetupTest() {
	s.Require().NoError(os.Setenv("KEY", "test_key"))
	server.ParseFlags()

	s.Require().NoError(server.Run(true))

	s.Server = httptest.NewServer(server.MakeRouter())
}

func TestHandlersSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
