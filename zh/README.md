# Go Web 编程 PDF 生成说明

本目录包含了生成中文版 PDF 的相关文件和脚本。

## 环境要求

1. Go 环境
2. pandoc（用于转换 HTML 到 PDF）
3. LaTeX 环境（建议使用 MacTeX）
4. GNU sed（在 macOS 上需要安装 `gsed`）
5. GitHub Token（用于访问 GitHub API）

## 字体要求

- SF Mono - 代码字体
- 宋体-简 (Songti SC) - 正文字体
- 苹方-简 (PingFang SC) - 标题字体
- 仿宋-简 (STFangsong) - 等宽中文字体

## 安装依赖

macOS 环境下：

```bash
# 安装 pandoc
brew install pandoc

# 安装 MacTeX
brew install --cask mactex

# 安装 GNU sed
brew install gnu-sed
```

## 配置

1. 创建 `.env` 文件并设置 GitHub Token：

```bash
echo "GITHUB_TOKEN=your_github_token" > .env
```

2. 确保所需字体已安装在系统中

## 生成 PDF

运行以下命令：

```bash
./build_pdf.sh
```

生成的 PDF 文件将保存为 `build-web-application-with-golang.pdf`。

## 目录结构

- `github_builder/` - 包含 HTML 生成工具
- `output/` - 临时生成的 HTML 文件
- `images/` - 图片资源
- `build_pdf.sh` - PDF 生成脚本
- `*.md` - 源码文件