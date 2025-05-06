# 本番用イメージ
FROM alpine:latest


# 作業ディレクトリ
WORKDIR /root/

# ビルド済みバイナリをコピー
COPY ./src/main ./



# 環境変数・ポート設定
ENV PORT=8080
EXPOSE 8080

# アプリケーション起動
CMD ["./main"]
