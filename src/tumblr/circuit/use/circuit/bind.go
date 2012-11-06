package circuit

import (
	"tumblr/circuit/kit/join"
	"tumblr/circuit/sys/lang/types"
)

var (
	link       = join.SetThenGet{Name: "circuit language runtime"}
	hostparser = join.SetThenGet{Name: "host parser"}
)

func Bind(v runtime, hp hostParser) {
	link.Set(v)
	hostparser.Set(hp)
}

func get() runtime {
	return link.Get().(runtime)
}

func getHostParser() hostParser {
	return hostparser.Get().(hostParser)
}

// Operators

func ParseHost(h string) Host {
	return getHostParser().Parse(h)
}

func RegisterType(v interface{}) {
	types.RegisterType(v)
}

func RegisterFunc(fn Func) {
	types.RegisterFunc(fn)
}

// Runtime-relative funcs

func Ref(v interface{}) X {
	return get().Ref(v)
}

func PermRef(v interface{}) XPerm {
	return get().PermRef(v)
}

func XAddr() Addr {
	return get().XAddr()
}

func SetBoot(v interface{}) {
	get().SetBoot(v)
}

func Spawn(host Host, anchor []string, fn Func, in ...interface{}) ([]interface{}, Addr, error) {
	return get().Spawn(host, anchor, fn, in...)
}

func Daemonize(fn func()) {
	get().Daemonize(fn)
}

func Kill(addr Addr) error {
	return get().Kill(addr)
}

func Dial(addr Addr, service string) X {
	return get().Dial(addr, service)
}

func Listen(service string, receiver interface{}) {
	get().Listen(service, receiver)
}

func TryDial(addr Addr, service string) (X, error) {
	return get().TryDial(addr, service)
}

func Export(val ...interface{}) interface{} {
	return get().Export(val...)
}

func Import(exported interface{}) ([]interface{}, string, error) {
	return get().Import(exported)
}
