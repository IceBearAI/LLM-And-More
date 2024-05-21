# openapi 服务

openapi 服务是兼容openai的sdk，可以通过openapi服务调用openai的接口。

目前已支持接口:

- [x] `/v1/chat/completions` 对话接口
- [x] `/v1/completions` 文本接口
- [x] `/v1/models` 图片接口
- [ ] `/v1/embeddings` embedding接口

## 使用方法

```bash
$ curl -X POST "http://localhost:8000/v1/chat/completions" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $API_KEY" \
    -d '{
        "model": "chatglm3-6b-32k",
        "messages": [{"role": "system", "content": "You are a helpful assistant."}, {"role": "user", "content": "Hello!"}]
    }'
```