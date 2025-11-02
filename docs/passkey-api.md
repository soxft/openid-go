# Passkey 前端接入指南

本文档面向浏览器/前端开发者，说明如何串联后端 Passkey API 完成注册、登录与管理流程。

## 基础信息

- 所有接口前缀：`/passkey`
- 鉴权方式：需要现有账号已登录，并在请求头附带 `Authorization: Bearer <JWT>`
- 统一响应结构：

```json
{
  "success": true,
  "message": "success",
  "data": {
    "...": "..."
  }
}
```

当 `success=false` 时，`message` 字段直接给出错误提示，可据此进行前端展示。

## 常用工具函数

浏览器 WebAuthn API 会返回 `ArrayBuffer`，需要转换成 base64url 字符串后再发送给后端。可参考：

```ts
const toBase64Url = (buffer: ArrayBuffer): string => {
  const bytes = new Uint8Array(buffer);
  let binary = "";
  bytes.forEach(b => (binary += String.fromCharCode(b)));
  return btoa(binary).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/g, "");
};

const publicKeyCredentialToJSON = (cred: PublicKeyCredential) => ({
  id: cred.id,
  rawId: toBase64Url(cred.rawId),
  type: cred.type,
  clientExtensionResults: cred.getClientExtensionResults(),
  response: {
    attestationObject: cred.response.attestationObject
      ? toBase64Url((cred.response as AuthenticatorAttestationResponse).attestationObject)
      : undefined,
    clientDataJSON: toBase64Url(cred.response.clientDataJSON),
    authenticatorData: cred.response.authenticatorData
      ? toBase64Url((cred.response as AuthenticatorAssertionResponse).authenticatorData)
      : undefined,
    signature: cred.response.signature
      ? toBase64Url((cred.response as AuthenticatorAssertionResponse).signature)
      : undefined,
    userHandle: cred.response.userHandle
      ? toBase64Url((cred.response as AuthenticatorAssertionResponse).userHandle!)
      : undefined,
  },
});
```

后文所有 POST 请求体均指向此 JSON 结构。

## 注册流程

`navigator.credentials.create()` 生成的 `PublicKeyCredential` 结果需完整序列化后发送到后端。

- **获取注册参数**
  - `GET /passkey/register/options`
  - 成功 `data`：WebAuthn `PublicKeyCredentialCreationOptions`，直接作为浏览器注册调用的参数。

  ```ts
  const { data } = await fetchJson("/passkey/register/options");
  const credential = await navigator.credentials.create({ publicKey: data });
  ```

- **提交注册结果**
  - `POST /passkey/register`
  - `Content-Type: application/json`
  - 请求体：`PublicKeyCredential`（`navigator.credentials.create` 的完整 JSON，包含 `id`、`rawId`、`response` 等字段）。
  - 成功 `data`: `{ "passkeyId": number }`，表示新绑定的 Passkey 记录 ID。

  ```ts
  const credJSON = publicKeyCredentialToJSON(credential as PublicKeyCredential);
  const { data } = await fetchJson("/passkey/register", {
    method: "POST",
    body: JSON.stringify(credJSON),
  });
  console.log("new passkey id", data.passkeyId);
  ```

常见失败响应：
- `挑战已过期，请重试`：Redis 中的注册会话过期，需要重新获取 options。
- `注册失败`：注册数据校验错误或数据库写入失败。

## 登录流程

`navigator.credentials.get()` 生成的 `PublicKeyCredential` 结果需完整序列化后发送到后端。

- **获取登录参数**
  - `GET /passkey/login/options`
  - 成功 `data`：WebAuthn `PublicKeyCredentialRequestOptions`。
  - 若尚未绑定 Passkey，返回 `success=false`，`message=未绑定 Passkey`。

  ```ts
  const { data } = await fetchJson("/passkey/login/options");
  const credential = await navigator.credentials.get({ publicKey: data });
  ```

- **提交登录结果**
  - `POST /passkey/login`
  - `Content-Type: application/json`
  - 请求体：`PublicKeyCredential`（`navigator.credentials.get` 的完整 JSON）。
  - 成功 `data`: `{ "token": string, "passkeyId": number }`。
    - `token` 为新的 JWT，建议前端覆盖旧登录态。

  ```ts
  const credJSON = publicKeyCredentialToJSON(credential as PublicKeyCredential);
  const { data } = await fetchJson("/passkey/login", {
    method: "POST",
    body: JSON.stringify(credJSON),
  });
  updateToken(data.token);
  ```

常见失败响应：
- `挑战已过期，请重试`：Redis 中的登录会话失效，需要重新获取 options。
- `未绑定 Passkey`：用户无可用凭证。
- `登录失败`：签名校验失败或服务器异常。

## Passkey 管理

- **列表 Passkey**
  - `GET /passkey`
  - 成功 `data`: `{ "items": PasskeySummary[] }`
  - `PasskeySummary` 字段：
    - `id`: 记录 ID
    - `createdAt`: 创建时间
    - `lastUsedAt`: 最近使用时间（可能为 `null`）
    - `cloneWarning`: WebAuthn Clone Warning 标记
    - `signCount`: 签名计数
    - `transports`: 可用传输方式字符串数组

- **删除 Passkey**
  - `DELETE /passkey/:id`
  - 路径参数 `id`：待删除记录的整数 ID。
  - 成功 `message`: `success`。
  - 若 ID 不存在，返回 `success=false`，`message=Passkey 不存在`。

## 前端调用提示

- 建议直接使用 WebAuthn API 的 `publicKey` 参数与响应，避免二次加工导致字段缺失。
- `rawId`、`response.attestationObject`、`response.clientDataJSON` 等字段需要 base64url 编码；浏览器返回的 `ArrayBuffer` 需自行转换。
- 操作超时或页面刷新会导致会话丢失，需重新调用 `GET /passkey/.../options`。
- 若遇到 `未绑定 Passkey`，可提示用户先走注册流程再重试登录。
- 在同一页面内重复使用 `navigator.credentials.*` 前建议捕获 `AbortError`、`NotAllowedError` 并给出友好提示。
