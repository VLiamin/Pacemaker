package ocf

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
)

func runAction(ctx Context, action string, agent Agent, args Arguments) error {
	switch action {
	case ActionProbe:
		return agent.Probe(ctx, args)
	case ActionStart:
		return agent.Start(ctx, args)
	case ActionStop:
		return agent.Stop(ctx, args)
	case ActionStatus, ActionMonitor:
		if ctx.Interval() == 0 {
			return agent.Probe(ctx, args)
		}
		return agent.Monitor(ctx, args)
	case ActionPromote:
		return agent.Promote(ctx, args)
	case ActionDemote:
		return agent.Demote(ctx, args)
	case ActionNotify:
		return agent.Notify(ctx, args)
	case ActionMetaData:
		xml.NewEncoder(os.Stdout).Encode(agent.MetaData(ctx))
		return ErrSuccess
	default:
		return fmt.Errorf(`unexpected action %s. %w`, action, ErrGeneric)
	}
}

func getExitReason(err error) (Error, string) {
	if ocfErr, isOcfErr := err.(Error); isOcfErr {
		return ocfErr, ocfErr.Error()
	}

	if ocfErr, isOcfErr := errors.Unwrap(err).(Error); isOcfErr {
		return ocfErr, err.Error()
	}

	if err != nil {
		return ErrGeneric, err.Error()
	}

	return ErrSuccess, ErrSuccess.Error()
}

func Run(agent Agent, args []string, environ EnvironmentVariables) Error {
	ctx, cancel := context.WithTimeout(context.Background(), Context{nil, environ}.Timeout())
	defer cancel()

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "ocf_exit_reason: no action is specified\n")
		return ErrGeneric
	}

	err := runAction(Context{ctx, environ}, args[0], agent, Arguments{
		args, environ,
	})

	ocfErr, exitReason := getExitReason(err)
	if ErrSuccess != ocfErr {
		fmt.Fprintf(os.Stderr, "ocf_exit_reason: %s\n", exitReason)
	}

	return ocfErr
}
