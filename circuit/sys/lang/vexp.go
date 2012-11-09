package lang

import (
	"reflect"
	"tumblr/circuit/use/circuit"
)

type exportGroup struct {
	PtrPtr []*ptrPtrMsg
}

func (r *Runtime) exportValues(values []interface{}, importer circuit.Addr) ([]interface{}, []*ptrPtrMsg) {
	eg := &exportGroup{}
	rewriter := func(src, dst reflect.Value) bool {
		return r.exportRewrite(src, dst, importer, eg)
	}
	return rewriteInterface(rewriter, values).([]interface{}), eg.PtrPtr
}

func (r *Runtime) exportRewrite(src, dst reflect.Value, importer circuit.Addr, eg *exportGroup) bool {
	// Serialize cross-runtime pointers
	switch v := src.Interface().(type) {

	case *_permptr:
		pm := &permPtrPtrMsg{ID: v.impHandle().ID, TypeID: v.impHandle().Type.ID, Src: v.impHandle().Exporter}
		dst.Set(reflect.ValueOf(pm))
		return true

	case *_ptr:
		if importer == nil {
			panic("exporting non-perm ptrptr without importer")
		}
		pm := &ptrPtrMsg{ID: v.impHandle().ID, Src: v.impHandle().Exporter}
		dst.Set(reflect.ValueOf(pm))
		eg.PtrPtr = append(eg.PtrPtr, pm)
		return true

	case *_ref:
		if importer == nil {
			panic("exporting non-perm ptr without importer")
		}
		dst.Set(reflect.ValueOf(r.exportPtr(v.value, importer)))
		return true

	case *_permref:
		dst.Set(reflect.ValueOf(r.exportPtr(v.value, nil)))
		return true
	}

	return false
}

// If importer is nil, a permanent ptr is exported
func (r *Runtime) exportPtr(v interface{}, importer circuit.Addr) interface{} {
	exph := r.exp.Add(v, importer)

	if importer == nil {
		return &permPtrMsg{ID: exph.ID, TypeID: exph.Type.ID}
	}

	// Monitor the importer for liveness.
	// DropPtr the handles upon importer death.
	r.lk.Lock()
	defer r.lk.Unlock()
	_, ok := r.live[importer]
	if !ok {
		r.live[importer] = struct{}{}
		go func() {
			defer func() {
				r.lk.Lock()
				delete(r.live, importer)
				r.lk.Unlock()
				// DropPtr/forget all exported handles
				r.exp.RemoveImporter(importer)
			}()

			conn, err := r.dialer.Dial(importer)
			if err != nil {
				return
			}
			defer conn.Close()

			if conn.Write(&dontReplyMsg{}) != nil {
				return
			}
			// Read returns when the remote dies and 
			// runs the conn into an error
			conn.Read()
		}()
	}
	return &ptrMsg{ID: exph.ID, TypeID: exph.Type.ID}
}
