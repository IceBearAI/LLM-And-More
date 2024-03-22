package util

import (
	"context"
	"reflect"
	"unsafe"
)

type iface struct {
	itab, data uintptr
}

type valueCtx struct {
	context.Context
	key, value any
}

// CopyContext. copy context
func CopyContext(ctx context.Context) context.Context {
	newCtx := context.Background()
	kv := make(map[any]any)
	getKeyValues(ctx, kv)
	for k, v := range kv {
		newCtx = context.WithValue(newCtx, k, v)
	}

	return newCtx
}

func getKeyValues(ctx context.Context, kv map[any]any) {
	rtType := reflect.TypeOf(ctx).String()

	// 遍历到最底层，返回
	if rtType == "*context.emptyCtx" {
		return
	}

	ictx := *(*iface)(unsafe.Pointer(&ctx))

	if ictx.data == 0 {
		return
	}

	valCtx := (*valueCtx)(unsafe.Pointer(ictx.data))
	if valCtx != nil && valCtx.key != nil && valCtx.value != nil {
		kv[valCtx.key] = valCtx.value
	}

	getKeyValues(valCtx.Context, kv)
}
