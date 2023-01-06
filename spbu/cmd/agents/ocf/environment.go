package ocf

import (
	"os"
)

type (
	OSEnvironment struct{}
)

func (env OSEnvironment) Get(key string) (string, bool) { return os.LookupEnv(key) }
