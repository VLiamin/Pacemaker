package main

import (
	"time"

	"yadro.com/m.konovalov/pcmkr/cmd/agents/ocf"
)

type (
	Agent struct {
		ocf.Pseudo
	}
)

func (agent Agent) MetaData(ocf.Context) ocf.MetaData {
	return ocf.MetaData{
		Name:            `traid-config`,
		Version:         `0.0.0`,
		LongDesription:  `demo-resource`,
		ShortDesription: `demo-resource`,
		Parameters: ocf.Parameters{{
			Name:     `demo_parameter_1`,
			Uniq:     false,
			Required: false,
			Kind:     `integer`,
			Default:  `10`,
		}, {
			Name:     `demo_parameter_2`,
			Uniq:     false,
			Required: false,
			Kind:     `boolean`,
			Default:  `false`,
		}},
		Actions: ocf.Actions{{
			Name:    ocf.ActionStart,
			Timeout: 20 * time.Second,
		}, {
			Name:    ocf.ActionStop,
			Timeout: 20 * time.Second,
		}, {
			Name:     ocf.ActionMonitor,
			Interval: 20 * time.Second,
			Timeout:  20 * time.Second,
		}},
	}
}
