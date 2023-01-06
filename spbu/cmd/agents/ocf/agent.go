package ocf

type (
	Agent interface {
		MetaData(Context) MetaData
		Validate(Context, Arguments) error
		Probe(Context, Arguments) error
		Start(Context, Arguments) error
		Stop(Context, Arguments) error
		Reload(Context, Arguments) error
		Monitor(Context, Arguments) error
		Promote(Context, Arguments) error
		Demote(Context, Arguments) error
		Notify(Context, Arguments) error
	}
)

type (
	NilAgent struct{}
)

func (agent NilAgent) MetaData(Context) MetaData         { return MetaData{} }
func (agent NilAgent) Validate(Context, Arguments) error { return ErrNotImplemented }
func (agent NilAgent) Probe(Context, Arguments) error    { return ErrNotImplemented }
func (agent NilAgent) Start(Context, Arguments) error    { return ErrNotImplemented }
func (agent NilAgent) Stop(Context, Arguments) error     { return ErrNotImplemented }
func (agent NilAgent) Reload(Context, Arguments) error   { return ErrNotImplemented }
func (agent NilAgent) Monitor(Context, Arguments) error  { return ErrNotImplemented }
func (agent NilAgent) Promote(Context, Arguments) error  { return ErrNotImplemented }
func (agent NilAgent) Demote(Context, Arguments) error   { return ErrNotImplemented }
func (agent NilAgent) Notify(Context, Arguments) error   { return ErrNotImplemented }
