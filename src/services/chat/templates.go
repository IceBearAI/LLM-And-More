package chat

import "context"

type Templates interface {
	Register(ctx context.Context, name string, conv Conv)
	GetConv(ctx context.Context, name string) (Conv, bool)
}

type templates struct {
	conv map[string]Conv
}

func (t *templates) Register(ctx context.Context, name string, conv Conv) {
	if _, ok := t.conv[name]; ok {
		return
	}
	t.conv[name] = conv
}

func (t *templates) GetConv(ctx context.Context, name string) (Conv, bool) {
	conv, ok := t.conv[name]
	return conv, ok
}

func NewTemplates() Templates {
	return &templates{conv: map[string]Conv{}}
}
