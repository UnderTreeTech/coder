# 新建设计方案命令

## 描述
基于已通过门禁的需求文档，开始进行方案设计和任务拆解。

## 触发条件
用户输入 `/design:new`。

## 执行动作
1. 确认当前需求已通过需求门禁（检查 `gate-1-requirement-review.md`）。
2. **【交互确认层】强制阻断：**
   - 先读取 `.service-matrix/dependencies.yaml` 分析可能的依赖关系。
   - **必须明确询问用户**：“本次功能实现落在 `api/` 下的哪些具体服务中？这些接口是作为内部普通 Service 方法、还是对外暴露的 HTTP，或是跨服务的 gRPC？”
   - **等待用户回答后**，才允许进行下一步。
3. 基于用户确认的服务归属和通信协议，在 `requirements/{requirement-id}/` 目录下创建 `design.md`。
   - **注意**：设计文档 (`design.md`) 必须包含以下核心章节：本次需求的改动服务范围、整体流程图或流程说明、架构方案（包含改动代码结构）、可能带来的性能瓶颈点分析、潜在的风险点及应对措施。
   - **注意**：如果是普通 service 方法，则设计中不得包含 gRPC/proto 的相关章节。
   - **注意**：如果是 gRPC 或基于 IDL 的 HTTP 接口，必须在设计中涵盖 IDL 设计部分：
     - 到对应的服务 IDL 目录（`api/idl/{service_name}/`）找相应的 proto 文件。
     - 如果存在该文件，则在此基础上按团队的 protobuf 规范（`context/team/protobuf-style-guide.md`）定义新接口。
     - 如果不存在该文件，则需按照 protobuf 规范创建新的 proto 文件，再定义接口。
   - **注意**：如果是表结构设计，必须遵循 db设计规范 （`context/team/db_design.md`）。所有表设计除主键自带索引外，仅在设计文档中建议可以创建的索引，不要直接生成带索引的DDL。
4. 基于 `design.md` 引导生成 `tasks/features.json` 进行任务拆解。
