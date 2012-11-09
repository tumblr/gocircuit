package circuit

// X is a local proxy to a Go value residing on a remote runtime.
// It stands for cross-runtime (X) pointer (Ptr).
//
// X pointers to remote values are garbage-collected: When all Ptrs pointing
// to a receiver object have been abandoned (collected by their respective Go
// runtimes' garbage collector), the receiver's circuit runtime is notified and
// the receiver is released, unless still utilized in its native runtime.
//
// X is a user-facing interface.
type X interface {

	// Addr returns the address of the runtime, hosting the underlying object
	Addr() Addr

	// Call invokes the method named proc on the (possibly) remote circuit
	// runtime where the receiver resides and delivers the returned value
	// locally.
	// If the receiver runtime is unreachable, Call panics.
	Call(proc string, in ...interface{}) []interface{}

	// IsX is a nop. It is present to help distinguish Ptr values more clearly.
	IsX()

	String() string
}

// XPerm is much like X, except it is not garbage-collected.
//
// When a native object is permanently exported (using PermRef), resulting in a
// XPerm at the receiving runtime, the object can never be freed.
// Consequently, XPerms are not monitored by the circuit's garbage
// collector.
type XPerm interface {

	// PermPtr peruses Ptr's interface.
	// In particular, applications can treat a PermPtr as a Ptr semantically.
	X

	// IsPerm is a nop, used to distinguish perm and non-perm ptrs.
	IsXPerm()
}

//  Why XRef and XPermRef return X and XPerm?
//
//  By definition, XRef returns a special value that encloses a user value
//  together with an annotation instructing the circuit runtime to pass it as a
//  cross-runtime pointer, X. Let's call the type of this special value
//  XLicense.  (This type is fictitious for the purposes of this discussion.)
//  We would then define XRef as:
//
//		func XRef(...) XLicense { ... }
//
//  Now consider a common use pattern of XRef:
//
//		func F() X { return XRef(...) }
//
//  F creates some native Go value and returns it as an XLicense. F here is
//  intended to be called from another runtime, so that the returned value
//  arrives at the caller as a cross-runtime pointer, X.
//
//  From the caller's point of view, F should return an X value. Therefore,
//  XLicense must be the same as X.
//
//  Note that this line of reasoning is a hack made necessary by the fact
//  that we are using a naked runtime. In particular, because in native Go
//  a func can only have one signature, whereas in the circuit a func has
//  different signatures depending on the context: caller or callee.
//  This will be remedied (made transparent to the user) by a specialized
//  circuit compiler.

// Func is a symbolic type that refurs to circuit func types, which
// are types with one public method.
type Func interface{}
