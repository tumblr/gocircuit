package lang

// Addr is a unique representation of the identity of a remote runtime.
// The implementing type must be registered with gob.
type Addr interface {
	String() string

	// RuntimeID is tentatively part of the transport address so that, if
	// needed, we can verify the identity of the runtime we are talking to. 
	RuntimeID() RuntimeID
}

// Conn is a connection to a remote runtime.
type Conn interface {
	// The language runtime does not itself utilize timeouts on read/write
	// operations. Instead, it requires that calls to Read and Write be blocking
	// until success or irrecoverable failure is reached.
	//
	// The implementation of Conn must monitor the liveness of the remote
	// endpoint using an out-of-band (non-visible to the runtime) method. If
	// the endpoint is considered dead, all pending Read and Write request must
	// return with non-nil error.
	//
	// A non-nil error returned on any invokation of Read and Write signals to
	// the runtime that not just the connection, but the entire runtime
	// (identified by its address) behind the connection is dead.
	//
	// Such an event triggers various language runtime actions such as, for
	// example, releasing all values exported to that runtime. Therefore, a
	// typical Conn implementation might choose to attempt various physical
	// connectivity recovery methods, before it reports an error on any pending
	// connection. Such implentation strategies are facilitated by the fact
	// that the runtime has no semantic limits on the length of blocking waits.
	// In fact, the runtime has no notion of time altogether.

	// Read/Write operations must panic on any encoding/decoding errors.
	// Whereas they must return an error for any exernal (network) unexpected
	// conditions.  Encoding errors indicate compile-time errors (that will be
	// caught automatically once the system has its own compiler) but might be
	// missed by the bare Go compiler.
	//
	// Read/Write must be re-entrant.

	Read() (interface{}, error)
	Write(interface{}) error
	Close() error

	Addr() Addr
}

// Listener is a device for accepting incoming connections.
type Listener interface {
	Accept() Conn
	Close()
	Addr() Addr
}

// Dialer is a device for initating connections to addressed remote endpoints.
type Dialer interface {
	Dial(addr Addr) (Conn, error)
}

type Transport interface {
	Dialer
	Listener
}

// Host represents a host that can spawn circuit runtimes.
// Its specific implementation is up to the user.
// A Host must be gob encodable.
type Host interface{
	String() string
}
