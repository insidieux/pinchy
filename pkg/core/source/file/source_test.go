package file

import (
	"context"
	"testing"

	"github.com/insidieux/pinchy/pkg/core"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

// --- Tests ---

func TestNewSource(t *testing.T) {
	suite.Run(t, new(newSourceTestSuite))
}

func TestSource_Fetch(t *testing.T) {
	suite.Run(t, new(sourceFetchTestSuite))
}

func TestSource_WithLogger(t *testing.T) {
	suite.Run(t, new(sourceWithLoggerTestSuite))
}

// --- Suites ---

type newSourceTestSuite struct {
	suite.Suite
}

func (s *newSourceTestSuite) TestNewSource() {
	got := NewSource(nil, `filename`)
	s.Implements((*core.Source)(nil), got)
	s.Equal(&Source{nil, `filename`, nil}, got)
}

type sourceFetchTestSuite struct {
	suite.Suite
	source *Source
	reader afero.Afero
	hook   *test.Hook
}

func (s *sourceFetchTestSuite) SetupTest() {
	s.reader = afero.Afero{Fs: afero.NewMemMapFs()}
	s.source = NewSource(s.reader, `filename`)
	s.source.logger, s.hook = test.NewNullLogger()
}

func (s *sourceFetchTestSuite) TestErrorRead() {
	services, err := s.source.Fetch(context.Background())
	s.Nil(services)
	s.Error(err)
	s.EqualError(err, `failed read content from config file: open filename: file does not exist`)
}

func (s *sourceFetchTestSuite) TestErrorUnmarshal() {
	inMemoryFile, err := s.reader.Create(string(s.source.filename))
	if err != nil {
		panic(errors.Wrap(err, `failed to create in-memory file`))
	}
	if _, err := inMemoryFile.WriteString(`{"key": "value"}`); err != nil {
		panic(errors.Wrap(err, `failed to write to in-memory file`))
	}

	services, err := s.source.Fetch(context.Background())
	s.Nil(services)
	s.Error(err)
	s.Contains(err.Error(), `failed unmarshal content from config file`)
}

func (s *sourceFetchTestSuite) TestSkipServiceValidationCase() {
	inMemoryFile, err := s.reader.Create(string(s.source.filename))
	if err != nil {
		panic(errors.Wrap(err, `failed to create in-memory file`))
	}
	serviceBytes, err := yaml.Marshal(core.Services{
		{
			Name: `service`,
		},
	})
	if err != nil {
		panic(errors.Wrap(err, `failed to marshal service`))
	}
	if _, err := inMemoryFile.Write(serviceBytes); err != nil {
		panic(errors.Wrap(err, `failed to write to in-memory file`))
	}

	services, err := s.source.Fetch(context.Background())
	s.NotNil(services)
	s.NoError(err)
	s.Equal(core.Services{}, services)
	s.Equal(s.hook.LastEntry().Level, logrus.WarnLevel)
	s.Equal(s.hook.LastEntry().Message, `Failed to validate service #0: service "service" field "address" is required and cannot be empty`)
}

func (s *sourceFetchTestSuite) TestSuccess() {
	inMemoryFile, err := s.reader.Create(string(s.source.filename))
	if err != nil {
		panic(errors.Wrap(err, `failed to create in-memory file`))
	}
	expected := core.Services{
		{
			Name:    `service-1`,
			Address: `127.0.0.1`,
		},
		{
			Name:    `service-2`,
			Address: `127.0.0.2`,
		},
	}
	servicesBytes, err := yaml.Marshal(expected)
	if err != nil {
		panic(errors.Wrap(err, `failed to marshal services`))
	}
	if _, err := inMemoryFile.Write(servicesBytes); err != nil {
		panic(errors.Wrap(err, `failed to write to in-memory file`))
	}

	services, err := s.source.Fetch(context.Background())
	s.NotNil(services)
	s.NoError(err)
	s.Equal(expected, services)
}

type sourceWithLoggerTestSuite struct {
	suite.Suite
}

func (s *sourceWithLoggerTestSuite) TestWithLogger() {
	logger, _ := test.NewNullLogger()
	src := NewSource(nil, `filename`)
	src.WithLogger(logger)
}
