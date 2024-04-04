package handlers_test

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/andreamper220/metrics.git/internal/server"
)

type HandlerTestSuite struct {
	suite.Suite
	Server *httptest.Server
}

func (s *HandlerTestSuite) SetupTest() {
	r := server.MakeRouter()
	s.Server = httptest.NewServer(r)
}

func TestHandlersSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
