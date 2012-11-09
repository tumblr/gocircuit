package lang

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"circuit/sys/lang/types"
)

// call invokes the method of r encoded by f with respect to t, with arguments a
func call(recv reflect.Value, t *types.TypeChar, id types.FuncID, arg []interface{}) (reply []interface{}, err error) {
	// Recover panic in user code and return it in error argument
	defer func() {
		p := recover()
		if p == nil {
			return
		}
		t := string(debug.Stack())
		switch q := p.(type) {
		case error:
			err = NewError(q.Error() + "\n" + t)
		default:
			err = NewError(fmt.Sprintf("%#v\n%s", q, t))
		}
	}()

	fn := t.Func[id]
	if fn == nil {
		return nil, NewError("no func")
	}
	av := make([]reflect.Value, 0, 1 + len(arg))
	av = append(av, recv)
	for _, a := range arg {
		av = append(av, reflect.ValueOf(a))
	}
	rv := fn.Method.Func.Call(av)
	reply = make([]interface{}, len(rv))
	for i, r := range rv {
		reply[i] = r.Interface()
	}
	return reply, nil
}
