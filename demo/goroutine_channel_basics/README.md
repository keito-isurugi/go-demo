# goroutineとchannelの基礎

このプログラムは、Goのgoroutineとchannelの基本的な構文と使い方を学ぶためのデモです。

## 実行方法

```bash
cd demo/goroutine_channel_basics
go run main.go
```

## 学習内容

### Example 1: 基本的なgoroutine
- `go` キーワードで関数を並行実行
- goroutineの基本的な動作を理解

```go
go sayHello("Bob")  // goroutineで実行
```

### Example 2: 基本的なchannel
- channelの作成: `make(chan string)`
- 送信: `ch <- "メッセージ"`
- 受信: `message := <-ch`

### Example 3: channelで複数goroutineを同期
- 複数のworker goroutineが協調動作
- `for range` でchannelから連続受信
- `close()` でchannelの終了を通知

### Example 4: buffered channelとunbuffered channel
- **unbuffered**: `make(chan int)` - 送信側と受信側が揃うまでブロック
- **buffered**: `make(chan int, 2)` - バッファが満杯になるまでブロックしない

### Example 5: select文で複数channelを処理
- 複数のchannelを同時に待つ
- どれか1つが準備できたら処理

```go
select {
case msg1 := <-ch1:
    // ch1から受信
case msg2 := <-ch2:
    // ch2から受信
}
```

### Example 6: タイムアウト処理
- `time.After()` を使ったタイムアウト実装
- 長時間かかる処理の制限

### Example 7: done channelパターン
- goroutineの完了を待つパターン
- シンプルな同期方法

### Example 8: close()でchannelを閉じる
- `close(ch)` でchannelを閉じる
- `for range` でchannelが閉じられるまでループ
- 受信時の2値返却: `val, ok := <-ch`

## 重要なポイント

1. **goroutine**: `go` キーワードで関数を並行実行
2. **channel**: goroutine間の通信と同期に使用
3. **矢印の向き**:
   - `ch <- value` (送信)
   - `value := <-ch` (受信)
4. **バッファの有無**: 同期の挙動が変わる
5. **select**: 複数channelの多重化
6. **close**: channelの終了を通知（送信側が行う）

## 学習の進め方

1. まず全体を実行して動作を確認
2. 各exampleのコードを読んで理解
3. 値やタイミングを変更して挙動を観察
4. コメントアウトして1つずつ実行してみる
