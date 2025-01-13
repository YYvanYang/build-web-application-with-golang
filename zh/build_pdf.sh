#!/bin/bash

SED='sed'

if [ `uname -s` == 'Darwin' ] ; then
  SED='gsed'
fi

bn="`basename $0`"
WORKDIR="$(cd $(dirname $0); pwd -P)"
OUTPUT_DIR="$WORKDIR/output"

# 从 .env 文件加载环境变量
if [ -f "$WORKDIR/.env" ]; then
    export $(cat "$WORKDIR/.env" | xargs)
else
    echo "未找到 .env 文件，请创建 .env 文件并设置 GITHUB_TOKEN"
    echo "文件格式示例："
    echo "GITHUB_TOKEN=your_github_token"
    exit 1
fi

# 检查 GitHub Token
if [ -z "$GITHUB_TOKEN" ]; then
    echo "在 .env 文件中未找到 GITHUB_TOKEN"
    echo "请在 .env 文件中设置 GITHUB_TOKEN"
    exit 1
fi

#
# Default language: zh
# You can overwrite following variables in config file.
#
MSG_INSTALL_PANDOC_FIRST='请先安装pandoc，然后再次运行'
MSG_SUCCESSFULLY_GENERATED='build-web-application-with-golang.pdf 已经建立'
MSG_CREATOR='Astaxie'
MSG_DESCRIPTION='一本开源的Go Web编程书籍'
MSG_LANGUAGE='zh-CN'
MSG_TITLE='Go Web编程'
[ -e "$WORKDIR/config" ] && . "$WORKDIR/config"

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 使用build_new.go生成HTML文件
(
cd "$WORKDIR/github_builder" && go run build_new.go
)

if ! type pandoc >/dev/null 2>&1; then
    echo "$MSG_INSTALL_PANDOC_FIRST"
    exit 1
fi

# 复制图片
mkdir -p "$OUTPUT_DIR/images"
cp -r "$WORKDIR/images/"* "$OUTPUT_DIR/images/" 2>/dev/null || true
ls "$OUTPUT_DIR/"[0-9]*.html 2>/dev/null | xargs $SED -i "s/png?raw=true/png/g" || true

echo "工作目录: $WORKDIR"
echo "输出目录: $OUTPUT_DIR"

# 使用 pandoc 将 HTML 转换为 PDF
echo "正在生成 PDF 文件..."
pandoc -s --pdf-engine=xelatex \
  --toc \
  --toc-depth=2 \
  -V documentclass=report \
  -V CJKmainfont="Songti SC" \
  -V CJKsansfont="PingFang SC" \
  -V monofont="SF Mono" \
  -V CJKmonofont="STFangsong" \
  -V geometry:margin=1in \
  --verbose \
  "$OUTPUT_DIR/"*.html -o build-web-application-with-golang.pdf

echo "$MSG_SUCCESSFULLY_GENERATED" 