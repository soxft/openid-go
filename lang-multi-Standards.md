# 多语言开发规范

## 概述

本文档面向 LLM 代码代理与团队成员，提供一套覆盖 Go、TypeScript、Python、Shell 以及配置文件的统一开发标准。目标是在多语言协作场景下保持一致的编码风格、可靠的工程质量与高效的交付节奏。所有规范基于团队既有实践，并结合模块化架构、防御性编程、TDD 驱动与高覆盖率测试的要求制定。

---

## 通用开发原则

### 项目结构与模块化

- 采用领域驱动的目录结构，以清晰职责边界划分模块，例如 `api/`、`app/`、`library/`、`process/`、`config/`。
- 新增语言组件时，在根目录创建语义明确的顶层目录（例如 `frontend/`, `scripts/`）。
- 共享逻辑封装为独立的内部包，禁止跨层直接依赖底层细节。
- 所有服务入口放置在 `cmd/<service-name>/main.(go|ts|py)` 内，避免根目录臃肿。

### 版本控制与分支策略

- 采用 Git Flow 派生模型：`main`（生产）、`develop`（集成）、`feature/*`、`hotfix/*`、`release/*`。
- 提交信息遵循 Conventional Commits，结合语义化版本（SemVer）管理发布。
- PR 命名建议：`[type] scope - summary`，并关联追踪任务或 issue。
- 统一使用 `git rebase` 同步主干，避免无意义的 merge commit。

### 文档与知识沉淀

- 所有新模块必须更新 `README.md` 或 `docs/` 下相应章节，记录架构、依赖、运行方式。
- 复杂设计新增简短 ADR（Architecture Decision Record）。
- 公有 API、脚本和 CLI 工具应提供 `--help` 或 README 文档。
- 在代码中使用结构化注释（如 GoDoc、JSDoc、Google-Style Docstring）。

### 测试与质量门禁

- TDD 优先：先写失败测试再实现功能。
- 要求核心业务模块覆盖率 ≥ 80%，工具代码 ≥ 60%。
- CI 中运行静态检查、单元测试、集成测试，必要时执行安全扫描（如 `gosec`, `npm audit`, `pip-audit`）。
- 引入性能敏感改动需提供 Benchmark/Profiling 数据。

### 代码审查流程

- PR 至少由一名非作者 Reviewer 审核，跨语言改动需对应语言 Reviewer 参与。
- Reviewer 关注点：可读性、一致性、防御性、边界处理、回归风险、测试完备性、性能与安全。
- 代码审查提出的 TODO 必须转化为 issue 或后续任务，避免 PR 长期挂起。

### 安全与隐私

- 所有秘钥与凭证存储在 `config.tpl` 或环境变量模板中，禁入库。
- 默认启用输入校验、输出编码、防重放、防注入、防 XXE。
- 对外接口添加速率限制与审计日志。
- 日志中屏蔽敏感字段（密码、Token、身份标识）。

---

## 语言特定规范

### Go

- **文件组织**: 每个包仅暴露最小 API；测试文件与被测文件同目录，命名 `*_test.go`；命令行入口置于 `cmd/<service>/main.go`。
- **命名约定**: 包名小写且具描述性；导出符号使用 PascalCase，错误变量 `ErrXxx`；接口偏向动名词结尾，如 `TokenReader`。
- **代码风格**: 强制 `goimports`; 使用卫语句减少嵌套；避免全局状态，若必须则以 `sync.Once` 控制初始化。
- **注释规范**: 导出符号使用 GoDoc，HTTP 处理函数补充路由与权限说明；复杂逻辑可使用 `// NOTE:` 对代理提示注意点。
- **导入与模块**: 导入分组（标准库 / 第三方 / 内部）；模块依赖由 `go.mod` 管理，更新依赖需执行 `go mod tidy`。
- **错误处理**: 统一返回 `error`，禁止忽略；使用 `errors.Join` 聚合，或自定义领域错误类型；HTTP 层返回结构化 JSON 错误。
- **测试规范**: 使用 `testing`, `testify`；基准测试放在 `benchmark_test.go`；集成测试依赖 Docker Compose 时标注 `// +build integration`。
- **性能考虑**: 避免在热路径中使用反射；使用上下文超时；利用 `sync.Pool` 复用大对象。

### TypeScript / JavaScript

- **文件组织**: 根据层级划分 `src/domain`, `src/application`, `src/infrastructure`; UI 组件分 `components`, `hooks`, `pages`。
- **命名约定**: 文件使用 `kebab-case.tsx`; 类、React 组件用 PascalCase；常量全大写加下划线。
- **代码风格**: 使用 `2` 空格缩进；强制 `strict` 模式；拒绝隐式 `any`；Prefer `const`。
- **注释规范**: 使用 JSDoc，关键函数记录参数、返回值和副作用；TODO 需附 issue，如 `// TODO(#456): refactor state machine`。
- **导入管理**: 使用路径别名（配置 `tsconfig.json`）限定跨域导入；禁止循环依赖。
- **错误处理**: 使用 `Result` 类型或 `try/catch` 包装；Promise 链必须 `catch`；API 请求统一封装错误码到业务异常。
- **测试规范**: 使用 `vitest`/`jest`; 文件命名 `*.spec.ts`; 组件测试使用 `testing-library`；集成测试跑在 `playwright`。
- **性能考虑**: 避免在渲染函数创建匿名闭包；使用 Suspense/Lazy 加载大组件；服务端启用缓存与 CDN。

### Python

- **文件组织**: 脚本放 `scripts/`; 应用结构 `app/` (package) + `tests/`; CLI 入口在 `app/__main__.py`。
- **命名约定**: 模块与包使用 `snake_case`; 类用 PascalCase; 函数、变量采用 `snake_case`; 常量全大写。
- **代码风格**: 使用 `black`（行宽 100）+ `isort`；启用 `mypy` 做静态检查；避免可变默认参数。
- **注释规范**: Google 风格 Docstring；必要复杂逻辑使用 `# NOTE:`；TODO 带责任人或 issue。
- **依赖管理**: 使用 `poetry` 或 `pip-tools`; 锁定版本并开启 Hash 校验；虚拟环境统一 `./.venv`。
- **错误处理**: 捕获指定异常；日志中包含 `exc_info=True`; 自定义异常继承自 `Exception` 并提供语义。
- **测试规范**: 使用 `pytest`; 测试文件 `tests/test_<module>.py`; fixture 放在 `conftest.py`；集成测试打 `@pytest.mark.integration`。
- **性能考虑**: 使用生成器减少内存；必要时启用 `multiprocessing` 或 `asyncio`; 提供 `profiling/` 脚本。

### Shell (Bash/Zsh)

- **文件组织**: 脚本放在 `scripts/` 或 `ops/`; 可执行文件加执行权限并以 `#!/usr/bin/env bash` 开头。
- **命名约定**: 文件 `kebab-case.sh`; 变量大写；函数 `snake_case`。
- **代码风格**: 使用 `set -euo pipefail`; 依赖 `shellcheck`；复杂逻辑拆分函数。
- **注释规范**: 顶部描述脚本用途与参数；关键步骤前添加注释说明副作用。
- **依赖管理**: 所需工具在 README 标明；必要时在脚本中做版本检测。
- **错误处理**: 使用 `trap` 捕捉错误；对外暴露 EXIT 码；输出日志到 stderr。
- **测试规范**: 使用 `bats` 编写自动化测试。
- **性能考虑**: 避免在循环中调用外部命令；优先使用内建。

### 配置文件 (YAML/JSON/Docker)

- **组织**: 环境配置按 `config/<env>/` 分类；共享模板 `.tpl` 置于 `config/templates/`。
- **命名**: 使用环境后缀，例如 `app.production.yaml`；敏感字段用占位符标记。
- **风格**: YAML 两空格缩进；所有键小写带连字符；JSON 使用 `jq` 格式化。
- **注释**: YAML 通过 `#` 说明目的与默认值；JSON 使用旁注文档。
- **依赖管理**: Docker 镜像固定 tag；Compose 文件遵循 v3 以上。
- **错误处理**: 配置解析失败需在应用启动阶段终止，记录错误。
- **测试**: 使用 `kubeval`/`yamllint` 验证；CI 中对 Helm Chart 执行 `helm template --strict`。
- **性能/运维**: 指定资源限制与健康检查；记录日志收集方式。

---

## 工具配置

### VS Code 推荐设置 (`.vscode/settings.json`)

```json
{
  "editor.formatOnSave": true,
  "editor.rulers": [100],
  "files.trimTrailingWhitespace": true,
  "files.insertFinalNewline": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go",
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  },
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[typescriptreact]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[python]": {
    "editor.defaultFormatter": "ms-python.black-formatter"
  },
  "[shellscript]": {
    "editor.defaultFormatter": "foxundermoon.shell-format"
  },
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "python.formatting.provider": "none",
  "typescript.tsserver.experimental.enableProjectDiagnostics": true
}
```

### 通用 `.editorconfig`

```ini
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = space
indent_size = 2

[*.go]
indent_style = tab
indent_size = 4

[*.py]
indent_size = 4

[*.sh]
indent_size = 2

[Makefile]
indent_style = tab
```

### 格式化与 Lint 配置

- **Go**: `gofmt`, `goimports`, `golangci-lint` (`.golangci.yml` 见现有规范)。
- **TypeScript**: `eslint` + `prettier` 结合；`eslint.config.js` 示例：

```js
import eslint from "@eslint/js";
import tsParser from "@typescript-eslint/parser";
import tsPlugin from "@typescript-eslint/eslint-plugin";
import prettierPlugin from "eslint-plugin-prettier";

export default [
  {
    ignores: ["dist/**", "node_modules/**"],
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        project: "./tsconfig.json"
      }
    },
    plugins: {
      "@typescript-eslint": tsPlugin,
      prettier: prettierPlugin
    },
    rules: {
      ...eslint.configs.recommended.rules,
      ...tsPlugin.configs.recommended.rules,
      "prettier/prettier": "error",
      "@typescript-eslint/no-explicit-any": "error",
      "@typescript-eslint/consistent-type-imports": "warn"
    }
  }
];
```

- **Python**: `pyproject.toml` 片段：

```toml
[tool.black]
line-length = 100
target-version = ["py311"]

[tool.isort]
profile = "black"

[tool.mypy]
python_version = "3.11"
strict = true
warn_unused_configs = true

[tool.pytest.ini_options]
minversion = "7.0"
addopts = "-ra -q"
testpaths = ["tests"]
```

- **Shell**: 在 `package.json` 或 `Makefile` 中加入 `shellcheck scripts/*.sh` 任务。

### 预提交钩子 (`.pre-commit-config.yaml`)

```yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: check-yaml
      - id: check-case-conflict
      - id: end-of-file-fixer
      - id: trailing-whitespace
  - repo: https://github.com/golangci/golangci-lint
    rev: v1.60.1
    hooks:
      - id: golangci-lint
  - repo: https://github.com/pre-commit/mirrors-eslint
    rev: v8.57.0
    hooks:
      - id: eslint
        additional_dependencies:
          - eslint
          - prettier
          - typescript
          - @typescript-eslint/parser
          - @typescript-eslint/eslint-plugin
  - repo: https://github.com/psf/black
    rev: 24.10.0
    hooks:
      - id: black
  - repo: https://github.com/pycqa/isort
    rev: 5.13.2
    hooks:
      - id: isort
  - repo: https://github.com/igorshubovych/markdownlint-cli
    rev: v0.42.0
    hooks:
      - id: markdownlint
```

### CI/CD 建议

- 使用 GitHub Actions：

```yaml
name: ci

on:
  pull_request:
    branches: ["main", "develop"]
  push:
    branches: ["main"]

jobs:
  lint-test:
    runs-on: ubuntu-latest
    services:
      redis:
        image: redis:7-alpine
        ports: ["6379:6379"]
      db:
        image: postgres:15-alpine
        env:
          POSTGRES_PASSWORD: example
        ports: ["5432:5432"]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - uses: actions/setup-node@v4
        with:
          node-version: "20"
      - uses: actions/setup-python@v5
        with:
          python-version: "3.11"
      - run: go test ./...
      - run: pnpm install && pnpm test --if-present
      - run: pip install -r requirements.txt && pytest
      - run: make lint
```

- 发布流水线需包含镜像构建、扫描（Trivy）、签名（cosign）与部署。

---

## 代码审查标准

- **一致性**: 是否遵守本规范及现有代码风格。
- **鲁棒性**: 边界条件、异常路径、超时、重试策略是否完善。
- **安全性**: 鉴权、授权、输入验证、敏感信息处理是否妥当。
- **可维护性**: 模块间依赖是否合理，是否存在隐性耦合。
- **测试充分性**: 单元、集成、端到端测试是否覆盖关键逻辑与回归场景。
- **性能影响**: 新增算法复杂度、资源占用是否符合预期，并提供度量。
- **可观测性**: 日志、指标、追踪埋点是否齐全，告警阈值是否更新。

---

## 最佳实践示例

### Go

**推荐写法**

```go
// IssueToken issues a signed JWT for the given account ID.
func IssueToken(ctx context.Context, accountID int64, signer jwt.Signer) (string, error) {
    if accountID <= 0 {
        return "", fmt.Errorf("issue token: invalid account id %d", accountID)
    }

    token, err := signer.Sign(jwt.MapClaims{
        "sub": accountID,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
    })
    if err != nil {
        return "", fmt.Errorf("issue token: sign: %w", err)
    }

    return token, nil
}
```

**需避免**

```go
func issueToken(id int64, signer jwt.Signer) string {
    token, _ := signer.Sign(jwt.MapClaims{"sub": id})
    return token
}
```

### TypeScript

**推荐写法**

```ts
export interface PasskeyChallenge {
  challenge: string;
  expiresAt: number;
}

export async function createChallenge(userId: string, repo: ChallengeRepository): Promise<PasskeyChallenge> {
  const challenge = await repo.create(userId);
  if (!challenge) {
    throw new AppError("challenge:create", "failed to create challenge");
  }
  return challenge;
}
```

**需避免**

```ts
export async function createChallenge(userId) {
  return await repo.create(userId);
}
```

### Python

**推荐写法**

```python
@dataclass(slots=True)
class PasskeyChallenge:
    challenge: str
    expires_at: datetime


def issue_challenge(user_id: str, repo: ChallengeRepo) -> PasskeyChallenge:
    if not user_id:
        raise ValueError("user_id is required")
    challenge = repo.create(user_id)
    if challenge is None:
        msg = "failed to persist challenge"
        raise RepositoryError(msg)
    return challenge
```

**需避免**

```python
def issue_challenge(user_id, repo):
    return repo.create(user_id)
```

### Shell

**推荐写法**

```bash
#!/usr/bin/env bash
set -euo pipefail

log() {
  printf '%%s\n' "$*" >&2
}

if [[ $# -lt 1 ]]; then
  log "usage: $0 <env>"
  exit 1
fi

ENV="$1"
docker compose -f "deploy/${ENV}.yaml" up -d
```

**需避免**

```bash
#!/bin/bash
ENV=$1
docker-compose up -d
```

---

## 升级与扩展指南

- 添加新语言前需扩展本规范并在 `AGENTS.md` 登记。
- 升级依赖需通过自动化工具（`renovate`, `dependabot`）追踪，并在 Changelog 中记录破坏性变更。
- 关键组件升级前应创建 PoC 分支验证兼容性，必要时编写迁移脚本。

## 常见问题解答 (FAQ)

- **如何快速启动多语言开发环境？**
  - 使用 `make setup` 安装 Go/Node/Python 依赖；必要服务通过 `docker compose up` 启动。
- **LLM 如何选择合适规范？**
  - 根据文件扩展名匹配语言章节，并查阅通用原则获取共享约束。
- **测试过慢怎么办？**
  - 拆分测试套件并在 CI 中并行运行；对慢测试打标签按需执行。
- **如何添加新的预提交检查？**
  - 修改 `.pre-commit-config.yaml` 并运行 `pre-commit install`; 在 PR 描述中说明新增条目。

---

本规范为活文档，建议至少每季度回顾一次，并在重要变更后同步更新相关自动化配置与培训资料。
