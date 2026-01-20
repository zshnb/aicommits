#!/bin/bash
set -e

# ================= 配置项 =================
# 请修改这里为你自己的仓库信息
REPO_OWNER="zshnb"
REPO_NAME="aicommits"
BIN_NAME="aicommits"
# =========================================

# 检测操作系统和架构
OS="$(uname -s)"
ARCH="$(uname -m)"

case $OS in
    "Linux")
        case $ARCH in
        "x86_64")
            if [ "$(getconf LONG_BIT)" = "64" ]; then
                FILE_OS="Linux"
                FILE_ARCH="amd64"
            else
                echo "不支持 32 位 Linux"
                exit 1
            fi
            ;;
        "aarch64" | "arm64")
            FILE_OS="Linux"
            FILE_ARCH="arm64"
            ;;
        *)
            echo "不支持的架构: $ARCH"
            exit 1
            ;;
        esac
        ;;
    "Darwin")
        FILE_OS="Darwin"
        case $ARCH in
        "x86_64")
            FILE_ARCH="amd64"
            ;;
        "arm64")
            FILE_ARCH="arm64"
            ;;
        *)
            echo "不支持的架构: $ARCH"
            exit 1
            ;;
        esac
        ;;
    *)
        echo "不支持的系统: $OS"
        exit 1
        ;;
esac

# 构建下载 URL (GoReleaser 的默认命名格式)
# 格式示例: aicommits_Darwin_arm64.tar.gz
FILE_NAME="${REPO_NAME}_${FILE_OS}_${FILE_ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/latest/download/${FILE_NAME}"

INSTALL_PATH="/usr/local/bin/$BIN_NAME"
if [ -f "$INSTALL_PATH" ]; then
    echo "🔄 检测到已安装版本，准备升级..."
    IS_UPGRADE=true
else
    IS_UPGRADE=false
fi

echo "⬇️  正在下载 ${DOWNLOAD_URL}..."
tmp_dir=$(mktemp -d)
curl -sL "$DOWNLOAD_URL" -o "$tmp_dir/$FILE_NAME"

echo "📦 正在解压..."
tar -xzf "$tmp_dir/$FILE_NAME" -C "$tmp_dir"

echo "🚀 安装到 /usr/local/bin..."
# 检查是否有写权限
if [ -w "/usr/local/bin" ]; then
    mv "$tmp_dir/$BIN_NAME" "$INSTALL_PATH"
else
    sudo mv "$tmp_dir/$BIN_NAME" "$INSTALL_PATH"
fi

chmod +x "$INSTALL_PATH"
rm -rf "$tmp_dir"

if [ "$IS_UPGRADE" = true ]; then
    echo "✅ 升级成功！"
else
    echo "✅ 安装成功！请运行 '$BIN_NAME config' 进行初始化。"
fi
