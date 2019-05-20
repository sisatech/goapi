package graphqlws

import "fmt"

// Logger provides an interface allowing the library's internal logging to be
// handled by third party logging libraries.
type Logger interface {

	// Info logs a general information message. This is verbose and should
	// be discarded except when troubleshooting.
	Info(message string)

	// Error logs a message describing a non-critical error that has
	// occurred. This is not verbose but is also for troubleshooting
	// purposes. It is okay to discard these messages because the library
	// only logs errors that can be recovered or ignored. Critical errors
	// will always be returned as an error somewhere.
	Error(message string)
}

type logger struct {
	logger Logger
}

func (l logger) Info(message string, v ...interface{}) {
	if l.logger == nil {
		return
	}
	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}
	l.logger.Info(message)
}

func (l logger) Error(message string, v ...interface{}) {
	if l.logger == nil {
		return
	}
	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}
	l.logger.Error(message)
}

type operationLogger struct {
	logger logger
	suffix string
}

func (l operationLogger) Info(message string, v ...interface{}) {
	l.logger.Error(message+l.suffix, v...)
}

func (l operationLogger) Error(message string, v ...interface{}) {
	l.logger.Error(message+l.suffix, v...)
}
