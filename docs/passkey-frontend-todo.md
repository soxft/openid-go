# Passkey 前端改造清单

面向现有业务前端维护者，列出上线 Passkey 所需的改动项与核对项。

## 1. 路由与鉴权
- 确保用户中心（或登录页）在调用 Passkey 接口前已获取有效 JWT，并在请求头携带 `Authorization: Bearer <token>`。
- 如果后端域名为 `https://local.bsz.com:3000`，请核对浏览器调用时的 `origin` 与后端配置一致（配置项 `config.yaml` → `server.frontUrl`）。

## 2. 注册流程（账号设置界面）
1. 点击“绑定 Passkey”触发注册流程：
   - `GET /passkey/register/options` → 取得 `PublicKeyCredentialCreationOptions`
   - 调用 `navigator.credentials.create({ publicKey: options })`
2. 将返回的 `PublicKeyCredential` 转换为 JSON 并提交：
   - `POST /passkey/register`，`Content-Type: application/json`
3. 注册成功后刷新列表或提示“绑定成功”。

> ⚠️ Safari、Chrome 等浏览器要求页面为 HTTPS 且顶级域一致；请在本地调试时使用 HTTPS。

## 3. 登录流程（登录页 Passkey 按钮）
1. 登录表单新增“使用 Passkey 登录”按钮。
2. 点击按钮流程：
   - `GET /passkey/login/options`
   - `navigator.credentials.get({ publicKey: options })`
   - `POST /passkey/login`
   - 如成功，后端会返回新的 JWT（`data.token`），应覆盖现有登录态并跳转后台首页。
3. 若接口返回 `未绑定 Passkey`，提示用户先在个人中心绑定。

## 4. 凭证管理页（可选）
- 通过 `GET /passkey` 展示当前账号绑定的 Passkey 列表：
  - 展示字段建议：创建时间、最近使用时间、是否存在 Clone Warning。
- 删除按钮调用 `DELETE /passkey/:id`。
- 注册成功后自动刷新本列表。

## 5. JS 工具函数
- 确保项目中存在以下工具或等效处理，用于序列化 WebAuthn `ArrayBuffer`：

```ts
const toBase64Url = (buffer: ArrayBuffer): string => {
  const bytes = new Uint8Array(buffer);
  let binary = "";
  bytes.forEach(b => (binary += String.fromCharCode(b)));
  return btoa(binary).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/g, "");
};

export const publicKeyCredentialToJSON = (cred: PublicKeyCredential) => ({
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

## 6. 错误处理与 UX
- 注册/登录出现 `挑战已过期，请重试` → 自动重新获取 options 并提醒用户重新操作。
- 捕获 `navigator.credentials.*` 抛出的 `NotAllowedError`（用户取消）等异常，提示“操作已取消”。
- 若浏览器不支持 WebAuthn（`window.PublicKeyCredential` 不存在），隐藏按钮或给出提示。

## 7. 联调与测试
- 浏览器控制台搜索网络请求是否正确发送 JSON（非 `application/x-www-form-urlencoded`）。
- 本地调试请确保后端 Redis 已启动，避免接口返回挑战不存在。
- 常用测试场景：
  1. 首次绑定 Passkey
  2. 重复绑定（应覆盖旧记录）
  3. 登录成功获取新 JWT
  4. 删除 Passkey 后重新登录（应提示无绑定）

完成以上事项后，即可在前端完整支持 Passkey 注册与登录。