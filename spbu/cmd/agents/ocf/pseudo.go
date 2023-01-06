package ocf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type (
	Pseudo struct {
		NilAgent
	}

	PseudoState struct {
		Time   time.Time `json:"time"`
		Action string    `json:"action"`
	}
)

func (agent Pseudo) Probe(ctx Context, args Arguments) error { return agent.Monitor(ctx, args) }

func (agent Pseudo) saveState(ctx Context, state PseudoState) error {
	dirName := filepath.Join(ctx.StateDir(), ctx.ResourceKind())
	fileName := filepath.Join(dirName, ctx.ResourceInstance()) + `-pseudo.json`

	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return fmt.Errorf(`failed to create state file %s directory. %w`, fileName, err)
	}

	stateFile, err := os.OpenFile(
		fileName,
		os.O_CREATE|os.O_TRUNC|os.O_RDWR,
		0660,
	)
	if err != nil {
		return fmt.Errorf(`failed to create state file %s. %w`, fileName, err)
	}
	defer stateFile.Close()

	if err = json.NewEncoder(stateFile).Encode(state); err != nil {
		return fmt.Errorf(`failed to create state file %s. %w`, fileName, err)
	}

	if err = stateFile.Sync(); err != nil {
		return fmt.Errorf(`failed to create state file %s. %w`, fileName, err)
	}

	return ErrSuccess
}

func (agent Pseudo) Start(ctx Context, _ Arguments) error {
	return agent.saveState(ctx, PseudoState{
		time.Now().UTC(),
		ActionStart,
	})
}

func (agent Pseudo) dropState(ctx Context) error {
	fileName := filepath.Join(ctx.StateDir(), ctx.ResourceKind(), ctx.ResourceInstance()) + `-pseudo.json`

	err := os.Remove(fileName)
	if err == nil {
		return ErrSuccess
	}

	if os.IsNotExist(err) {
		return fmt.Errorf(`no state file %s. %w`, fileName, ErrNotRunning)
	}
	return fmt.Errorf(`failed to remove state file %s. %w`, fileName, err)
}

func (agent Pseudo) Stop(ctx Context, _ Arguments) error { return agent.dropState(ctx) }

func (agent Pseudo) loadState(ctx Context) (PseudoState, error) {
	fileName := filepath.Join(ctx.StateDir(), ctx.ResourceKind(), ctx.ResourceInstance()) + `-pseudo.json`

	stateFile, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return PseudoState{}, ErrNotRunning
		}
		return PseudoState{}, fmt.Errorf(`failed to open state file %s. %w`, fileName, err)
	}
	defer stateFile.Close()

	state := PseudoState{}
	if err := json.NewDecoder(stateFile).Decode(&state); err == nil {
		return state, nil
	}

	return PseudoState{}, fmt.Errorf(`failed to read state file %s. %w`, fileName, err)
}

func (agent Pseudo) Monitor(ctx Context, _ Arguments) error {
	state, err := agent.loadState(ctx)

	if err != nil {
		return err
	}

	switch state.Action {
	case ActionStart, ActionDemote:
		return ErrSuccess

	case ActionPromote:
		return ErrRunningMaster
	}

	return ErrGeneric
}

func (agent Pseudo) Promote(ctx Context, _ Arguments) error {
	return agent.saveState(ctx, PseudoState{
		time.Now().UTC(),
		ActionPromote,
	})
}

func (agent Pseudo) Demote(ctx Context, _ Arguments) error {
	return agent.saveState(ctx, PseudoState{
		time.Now().UTC(),
		ActionDemote,
	})
}
