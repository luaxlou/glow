#!/bin/bash
# Glow 应用测试脚本
# 用于在 glow-server 重启后测试完整的发布流程

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
GLOW_CLI="$PROJECT_ROOT/glow"
APP_NAME="simple-app"

echo "=========================================="
echo "Glow 应用发布测试"
echo "=========================================="
echo ""

# 1. 检查 glow-server 是否运行
echo "1️⃣  检查 glow-server 状态..."
if ! curl -s http://localhost:32102/health > /dev/null 2>&1; then
    echo "❌ 错误: glow-server 未运行"
    echo "请先执行以下步骤重启服务器："
    echo "  cd $PROJECT_ROOT"
    echo "  sudo pkill -f glow-server"
    echo "  sudo mv glow-server-new glow-server"
    echo "  sudo ./glow-server serve"
    exit 1
fi
echo "✅ glow-server 正在运行"
echo ""

# 2. 测试新的 API 路由
echo "2️⃣  测试新的 API 路由 (PUT /apps/:name)..."
TEST_RESPONSE=$(curl -s -X PUT http://localhost:32102/apps/test-route-check \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer test" \
    -d '{"name":"test"}' 2>&1)

if echo "$TEST_RESPONSE" | grep -q "404"; then
    echo "❌ 错误: PUT /apps/:name 路由不存在"
    echo "请确认已重启 glow-server 到最新版本"
    echo "响应: $TEST_RESPONSE"
    exit 1
fi
echo "✅ PUT /apps/:name 路由已加载"
echo ""

# 3. 清理旧应用（如果存在）
echo "3️⃣  清理旧的应用实例（如果存在）..."
if "$GLOW_CLI" get app "$APP_NAME" > /dev/null 2>&1; then
    echo "发现旧的应用，停止中..."
    "$GLOW_CLI" stop "$APP_NAME" > /dev/null 2>&1 || true
    sleep 1
fi
echo "✅ 清理完成"
echo ""

# 4. 应用配置
echo "4️⃣  应用配置 (glow apply)..."
cd "$SCRIPT_DIR"
if ! "$GLOW_CLI" apply -f app.yaml; then
    echo "❌ 配置应用失败"
    exit 1
fi
echo ""

# 5. 启动应用
echo "5️⃣  启动应用..."
if ! "$GLOW_CLI" start app "$APP_NAME"; then
    echo "❌ 启动应用失败"
    exit 1
fi
echo ""

# 6. 等待应用启动
echo "6️⃣  等待应用就绪..."
sleep 3

# 7. 检查应用状态
echo "7️⃣  检查应用状态..."
echo "----------------------------------------"
"$GLOW_CLI" get app "$APP_NAME"
echo "----------------------------------------"
echo ""

# 8. 查看应用日志
echo "8️⃣  应用日志（最近 10 行）..."
echo "----------------------------------------"
"$GLOW_CLI" logs "$APP_NAME" 2>&1 | tail -10
echo "----------------------------------------"
echo ""

# 9. 测试应用访问
echo "9️⃣  测试应用访问..."
if curl -s http://app.local/ > /dev/null 2>&1; then
    RESPONSE=$(curl -s http://app.local/)
    echo "✅ 应用响应成功"
    echo "响应内容: $RESPONSE"
else
    echo "⚠️  无法访问 http://app.local/"
    echo "   请确保 /etc/hosts 包含: 127.0.0.1 app.local"
fi
echo ""

# 10. 检查配置文件
echo "🔟 检查生成的配置文件..."
CONFIG_FILE="/var/lib/glow-server/apps/$APP_NAME/${APP_NAME}_local_config.json"
if [ -f "$CONFIG_FILE" ]; then
    echo "✅ 配置文件已生成: $CONFIG_FILE"
    echo "内容预览:"
    cat "$CONFIG_FILE" | head -5
else
    echo "⚠️  配置文件不存在: $CONFIG_FILE"
fi
echo ""

# 11. 完成
echo "=========================================="
echo "✅ 测试完成！"
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
echo "  停止应用: $GLOW_CLI stop $APP_NAME"
echo ""
echo "如需完全清理:"
echo "  $GLOW_CLI delete app $APP_NAME"
echo ""
