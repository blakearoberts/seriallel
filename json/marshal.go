package json

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sync"
)

func Marshal(v any) ([]byte, error) {
	return marshal(reflect.ValueOf(v))
}

func marshal(v reflect.Value) ([]byte, error) {
	switch v.Kind() {
	case reflect.Pointer:
		return marshal(v.Elem())
	case reflect.Array, reflect.Slice:
		return marshalList(v)
	default:
		if !v.IsValid() {
			return []byte("null"), nil
		}
		return json.Marshal(v.Interface())
	}
}

func marshalList(v reflect.Value) ([]byte, error) {
	b := bytes.NewBuffer([]byte{'['})

	// serialize all items in parallel
	ch := make(chan []byte)
	wg := sync.WaitGroup{}
	for i := 0; i < v.Len(); i++ {
		wg.Add(1)
		go func(v reflect.Value) {
			defer wg.Done()
			bytes, err := marshal(v)
			if err != nil {
				// TODO
				panic(err)
			}
			ch <- bytes
		}(v.Index(i))
	}

	// close channel once all items have been serialized
	go func() {
		wg.Wait()
		close(ch)
	}()

	// write each item to buffer
	i := 0
	for s := range ch {
		if i > 0 {
			if err := b.WriteByte(','); err != nil {
				return nil, err
			}
		}
		if _, err := b.Write(s); err != nil {
			return nil, err
		}
		i++
	}
	wg.Wait()

	if err := b.WriteByte(']'); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
