# ベースイメージとして公式のGoイメージを使用
FROM golang:1.20-alpine

# 作業ディレクトリを設定
WORKDIR /app

# Go Modulesの利用を考慮して、go.modとgo.sumをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# Goアプリケーションをビルド
RUN go build -o main .

# アプリケーションを実行
CMD ["./main"]
