#!/bin/bash
set -e

echo "=== t2.micro上でのビルドを最適化 ==="

# システム情報を表示
echo "システム情報:"
free -h
df -h
echo "利用可能なメモリとディスク容量を確認しました"

# ディスククリーンアップ (最初に実行して容量を確保)
echo "ディスク容量をクリーンアップ中..."
sudo docker system prune -af --volumes
sudo apt-get clean || sudo yum clean all || true
sudo rm -rf /tmp/* /var/tmp/* || true
echo "クリーンアップ完了"

# EC2のスワップファイル設定 (メモリ不足対策)
if [ ! -f /swapfile ]; then
    echo "スワップファイルを作成中..."
    sudo swapoff -a || true
    sudo rm -f /swapfile || true
    sudo dd if=/dev/zero of=/swapfile bs=64M count=16
    sudo chmod 600 /swapfile
    sudo mkswap /swapfile
    sudo swapon /swapfile
    echo "スワップファイルを有効化しました (1GB)"
    free -h
else
    echo "既存のスワップファイルを使用します"
    sudo swapon /swapfile || true
fi

# Dockerのメモリ設定の最適化
echo "Dockerの設定を最適化中..."
mkdir -p ~/.docker
echo '{
  "experimental": true,
  "builder": {
    "gc": {
      "enabled": true,
      "defaultKeepStorage": "1GB"
    }
  },
  "storage-driver": "overlay2",
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}' > ~/.docker/config.json

# Docker関連のファイルシステムの最適化
sudo rm -rf /var/lib/docker/tmp/* || true

# Docker設定を再起動して適用
echo "Dockerを再起動して設定を適用中..."
sudo systemctl restart docker
echo "Dockerを再起動しました"
sleep 3

# ビルドに必要な環境変数を設定
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1
export DOCKER_MEMORY_LIMIT=1024m
export GOGC=20
export GOPROXY=https://proxy.golang.org,direct

# ビルド実行
echo "最適化されたビルドを開始します..."
echo "注意: 初回ビルドは時間がかかります。コーヒーを飲みながらお待ちください☕"

# 実際のビルドコマンド
# ディスク問題を回避するため、--no-cacheは使用しない
docker-compose -f docker-compose.prod.yml build --build-arg DOCKER_BUILDKIT=1 --build-arg GOGC=20

# ビルド後のクリーンアップ
echo "中間コンテナやキャッシュをクリーンアップ中..."
sudo docker image prune -f
# キャッシュは消さないように --volumes は指定しない
sudo docker container prune -f

echo ""
echo "✅ ビルドが完了しました！"
echo "make prod を実行して起動できます"