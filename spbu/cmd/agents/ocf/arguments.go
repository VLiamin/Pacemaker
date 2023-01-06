package ocf

type (
	EnvironmentVariables interface {
		Get(string) (string, bool)
	}
)

type (
	Arguments struct {
		args []string
		env  EnvironmentVariables
	}
)

func (a Arguments) Get(key string) (string, bool) { return a.env.Get(`OCF_RESKEY_` + key) }
