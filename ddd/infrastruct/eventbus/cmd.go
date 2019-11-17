package eventbus

import (
	"errors"
)

var busCommandPublisher CommandPublisher
var busCommandSubscriber CommandSubscriber

type Command interface {
	Action() string
	Data() interface{}
}

type CommandHandler func(cmd Command) CommandResult

type CommandResult interface {
	Result() interface{}
	Error() error
}

type SimpleCommandResult struct {
	result interface{}
	err    error
}

func (c *SimpleCommandResult) Result() interface{} {
	return c.result
}

func (c *SimpleCommandResult) Error() error {
	return c.err
}

func NewCommandResult(res interface{}, err error) CommandResult {
	return &SimpleCommandResult{
		res,
		err,
	}
}

type Feature chan CommandResult

type CommandPublisher interface {
	SendCommand(cmd Command) CommandResult
	SendCommandAsync(cmd Command) (Feature, error)
}

type CommandSubscriber interface {
	RegisterCommandHandler(action string, handler CommandHandler)
}

func SendCommand(cmd Command) CommandResult {
	if busCommandPublisher != nil {
		return busCommandPublisher.SendCommand(cmd)
	}
	return nil
}

func SendCommandAsync(cmd Command) (Feature, error) {
	if busCommandPublisher != nil {
		return busCommandPublisher.SendCommandAsync(cmd)
	}
	return nil, errors.New(" Bus not registed ")
}

func RegisterCommandHandler(action string, handler CommandHandler) {
	if busCommandSubscriber != nil {
		busCommandSubscriber.RegisterCommandHandler(action, handler)
	}
}

type SimpleCommand struct {
	action string
	data   interface{}
}

func NewCommand(action string, data interface{}) Command {
	return &SimpleCommand{
		action,
		data,
	}
}

func (s *SimpleCommand) Action() string {
	return s.action
}

func (s *SimpleCommand) Data() interface{} {
	return s.data
}
