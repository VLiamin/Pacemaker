package ocf

import (
	"context"
	"strconv"
	"time"
)

type (
	Context struct {
		context.Context
		env EnvironmentVariables
	}
)

func (c Context) Pseudo() Pseudo {
	return Pseudo{}
}

func (c Context) parseDuration(key string) time.Duration {
	value, hasValue := c.env.Get(key)
	if !hasValue {
		return time.Duration(0)
	}

	timeout, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Duration(-1)
	}

	return time.Duration(time.Duration(timeout) * time.Millisecond)
}

func (c Context) Interval() time.Duration {
	return c.parseDuration(`OCF_RESKEY_CRM_meta_interval`)
}

func (c Context) Timeout() time.Duration {
	return c.parseDuration(`OCF_RESKEY_CRM_meta_timeout`)
}

func (c Context) NodeID() (value string) {
	value, _ = c.env.Get(`OCF_RESKEY_CRM_meta_on_node_uuid`)
	return
}

func (c Context) NodeName() (value string) {
	value, _ = c.env.Get(`OCF_RESKEY_CRM_meta_on_node`)
	return
}

func (c Context) ResourceKind() (value string) {
	value, _ = c.env.Get(`OCF_RESOURCE_TYPE`)
	return
}

func (c Context) ResourceInstance() (value string) {
	value, _ = c.env.Get(`OCF_RESOURCE_INSTANCE`)
	return
}

func (c Context) StateDir() (value string) {
	hasValue := false
	if value, hasValue = c.env.Get(`HA_RSCTMP`); !hasValue {
		value = "/run/resource-agents"
	}

	return
}
