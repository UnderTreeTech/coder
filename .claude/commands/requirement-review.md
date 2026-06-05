# 需求评审门禁执行命令

## 描述
执行需求阶段门禁检查。

## 触发条件
用户输入 `/requirement:review`。

## 执行动作
1. 加载当前活动的 requirement id（可从 `.harness/local.yaml` 或询问用户获得）。
2. 调用 `requirement-quality-reviewer` Agent。
3. 读取对应的 `requirement.md` 文件。
4. 校验是否满足需求评审门禁所有条件。
5. 生成报告并写入 `requirements/{requirement-id}/gate-1-requirement-review.md`。
6. 输出结果给用户。