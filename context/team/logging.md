# 日志规范

## 1. 日志级别定义

| 级别 | 使用场景 | 示例 |
|-----|---------|------|
| `DEBUG` | 本地开发调试，**生产禁用** | 变量中间值、SQL 语句 |
| `INFO` | 关键业务路径、请求入口/出口 | 请求开始、接口调用成功 |
| `WARN` | 可恢复的异常、降级情况 | 缓存未命中、重试触发 |
| `ERROR` | 需要人工介入的错误 | 数据库连接失败、第三方超时 |
| `FATAL` | 系统无法继续运行 | 配置加载失败、端口冲突 |

---

## 2. 日志必填字段

每条日志**必须**包含以下字段：

```json
{
    "level": "info",
    "timestamp": "2026-06-04T10:00:00.000+08:00",
    "service": "{service-name}",
    "trace_id": "xxx-yyy-zzz",
    "request_id": "req-xxx-yyy",
    "msg": "日志内容"
}
```

| 字段 | 类型 | 必填 | 说明 |
|-----|------|------|-----|
| `level` | string | ✅ | 日志级别（小写） |
| `timestamp` | string | ✅ | RFC3339 格式，包含时区 |
| `service` | string | ✅ | 服务名（来自配置） |
| `trace_id` | string | ✅ | 链路追踪 ID，跨服务传递 |
| `request_id` | string | ✅ | 单次请求 ID |
| `msg` | string | ✅ | 日志内容（中文或英文，不混用） |
| `user_id` | string | ❌ | 涉及用户操作时必填，**禁止记录明文密码/手机号** |
| `error` | string | ❌ | ERROR 级别时必填，记录错误信息 |
| `duration_ms` | int | ❌ | 接口/操作耗时（毫秒） |

---

## 3. 禁止记录的字段

以下信息**禁止**出现在日志中：

- 用户密码（明文或密文均不可）
- 手机号（脱敏后可，格式：`138****8888`）
- 身份证号
- 银行卡号
- 支付密码
- Token 原文（可记录前 8 位）

---

## 4. Go 日志使用示例

```go
import "github.com/UnderTreeTech/waterdrop/pkg/log"

// 正确：结构化日志
log.Debug(
	ctx.Request.Context(), 
	"reply", 
	log.Int("code", estatus.Code()),
	log.String("msg", reply.Message),
	log.Any("data", reply.Datas), 
	)

// 错误：使用 fmt.Println（生产禁用）
// fmt.Println("用户登录:", userID)

// 错误：记录敏感信息
// logger.Info("登录请求", zap.String("password", req.Password))
```

---

## 5. 日志采集规范

- 日志输出到 **stdout**（容器化部署），不直接写文件
- 本地开发可输出到文件，路径：`logs/{service-name}.log`（已加入 `.gitignore`）
- 日志文件大小超过 **100MB** 自动轮转，保留最近 **7** 天

---

## 6. 性能约束

- 热路径（QPS > 1000）的 DEBUG 日志必须使用条件判断
- 单条日志内容不超过 **4KB**
- 禁止在循环内打印日志（除非有采样）
