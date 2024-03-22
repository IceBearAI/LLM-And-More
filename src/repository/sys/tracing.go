package sys

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (t *tracing) ListFineTuningTemplate(ctx context.Context, page, pageSize int, name, templateType string) (res []types.FineTuningTemplate, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("page", page, "pageSize", pageSize, "name", name, "templateType", templateType, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.ListFineTuningTemplate(ctx, page, pageSize, name, templateType)
}

func (t *tracing) SaveFineTuningTemplate(ctx context.Context, data *types.FineTuningTemplate) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.SaveFineTuningTemplate(ctx, data)
}

func (t *tracing) DeleteFineTuningTemplate(ctx context.Context, name string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("name", name, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.DeleteFineTuningTemplate(ctx, name)
}

func (t *tracing) IsExistTuningTemplate(ctx context.Context, name string) (isExist bool, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("name", name, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.IsExistTuningTemplate(ctx, name)
}

func (t *tracing) GetFineTuningTemplate(ctx context.Context, name string) (res types.FineTuningTemplate, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("name", name, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.GetFineTuningTemplate(ctx, name)
}

func (t *tracing) CreateAudit(ctx context.Context, data *types.SysAudit) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CreateAudit", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.CreateAudit(ctx, data)
}

func (t *tracing) ListAudit(ctx context.Context, request ListAuditRequest) (items []types.SysAudit, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListAudit", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("request", request, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.ListAudit(ctx, request)
}

func (t *tracing) ListDict(ctx context.Context, request ListDictRequest) (items []types.SysDict, total int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListDict", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("request", request, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.ListDict(ctx, request)
}

func (t *tracing) CreateDict(ctx context.Context, data *types.SysDict) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CreateDict", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("data", data, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.CreateDict(ctx, data)
}

func (t *tracing) GetDict(ctx context.Context, id uint) (res types.SysDict, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetDict", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.GetDict(ctx, id)
}

func (t *tracing) UpdateDict(ctx context.Context, data *types.SysDict, updateChildrenCode bool) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "UpdateDict", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("data", data, "err", err, "updateChildrenCode", updateChildrenCode)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.UpdateDict(ctx, data, updateChildrenCode)
}

func (t *tracing) DeleteDict(ctx context.Context, id uint) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DeleteDict", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("id", id, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.DeleteDict(ctx, id)
}

func (t *tracing) GetDictByCode(ctx context.Context, code string) (res types.SysDict, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetDictByCode", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("code", code, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.GetDictByCode(ctx, code)
}

func (t *tracing) FindDictTreeByParentId(ctx context.Context, parentId uint, parentDictType ...string) (res []types.SysDict, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindDictTreeByParentId", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("parentId", parentId, "err", err, "parentDictType", parentDictType)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindDictTreeByParentId(ctx, parentId, parentDictType...)
}

func (t *tracing) GetDictByDictValue(ctx context.Context, parentId uint, dictValue string) (res types.SysDict, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetDictByDictValue", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("code", parentId, "dictValue", dictValue, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.GetDictByDictValue(ctx, parentId, dictValue)
}

func (t *tracing) FindDictTreeByCode(ctx context.Context, code []string) (res []types.SysDict, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "FindDictTreeByCode", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "repository.sys",
	})
	defer func() {
		span.LogKV("code", code, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.FindDictTreeByCode(ctx, code)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}
