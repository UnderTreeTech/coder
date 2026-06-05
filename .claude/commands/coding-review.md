# 代码审查门禁执行命令

## 描述
执行编码阶段完成后的门禁检查（即代码审查）。

## 触发条件
用户输入 `/coding:review`。

## 执行动作
1. 收集当前工作区未提交的代码变更（diff）。
2. 调用 `code-review-report` Skill 或相关审查 Agent。
3. 校验代码是否符合团队规范、业务逻辑是否符合设计文档、核心逻辑是否有单测覆盖。
4. **IDL 定义与生成**：IDL定义在 `api/idl/{service_name}/`，务必校验IDL是否有正确生成，如未生成报错并反馈需执行 `waterdrop protoc --grpc --swagger *.proto` 来生成。
5. 在api目录下执行 `go fmt ./...` 及 `go build ./...` 做格式化统一及验证编译是否能通过（api/idl目录除外）。
6. 生成审查报告并写入 `requirements/{requirement-id}/gate-3-code-review.md`。
7. 若通过，则提示用户可以提交代码；若不通过，则指出需要修改的地方。