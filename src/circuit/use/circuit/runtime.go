package circuit

type runtime interface {
	// Low-level
	WorkerAddr() Addr
	SetBoot(interface{})
	Kill(Addr) error

	// Spawn mechanism
	Spawn(Host, []string, Func, ...interface{}) ([]interface{}, Addr, error)
	Daemonize(func())

	// Cross-services
	Dial(Addr, string) X
	DialSelf(string) interface{}
	TryDial(Addr, string) (X, error)
	Listen(string, interface{})

	// Persistence of XPerm values
	Export(...interface{}) interface{}
	Import(interface{}) ([]interface{}, string, error)

	// Cross-interfaces
	Ref(interface{}) X
	PermRef(interface{}) XPerm
}
