# reconnect-resume

断线重连的续局会话服务，用 `net/http` 暴露两个接口。只用 Go 标准库。

```
go test ./...
go run ./cmd/serve --sample=grace
go run ./cmd/serve --sample=version
go run ./cmd/serve --sample=idempotent
go run ./cmd/serve --sample=work
```

## 接口

- `POST /session/open` 请求体 `{"player": "...", "atMs": N}`，返回会话快照。
- `POST /session/resume` 请求体 `{"requestId": "...", "token": "...", "version": N, "atMs": N}`。

## 口径

- **宽限窗口**：令牌最后一次活动后 `GraceMs` 毫秒内可以续局成功；超过窗口返回 `410` 与 `{"reason":"expired"}`。
- **版本校验**：请求带的 `version` 必须等于服务端当前版本，否则返回 `409` 与 `{"reason":"stale"}`。
- **幂等**：同一个 `requestId` 只生效一次；重复出现返回第一次的快照，不推进版本、不计入 `Resumes`。
- **未知令牌**：返回 `409` 与 `{"reason":"unknown"}`。
- **索引**：令牌必须按索引查找，`Scanned` 不得随会话总数放大。

## 不变量

- 一次成功续局只把版本推进一步。
- 被拒绝的请求不计入 `Resumes`。
- 同一 `requestId` 重复提交得到的快照逐字段相同。

## 输出契约

`--sample` 模式每行打印 `状态码 响应体`；
`--sample=work` 只打印 `{"resumes": N, "scanned": N}`。
