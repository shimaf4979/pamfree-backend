#!/bin/bash
# save as scripts/cleanup-disk.sh

echo "=== ディスク容量クリーンアップ ==="
df -h
echo "未使用のDockerリソースを削除中..."
sudo docker system prune -a -f --volumes

echo "一時ファイルを削除中..."
sudo rm -rf /tmp/*
sudo rm -rf /var/tmp/*

# キャッシュディレクトリのクリーンアップ
sudo find /var/cache -type f -delete

# ログファイルのクリーンアップ
sudo find /var/log -type f -name "*.gz" -delete
sudo find /var/log -type f -name "*.old" -delete
sudo find /var/log -type f -name "*.1" -delete

# 古いカーネルの削除
sudo package-cleanup --oldkernels --count=1 -y || true

echo "クリーンアップ後のディスク容量:"
df -h