package file

import (
	"context"

	"github.com/insidieux/pinchy/pkg/core"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type (
	// Reader tries to read content from file by name
	Reader interface {
		ReadFile(name string) ([]byte, error)
	}

	// Path is custom type for file path string
	Path string

	// Source is implementation of core.Source interface
	Source struct {
		reader   Reader
		filename Path
		logger   core.LoggerInterface
	}
)

// NewSource provide Source as core.Source implementation
func NewSource(reader Reader, filename Path) *Source {
	return &Source{
		reader:   reader,
		filename: filename,
	}
}

// Fetch provide information about core.Services from file
// - call Reader.ReadFile
// - yaml.Unmarshal contents
// - validate core.Service
// - return core.Services
func (s *Source) Fetch(ctx context.Context) (core.Services, error) {
	s.logger.Infof(`Reading file "%s"`, s.filename)
	contents, err := s.reader.ReadFile(string(s.filename))
	if err != nil {
		return nil, errors.Wrap(err, `failed read content from config file`)
	}

	s.logger.Infoln(`Decoding yml config`)
	items := make([]*core.Service, 0)
	if err := yaml.Unmarshal(contents, &items); err != nil {
		return nil, errors.Wrap(err, `failed unmarshal content from config file`)
	}

	s.logger.Infoln(`Collecting services list with service validation`)
	result := make([]*core.Service, 0)
	for index, item := range items {
		if err := item.Validate(ctx); err != nil {
			s.logger.Warningln(errors.Wrapf(err, `Failed to validate service #%d`, index).Error())
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

// WithLogger is implementation of core.Loggable interface
func (s *Source) WithLogger(logger core.LoggerInterface) {
	s.logger = logger
}
