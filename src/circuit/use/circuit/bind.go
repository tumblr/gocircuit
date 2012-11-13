package circuit

import (
	"circuit/kit/join"
	"circuit/sys/lang/types"
)

var link = join.SetThenGet{Name: "circuit language runtime"}

func Bind(v runtime) {
	link.Set(v)
}

func get() runtime {
	return link.Get().(runtime)
}

// Operators

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

func WorkerAddr() Addr {
	return get().WorkerAddr()
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
