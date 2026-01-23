#!/bin/bash
# Glow 应用发布脚本
# 用于部署 simple-app 示例应用

set -e  # 遇到错误立即退出

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
GLOW_CLI="$PROJECT_ROOT/glow"
APP_NAME="simple-app"
APP_YAML="$SCRIPT_DIR/app.yaml"

echo "=========================================="
echo "Glow 应用发布流程"
echo "=========================================="
echo ""

# 1. 检查 glow CLI 是否存在
echo "1️⃣  检查 glow CLI..."
if [ ! -f "$GLOW_CLI" ]; then
    echo "❌ 错误: glow CLI 不存在于 $GLOW_CLI"
    echo "请先运行: go build -o glow ./cmd/glow"
    exit 1
fi
echo "✅ glow CLI 存在: $GLOW_CLI"
echo ""

# 2. 检查 glow-server 是否运行
echo "2️⃣  检查 glow-server 状态..."
if ! curl -s http://localhost:32102/health > /dev/null 2>&1; then
    echo "❌ 错误: glow-server 未运行"
    echo "请先启动: sudo ./glow-server serve"
    exit 1
fi
echo "✅ glow-server 正在运行"
echo ""

# 3. 检查 API 配置
echo "3️⃣  检查 CLI 配置..."
if ! "$GLOW_CLI" context list > /dev/null 2>&1; then
    echo "❌ 错误: CLI 未配置"
    echo "请运行: ./glow context add default --server-url http://localhost:32102 --api-key <your-key>"
    exit 1
fi
echo "✅ CLI 已配置"
echo ""

# 4. 检查应用配置文件
echo "4️⃣  检查应用配置文件..."
if [ ! -f "$APP_YAML" ]; then
    echo "❌ 错误: app.yaml 不存在于 $APP_YAML"
    exit 1
fi
echo "✅ 配置文件: $APP_YAML"
echo ""

# 5. 显示配置内容
echo "📄 应用配置内容:"
echo "----------------------------------------"
cat "$APP_YAML"
echo "----------------------------------------"
echo ""

# 6. 确认发布
read -p "是否继续发布应用? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 发布已取消"
    exit 0
fi

echo ""
echo "=========================================="
echo "开始发布应用: $APP_NAME"
echo "=========================================="
echo ""

# 7. 应用配置（创建资源、绑定域名等）
echo "7️⃣  应用配置 (glow apply)..."
cd "$SCRIPT_DIR"
if "$GLOW_CLI" apply -f app.yaml; then
    echo "✅ 配置应用成功"
else
    echo "❌ 配置应用失败"
    exit 1
fi
echo ""

# 8. 启动应用
echo "8️⃣  启动应用 (glow start app)..."
if "$GLOW_CLI" start app "$APP_NAME"; then
    echo "✅ 应用启动成功"
else
    echo "❌ 应用启动失败"
    exit 1
fi
echo ""

# 9. 等待应用启动
echo "9️⃣  等待应用就绪..."
sleep 2

# 10. 检查应用状态
echo "🔟 检查应用状态..."
"$GLOW_CLI" get app "$APP_NAME"
echo ""

# 11. 测试应用
echo "1️⃣1️⃣  测试应用访问..."
if curl -s http://app.local/ > /dev/null 2>&1; then
    RESPONSE=$(curl -s http://app.local/)
    echo "✅ 应用响应: $RESPONSE"
else
    echo "⚠️  警告: 无法访问 http://app.local/"
    echo "   请确保 /etc/hosts 包含: 127.0.0.1 app.local"
fi
echo ""

# 12. 完成
echo "=========================================="
echo "✅ 应用发布完成！"
echo "=========================================="
echo ""
echo "应用信息:"
echo "  名称: $APP_NAME"
echo "  端口: 33203"
echo "  域名: http://app.local"
echo ""
echo "常用命令:"
echo "  查看状态: $GLOW_CLI get app $APP_NAME"
echo "  查看日志: $GLOW_CLI logs $APP_NAME"
echo "  重启应用: $GLOW_CLI restart app $APP_NAME"
echo "  停止应用: $GLOW_CLI stop app $APP_NAME"
echo "  删除应用: $GLOW_CLI delete app $APP_NAME"
echo ""
echo "配置文件位置:"
echo "  /var/lib/glow-server/apps/$APP_NAME/${APP_NAME}_local_config.json"
echo ""
