package testingh

import (
	"context"

	"github.com/stretchr/testify/suite"
)

type ContextSuite struct {
	suite.Suite

	Ctx           context.Context
	CtxCancelFunc context.CancelFunc
}

func (cs *ContextSuite) SetupTest() {
	cs.Ctx, cs.CtxCancelFunc = context.WithCancel(context.Background())
}

func (cs *ContextSuite) TearDownTest() {
	cs.CtxCancelFunc()
}
