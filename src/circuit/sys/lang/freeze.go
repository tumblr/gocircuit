package lang

import (
	"runtime/debug"
)

func (r *Runtime) Export(val ...interface{}) interface{} {
	expHalt, _ := r.exportValues(val, nil)
	return &exportedMsg{
		Value: expHalt,
		Stack: string(debug.Stack()),
	}
}

func (r *Runtime) Import(exported interface{}) ([]interface{}, string, error) {
	h, ok := exported.(*exportedMsg)
	if !ok {
		return nil, "", NewError("foreign saved message (msg=%T)", exported)
	}
	val, err := r.importValues(h.Value, nil, nil, false, nil)
	if err != nil {
		return nil, "", err
	}
	return val, h.Stack, nil
}
