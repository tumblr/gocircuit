package lang

/*

	FORKING A GO ROUTINE ON A REMOTE RUNTIME

		import . "circuit/use/circuit"

		type MyFunc struct{}
		func (MyFunc) AnyName(anyArg anyType) (anyReturn anyType) {
			...
		}
		func init() { types.RegisterFunc(MyFunc{}) }

		func main() {
			Go(conn, MyFunc{}, a1)
		}

*/
