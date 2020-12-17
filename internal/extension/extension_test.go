package extension

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestRegisterError_String(t *testing.T) {
	suite.Run(t, new(registerErrorStringTestSuite))
}

func TestRegisterError_Error(t *testing.T) {
	suite.Run(t, new(registerErrorErrorTestSuite))
}

// --- Suites ---

type registerErrorStringTestSuite struct {
	suite.Suite
	err RegisterError
}

func (s *registerErrorStringTestSuite) SetupTest() {
	s.err = RegisterError{}
}

func (s *registerErrorStringTestSuite) TestEmptyError() {
	s.Equal(``, s.err.String())
}

func (s *registerErrorStringTestSuite) TestNonEmptyError() {
	s.err = append(s.err, errors.New(`expected error 1`))
	s.err = append(s.err, errors.New(`expected error 2`))
	s.Equal(`expected error 1; expected error 2`, s.err.String())
}

type registerErrorErrorTestSuite struct {
	suite.Suite
	err RegisterError
}

func (s *registerErrorErrorTestSuite) SetupTest() {
	s.err = RegisterError{}
}

func (s *registerErrorErrorTestSuite) TestEmptyError() {
	s.Equal(``, s.err.String())
}

func (s *registerErrorErrorTestSuite) TestNonEmptyError() {
	s.err = append(s.err, errors.New(`expected error 1`))
	s.err = append(s.err, errors.New(`expected error 2`))
	s.Equal(`expected error 1; expected error 2`, s.err.Error())
}
