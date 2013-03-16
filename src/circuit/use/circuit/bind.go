// Package circuit exposes the core functionalities provided by the circuit programming environment
package circuit

import (
	"circuit/kit/join"
	"circuit/sys/lang/types"
)

var link = join.SetThenGet{Name: "circuit language runtime"}

// Bind is used internally.
func Bind(v runtime) {
	link.Set(v)
}

func get() runtime {
	return link.Get().(runtime)
}

// Operators

// RegisterValue regiseters the type of v so that cross-worker pointers to local values of the same type as v
// can be passed as arguments or return values of cross-worker function invokations. RegisterValue should
// typically be invoked within the init function of the package that defines the type.
func RegisterValue(v interface{}) {
	types.RegisterValue(v)
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

// WorkerAddr returns the address of this worker.
func WorkerAddr() Addr {
	return get().WorkerAddr()
}

func setBoot(v interface{}) {
	get().SetBoot(v)
}

// Spawn starts a new worker process on host; the worker is registered under
// all directories in the anchor file system named by anchor; the worker
// function fn is executed on the newly spawned worker with arguments in.
// Unless Daemonized is called during the execution of fn, the worker will be
// killed as soon as fn completes.
// Unless an error occurs, err equals nil, addr is the address of the newly
// spawened worker and retrn holds the values returned by fn.
func Spawn(host string, anchor []string, fn Func, in ...interface{}) (retrn []interface{}, addr Addr, err error) {
	return get().Spawn(host, anchor, fn, in...)
}

// Daemonize instructs the circuit runtime that this worker should not exit until fn completes,
// in addition to any other functions that might be daemonized already.
// Daemonize can be called only from within a worker function.
func Daemonize(fn func()) {
	get().Daemonize(fn)
}

// Kill kills the process of the worker with address addr.
func Kill(addr Addr) error {
	return get().Kill(addr)
}

// Dial contacts the worker specified by addr and requests a cross-worker
// pointer to service. If service is not being listened to at addr, nil is
// returned. Physical failures to contact the worker result in a panic.
func Dial(addr Addr, service string) X {
	return get().Dial(addr, service)
}

func DialSelf(service string) interface{} {
	return get().DialSelf(service)
}

// Listen registers receiver as a ..
func Listen(service string, receiver interface{}) {
	get().Listen(service, receiver)
}

// TryDial behaves like Dial, with the exception that instead of panicking 
// in the event of physical issues, an error is returned instead.
func TryDial(addr Addr, service string) (X, error) {
	return get().TryDial(addr, service)
}

func Export(val ...interface{}) interface{} {
	return get().Export(val...)
}

func Import(exported interface{}) ([]interface{}, string, error) {
	return get().Import(exported)
}
