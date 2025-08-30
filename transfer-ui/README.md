# Transfer UI

一个用于将网易云音乐歌单迁移到 Spotify 的前端应用。

## 功能特性

- 🎵 获取网易云音乐歌单
- 🔐 Spotify OAuth 授权
- 👤 Spotify 用户验证
- 🎯 歌曲选择
- 📥 歌曲迁移到 Spotify

## 安装依赖

```bash
npm install
```

## 开发运行

```bash
npm run dev
```

应用将在 `http://localhost:3000` 启动

## 构建

```bash
npm run build
```

## 项目结构

```
src/
├── components/           # React 组件
│   ├── StepIndicator.tsx    # 步骤指示器
│   ├── NetEaseInput.tsx     # 网易云歌单输入
│   ├── PlaylistDisplay.tsx  # 歌单展示
│   ├── SpotifyAuth.tsx      # Spotify 授权
│   ├── SpotifyUserConfirm.tsx # 用户确认
│   ├── TrackSelection.tsx   # 歌曲选择
│   └── TransferProcess.tsx  # 迁移处理
├── services/            # API 服务
│   └── api.ts          # 后端 API 调用
├── App.tsx             # 主应用组件
└── main.tsx           # 应用入口
```

## API 接口

前端通过 `/api` 代理与后端通信：

- `GET /api/netease/playlist` - 获取网易云歌单
- `GET /api/user/auth/spotify/login` - Spotify OAuth 登录
- `POST /api/user/auth/spotify/status` - 检查授权状态
- `GET /api/spotify/me` - 获取 Spotify 用户信息
- `GET /api/spotify/playlists` - 获取用户歌单
- `POST /api/spotify/playlists/:id/tracks` - 添加歌曲到歌单

## 使用流程

1. **输入网易云歌单ID** - 用户输入需要迁移的网易云音乐歌单ID
2. **查看歌单信息** - 系统获取并显示歌单详情
3. **Spotify 授权** - 用户授权应用访问 Spotify 账户
4. **确认用户身份** - 验证 Spotify 用户ID 并选择目标歌单
5. **选择歌曲** - 用户选择要迁移的歌曲
6. **执行迁移** - 系统将选中歌曲添加到 Spotify 歌单

## 注意事项

- 确保后端服务运行在 `http://localhost:8081`
- 需要有效的 Spotify 开发者账户和应用配置
- 网易云音乐歌单需要是公开可访问的
