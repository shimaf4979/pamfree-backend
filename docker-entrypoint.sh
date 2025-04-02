#!/bin/sh
set -e

echo "=== Processingプラットフォーム API サーバー ==="
echo "環境変数:"
echo "SERVER_PORT: $SERVER_PORT"
echo "DB_HOST: $DB_HOST"
echo "DB_NAME: $DB_NAME"
echo "UPLOAD_DIR: $UPLOAD_DIR"
echo "GIN_MODE: $GIN_MODE"
echo "========================================"

# 引数をそのまま実行
echo "サーバーを起動します..."
exec "$@"