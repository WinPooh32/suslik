// +build js

package file

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
)

func ReadAll(name string) ([]byte, error) {
	ch := make(chan error, 1)

	req := js.Global.Get("XMLHttpRequest").New()
	req.Set("responseType", "arraybuffer")
	req.Call("open", "GET", name, true)
	req.Call("addEventListener", "load", func(*js.Object) {
		go func() { ch <- nil }()
	}, false)
	req.Call("addEventListener", "error", func(o *js.Object) {
		go func() { ch <- &js.Error{Object: o} }()
	}, false)
	req.Call("send", nil)

	err := <-ch
	if err != nil {
		return nil, err
	}

	arraybuffer := req.Get("response")

	data, ok := js.Global.Get("Uint8Array").New(arraybuffer).Interface().([]byte)
	if !ok {
		return nil, fmt.Errorf("failed to get arraybuffer")
	}
	return data, nil
}
