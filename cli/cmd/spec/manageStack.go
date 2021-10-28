package spec

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	log "github.com/infracloudio/krius/pkg/logger"
	"github.com/spf13/cobra"
)

type AppRunner struct {
	log    *log.Logger
	status *Status
}

var debug bool

func NewAppRunner(log *log.Logger, status *Status) *AppRunner {
	return &AppRunner{
		log:    log,
		status: status,
	}
}

func manageApp(cmd *cobra.Command) error {
	logger := log.NewLogger(debug)
	runner := NewAppRunner(logger, NewStatus(logger))
	switch cmd.Use {
	case "apply":
		err := runner.applySpec(cmd)
		if err != nil {
			logger.Error(err.Error())
		}
	case "destroy":
		err := runner.uninstallSpec(cmd)
		if err != nil {
			logger.Error(err.Error())
		}

	}
	return nil
}
func (s *Status) Start(msg string) {
	s.Stop()
	s.message = msg
	if s.spinner != nil {
		s.spinner.Suffix = fmt.Sprintf(" %s ", s.message)
		s.spinner.Start()
	}
}

func (s *Status) Success(msg ...string) {
	if !s.spinner.Active() {
		return
	}
	if s.spinner != nil {
		s.spinner.Stop()
	}
	if msg != nil {
		s.logger.Infof(s.successStatusFormat, strings.Join(msg, " "))
		return
	}
	s.logger.Infof(s.successStatusFormat, s.message)
	s.message = ""
}
func (s *Status) Error(msg ...string) {
	if !s.spinner.Active() {
		return
	}
	if s.spinner != nil {
		s.spinner.Stop()
	}
	if msg != nil {
		s.logger.Infof(s.failureStatusFormat, strings.Join(msg, " "))
		return
	}
	s.logger.Infof(s.failureStatusFormat, s.message)
	s.message = ""
}
func (s *Status) Stop() {
	if !s.spinner.Active() {
		return
	}
	if s.spinner != nil {
		fmt.Fprint(s.logger.Writer, "\r")
		s.spinner.Stop()
	}
}

type Status struct {
	spinner *spinner.Spinner
	message string
	logger  *log.Logger

	successStatusFormat string
	failureStatusFormat string
}

func NewStatus(logger *log.Logger) *Status {
	return &Status{
		logger:              logger,
		spinner:             spinner.New(defaultCharSet, defaultDelay, spinner.WithWriter(os.Stdout)),
		successStatusFormat: successStatusFormat,
		failureStatusFormat: failureStatusFormat,
	}
}

var (
	defaultCharSet = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	defaultDelay   = 100 * time.Millisecond
)

const (
	successStatusFormat = " \x1b[32m✓\x1b[0m %s"
	failureStatusFormat = " \x1b[31m✗\x1b[0m %s"
)
