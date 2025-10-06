# RabbitMQ Producer-Consumer Demo

このデモは RabbitMQ を使った基本的な Producer-Consumer パターンの実装例です。

## 構成

- **Producer**: メッセージをキューに送信する
- **Consumer**: キューからメッセージを受信する
- **Queue**: `hello` という名前のキューを使用

## セットアップ

### 1. RabbitMQ を起動

```bash
# プロジェクトルートで実行
docker-compose up -d rabbitmq
```

RabbitMQ が起動すると以下のポートでアクセスできます：
- **AMQP**: `localhost:5672` (アプリケーション接続用)
- **Management UI**: `http://localhost:15672` (ブラウザから確認可能)
  - ユーザー名: `guest`
  - パスワード: `guest`

### 2. 依存関係をインストール

```bash
cd demo/rabbitmq
go mod download
```

## 使い方

### Producer を実行

メッセージをキューに送信します：

```bash
go run main.go -mode=producer
```

出力例：
```
Starting Producer...
 [x] Sent: Hello RabbitMQ! Message #1
 [x] Sent: Hello RabbitMQ! Message #2
 [x] Sent: Hello RabbitMQ! Message #3
 [x] Sent: Hello RabbitMQ! Message #4
 [x] Sent: Hello RabbitMQ! Message #5
Producer finished
```

### Consumer を実行

キューからメッセージを受信します（別のターミナルで実行）：

```bash
go run main.go -mode=consumer
```

出力例：
```
Starting Consumer...
 [*] Waiting for messages. To exit press CTRL+C
 [x] Received: Hello RabbitMQ! Message #1
 [x] Received: Hello RabbitMQ! Message #2
 [x] Received: Hello RabbitMQ! Message #3
 [x] Received: Hello RabbitMQ! Message #4
 [x] Received: Hello RabbitMQ! Message #5
```

終了する場合は `CTRL+C` を押してください。

## 動作確認の手順

1. Consumer を起動しておく
2. 別のターミナルで Producer を実行
3. Consumer がメッセージを受信することを確認

## Management UI での確認

ブラウザで `http://localhost:15672` にアクセスすると、RabbitMQ の管理画面が表示されます。
ここで以下の情報を確認できます：

- キューの状態
- メッセージ数
- Consumer の接続状態
- メッセージの送受信レート

## コード説明

### producer.go
- `NewProducer()`: RabbitMQ への接続とキューの宣言
- `Publish()`: メッセージをキューに送信
- `Close()`: 接続のクリーンアップ

### consumer.go
- `NewConsumer()`: RabbitMQ への接続とキューの宣言
- `Consume()`: メッセージを受信し続ける（ブロッキング）
- `Close()`: 接続のクリーンアップ

### main.go
- コマンドライン引数で Producer/Consumer を切り替え
- Producer: 5つのメッセージを1秒間隔で送信
- Consumer: メッセージを受信し続ける

## 学習ポイント

1. **メッセージキューの基本**: Producer と Consumer が非同期で動作
2. **疎結合**: Producer と Consumer は互いに依存しない
3. **スケーラビリティ**: 複数の Consumer を起動して負荷分散可能
4. **信頼性**: RabbitMQ がメッセージを永続化（durable 設定時）

## 次のステップ

- メッセージの永続化（durable queue）
- メッセージの確認応答（manual ack）
- Work Queue パターン（複数 Consumer）
- Pub/Sub パターン（Exchange 使用）
- Routing パターン
