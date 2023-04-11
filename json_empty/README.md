json.DisallowUnknownFields 用途
===

## Example

运行 main.go 会产生以下错误

```shell
{0 []} json: unknown field "name"
```

- `json.DisallowUnknownFields` 开关开启后，json 字节流中如果存在未知的字段，那么 `json.Decode` 会报错
- 但是，当 struct 中定义的 field 在 json 字节流中不存在时，`json.Decode` 和 `json.Unmarshal` 都不会有任何错误。
