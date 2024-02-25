package infra

import (
	"fmt"
	"strings"

	glog "github.com/flashmob/go-guerrilla/log"
	"golang.org/x/exp/slices"
)

type RootLogger glog.Logger

func NewRootLogger(rootCfg RootConfig) (RootLogger, error) {
	dest := strings.ToLower(rootCfg.LogDest)
	level := strings.ToLower(rootCfg.LogLevel)

	if !slices.Contains(glogDests, dest) {
		return nil, fmt.Errorf("invalid logging destination, got: %s, expected: %v", dest, glogDests)
	}
	if !slices.Contains(glogLevels, level) {
		return nil, fmt.Errorf("invalid logging level, got: %s, expected: %v", level, glogLevels)
	}

	log, err := glog.GetLogger(dest, level)
	return log, err
}

func NewLogger(root RootLogger, name string) LogFacade {
	return facadeImpl{
		root: root,
		staticFields: map[string]any{
			"logger.name": name,
		},
	}
}

type LogFacade interface {
	Debugf(msg string, args ...any)
	Infof(msg string, args ...any)
	Errorf(err error, msg string, args ...any)
}

type facadeImpl struct {
	root         RootLogger
	staticFields map[string]any
}

func (f facadeImpl) fields() map[string]any {
	return f.staticFields
}

func (f facadeImpl) Debugf(msg string, args ...any) {
	f.root.WithFields(f.fields()).Debugf(msg, args...)
}

func (f facadeImpl) Infof(msg string, args ...any) {
	f.root.WithFields(f.fields()).Infof(msg, args...)
}

func (f facadeImpl) Errorf(err error, msg string, args ...any) {
	f.root.WithFields(f.fields()).WithError(err).Errorf(msg, args...)
}

var (
	glogDests = []string{
		"off",
		"stdout",
		"stderr",
	}

	glogLevels = []string{
		"info",
		"panic",
		"fatal",
		"error",
		"warn",
		"info",
		"debug",
		"trace",
	}
)
