# 知识笔记软件

## 项目概述
本项目为个人使用的知识笔记记录软件，旨在为用户提供便捷的文档管理、编辑及预览功能，帮助用户高效整理和查阅个人知识内容。

## 核心功能模块
### 1. 文档管理
- 展示所有文档列表，支持查看与快速定位
- 编辑单个文档的基础信息（如标题、描述、标签等）

### 2. 文档编辑页面
- 从文档列表点击进入指定文档的编辑界面
- 左侧显示文档目录结构，辅助定位各章节
- 右侧提供 Markdown 编辑器，用于编写和修改文档内容

### 3. 文档全文预览页面
- 以全屏方式呈现文档内容
- 仅提供阅读体验，不包含任何编辑操作

## 开发

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

## Docker 部署

### 构建镜像
```bash
docker build -t zw-note:latest .
```

### 运行容器
```bash
docker run -d -p 8080:80 --name zw-note zw-note:latest
```

访问 http://localhost:8080 即可使用应用。

### 使用 Docker Compose（可选）
创建 `docker-compose.yml` 文件：

```yaml
version: '3.8'
services:
  zw-note:
    build: .
    ports:
      - "8080:80"
    restart: unless-stopped
```

运行：
```bash
docker-compose up -d
```
