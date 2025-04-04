.PHONY: help prod prod-logs stop restart status clean cleanup micro-build logs

# デフォルトのターゲット
help:
	@echo "== パンフレットアプリケーション Makefile =="
	@echo ""
	@echo "基本コマンド:"
	@echo "  make prod           - 本番用アプリケーションを起動"
	@echo "  make prod-logs      - 本番用アプリケーションのログを表示"
	@echo "  make stop           - アプリケーションを停止"
	@echo "  make restart        - アプリケーションを再起動"
	@echo "  make status         - アプリケーションの状態を確認"
	@echo ""
	@echo "ビルド関連:"
	@echo "  make micro-build    - t2.micro向けに最適化したビルドを実行"
	@echo ""
	@echo "メンテナンス:"
	@echo "  make cleanup        - ディスク容量をクリーンアップ"
	@echo "  make clean          - すべてのコンテナとボリュームを削除"
	@echo ""
	@echo "詳細情報: README.md を参照してください"

# 本番用アプリケーションを実行
prod:
	@echo "本番用アプリケーションを起動中..."
	@sudo docker network inspect pamfree_network >/dev/null 2>&1 || sudo docker network create pamfree_network
	sudo docker-compose -f docker-compose.prod.yml up -d
	@echo "アプリケーションが起動しました (http://localhost:8080)"
	@make status

# 本番用アプリケーションのログを表示
prod-logs:
	sudo docker-compose -f docker-compose.prod.yml logs -f api

# 短縮バージョン
logs: prod-logs

# アプリケーションを停止
stop:
	@echo "アプリケーションを停止中..."
	sudo docker-compose -f docker-compose.prod.yml down
	@echo "アプリケーションを停止しました"

# アプリケーションを再起動
restart: stop prod

# アプリケーションの状態を確認
status:
	@echo "アプリケーションの状態:"
	@sudo docker-compose -f docker-compose.prod.yml ps
	@echo ""
	@echo "システムリソース:"
	@echo "- メモリ:"
	@free -h | head -2
	@echo "- ディスク:"
	@df -h | grep -E "(Filesystem|/$$)"

# t2.micro向け最適化ビルド
micro-build:
	@echo "t2.micro向けの最適化ビルドを実行中..."
	@if [ -f ./build-on-micro.sh ]; then \
		./build-on-micro.sh; \
	else \
		echo "Error: build-on-micro.sh が見つかりません"; \
		exit 1; \
	fi
	@echo "ビルドが完了しました"

# ディスク容量クリーンアップ
cleanup:
	@echo "ディスク容量をクリーンアップ中..."
	@if [ -f ./cleanup.sh ]; then \
		./cleanup.sh; \
	else \
		echo "cleanup.sh が見つからないため、基本クリーンアップを実行します"; \
		sudo docker system prune -f; \
		sudo docker image prune -f; \
	fi
	@echo "クリーンアップが完了しました"
	@df -h | grep -E '(Filesystem|/$)'

# すべてのコンテナとボリュームを削除 (完全クリーンアップ)
clean:
	@echo "すべてのコンテナとボリュームを削除中..."
	sudo docker-compose -f docker-compose.prod.yml down -v
	sudo docker rmi -f pamfree-api:latest 2>/dev/null || true
	sudo docker volume rm pamfree_uploads 2>/dev/null || true
	@echo "クリーンアップが完了しました！"