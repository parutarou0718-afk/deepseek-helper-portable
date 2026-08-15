# DeepSeek 助手便携版（DeepSeek Helper Portable）

免安装、双击即用的 DeepSeek AI 助手一键包（Windows x64）。

## ✨ 特点
- **免安装**：解压即用，不用装 Node.js / npm，不用配环境
- **一键启动**：双击 `启动DeepSeek助手.exe` → 自动启动 → 浏览器自动打开
- **自带加载页**：启动过程有「正在启动…」提示，不会误以为卡死
- **本地保存**：聊天记录与设置存在本地
- **冻结依赖**：环境全部内置，不再被「版本号/依赖报错」劝退

## 📥 下载
👉 最新版本：[v1.0.0 下载](https://github.com/parutarou0718-afk/deepseek-helper-portable/releases/download/v1.0.0/DeepSeek-Helper-Portable-v1.0.0.zip)

## 🚀 使用
1. 解压 `DeepSeek-Helper-Portable-v1.0.0.zip`（别放 C 盘 Program Files）
2. 双击 `启动DeepSeek助手.exe`
3. 首次按提示填写自己的 DeepSeek API Key（在 [platform.deepseek.com](https://platform.deepseek.com) 注册获取）

## ⚠️ 说明
- 本工具是 DeepSeek Harness（`@deepseek-ai/dsh`，MIT 协议）的便携封装，**非 DeepSeek 官方出品**。
- 需要联网，调用 DeepSeek 官方 API，需自备 API Key，按量计费。
- 仅支持 Windows 10 / 11 64 位。

## 目录结构
- `main.go` / `go.mod`：启动器源码（Go）
- `loading.html`：启动加载页
- `landing/index.html`：下载落地页（可部署到你的网站）
- `使用说明.txt`：面向使用者的说明
