# Go 语言开发规范

## 概述

本文档旨在为基于 Go 语言的项目提供一套统一的开发规范，以确保代码的一致性、可读性和可维护性。本文档的目标读者是 LLM 代码代理和所有项目参与者。规范内容基于 Go 社区的最佳实践，并结合了本项目的特定风格。

## 1. 通用开发原则

### 1.1. 项目结构标准

项目应遵循分层架构，将不同职责的代码分离到独立的目录中。推荐使用以下结构：

```
.
├── api/         # 外部 API 接口定义 (例如 OpenAPI/Swagger)
├── app/         # 应用核心逻辑
│   ├── controller/ # 控制器/处理器层，处理 HTTP 请求
│   ├── middleware/ # 中间件
│   └── model/      # 数据模型 (GORM structs)
├── config/      # 配置加载与管理
├── core/        # 核心引导程序
├── docs/        # 项目文档
├── library/     # 通用库/工具函数
├── process/     # 后台进程、数据库、缓存等初始化
└── main.go      # 程序入口
```

### 1.2. 版本控制规范

- **分支模型**: 推荐使用 `Git Flow` 或类似的策略。
  - `main`: 稳定的主分支，用于生产发布。
  - `develop`: 开发分支，集成所有已完成的功能。
  - `feature/xxx`: 功能开发分支。
  - `hotfix/xxx`: 紧急修复分支。
- **提交信息**: 遵循 `Conventional Commits` 规范。
  - 格式: `<type>(<scope>): <subject>`
  - 示例: `feat(passkey): add support for discoverable login`
  - 常用 `type`: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`。

### 1.3. 文档编写要求

- 所有公开的函数、类型和常量都应有文档注释。
- API 端点应在对应的 `controller` 函数注释中清晰标明，包含请求方法和路径。
- 复杂的业务逻辑应在代码中添加必要的行内注释。

### 1.4. 代码审查流程

- 所有代码变更必须通过 Pull Request (PR) 提交。
- PR 必须至少有一位其他成员审查通过后才能合并。
- 审查重点：代码风格、逻辑正确性、测试覆盖率、文档完整性。

## 2. Go 语言特定规范

### 2.1. 文件组织

- **文件命名**: 使用小写下划线 `snake_case.go`。
- **目录结构**:
  - 按功能或领域划分包（`package`）。
  - 包名应简短、清晰，并使用小写。
  - 避免使用 `util` 或 `common` 等模糊的包名，优先使用功能性命名，如 `apiutil`, `userutil`。

### 2.2. 命名约定

- **变量**: 使用 `camelCase`。
- **函数/方法**: 公开函数使用 `PascalCase`，私有函数使用 `camelCase`。
- **接口**: 接口名以 `er` 结尾，如 `Reader`, `Writer`。
- **常量**: 使用 `PascalCase`。
- **包名**: 小写，简明扼要。
- **错误变量**: 以 `Err` 开头，如 `ErrSessionNotFound`。

### 2.3. 代码风格

- **格式化**: 所有代码必须使用 `gofmt` 或 `goimports` 进行格式化。
- **缩进**: 使用制表符 `\t`。
- **行长度**: 建议不超过 120 个字符。
- **导入/包管理**:
  - 使用 `goimports` 自动管理导入顺序。
  - 导入顺序：标准库、第三方库、项目内库，组间用空行分隔。
  - 依赖管理使用 `Go Modules` (`go.mod`, `go.sum`)。

### 2.4. 注释规范

- **文档字符串**:
  - 为所有公开的函数、类型、常量和变量提供注释。
  - 注释应以被注释对象的名字开头。
  - 示例:
    ```go
    // PasskeyRegistrationOptions 获取 Passkey 注册选项
    //
    //	GET /passkey/register/options
    func PasskeyRegistrationOptions(c *gin.Context) { ... }
    ```
- **行内注释**: 用于解释复杂或不直观的代码逻辑。
- **TODO 注释**:
  - 使用 `// TODO:` 格式标记待办事项。
  - 最好能附上相关 issue 编号或责任人。
  - 示例: `// TODO(#123): Refactor this to improve performance.`

### 2.5. 错误处理

- **总是检查错误**: 不要忽略函数返回的 `error`，除非明确知道可以安全地忽略。
- **错误信息**:
  - 错误信息应为小写，不以标点符号结尾。
  - 使用 `fmt.Errorf` 添加上下文信息: `fmt.Errorf("passkey begin registration: %w", err)`。
  - 使用 `%w` 来包装底层错误，以便上层可以使用 `errors.Is` 或 `errors.As`。
- **错误返回**:
  - 在函数或方法的开头处理 "happy path" 之前的卫语句（guard clauses）。
  - 错误应尽早返回。
  - 示例:
    ```go
    account, err := getAccount(c)
    if err != nil {
        api.Fail("user not found")
        return
    }
    ```

### 2.6. 测试规范

- **测试文件**: 测试文件命名为 `_test.go`。
- **测试函数**: 测试函数以 `Test` 开头，例如 `func TestMyFunction(t *testing.T)`。
- **测试覆盖率**: 核心业务逻辑的测试覆盖率应达到 80% 以上。
- **测试驱动开发 (TDD)**: 鼓励在新功能开发前先编写测试用例。
- **Mocking**: 使用 `gomock` 或 `testify/mock` 等库来模拟依赖。

## 3. 工具配置

### 3.1. 编辑器配置 (VS Code)

在项目根目录创建 `.vscode/settings.json`:

```json
{
  "go.formatTool": "goimports",
  "go.useLanguageServer": true,
  "editor.formatOnSave": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  },
  "gopls": {
    "ui.semanticTokens": true
  }
}
```

### 3.2. Lint 规则配置

使用 `golangci-lint` 作为代码检查工具。在项目根目录创建 `.golangci.yml`:

```yaml
run:
  timeout: 5m
  skip-dirs:
    - vendor/

linters:
  enable:
    - gofmt
    - goimports
    - revive
    - govet
    - staticcheck
    - unused
    - errcheck

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck # Don't require error checking in tests
```

### 3.3. 预提交钩子配置

使用 `pre-commit` 框架来自动化代码检查。在项目根目录创建 `.pre-commit-config.yaml`:

```yaml
repos:
-   repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.3.0
    hooks:
    -   id: check-yaml
    -   id: end-of-file-fixer
    -   id: trailing-whitespace
-   repo: https://github.com/golangci/golangci-lint
    rev: v1.50.1
    hooks:
    -   id: golangci-lint
```

## 4. 代码审查标准

- **可读性**: 代码是否清晰易懂？
- **一致性**: 是否遵循了本规范？
- **简洁性**: 是否有不必要的复杂性？
- **正确性**: 代码是否能正确实现需求？
- **测试**: 是否有足够的测试覆盖？
- **文档**: 注释和文档是否完整？

## 5. 最佳实践示例

### 5.1. 良好的代码示例

**清晰的函数定义和错误处理**

```go
// findUserByID finds a user by their ID.
// It returns the user and an error if one occurred.
func findUserByID(ctx context.Context, id int) (*model.User, error) {
    if id <= 0 {
        return nil, errors.New("invalid user id")
    }

    var user model.User
    if err := db.WithContext(ctx).First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("user with id %d not found", id)
        }
        return nil, fmt.Errorf("database error: %w", err)
    }

    return &user, nil
}
```

### 5.2. 应避免的代码示例

**忽略错误和模糊的命名**

```go
// Bad example
func GetUser(id int) *model.User {
    user := model.User{}
    // Error is ignored!
    db.First(&user, id)
    return &user
}
```

**嵌套过深**

```go
// Bad example
func processRequest(c *gin.Context) {
    if c.Request.Method == "POST" {
        err := c.Request.ParseForm()
        if err == nil {
            // ... more nested logic
        } else {
            // handle error
        }
    } else {
        // handle other methods
    }
}
```

---
这份规范旨在成为一个动态文档，随着项目的发展和团队的成长而不断完善。