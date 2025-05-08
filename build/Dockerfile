# 本番用イメージ
FROM alpine:latest

# 作業ディレクトリ
WORKDIR /root/

# ビルド済みバイナリをコピー
COPY src/main ./
RUN chmod +x ./main

# 環境変数・ポート設定
ENV PORT=80
EXPOSE 80

# アプリケーション起動
CMD ["./main"]
