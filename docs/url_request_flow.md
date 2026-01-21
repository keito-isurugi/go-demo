# URLを叩くと何が起きるか - 完全ガイド

ブラウザでURLを入力してからWebページが表示されるまでの全プロセスを解説します。

## 目次
1. [全体フロー概要](#全体フロー概要)
2. [Step 1: URL解析](#step-1-url解析)
3. [Step 2: DNS解決](#step-2-dns解決)
4. [Step 3: TCP接続](#step-3-tcp接続)
5. [Step 4: TLSハンドシェイク](#step-4-tlsハンドシェイク)
6. [Step 5: HTTPリクエスト送信](#step-5-httpリクエスト送信)
7. [Step 6: サーバー処理](#step-6-サーバー処理)
8. [Step 7: HTTPレスポンス受信](#step-7-httpレスポンス受信)
9. [Step 8: ブラウザレンダリング](#step-8-ブラウザレンダリング)
10. [まとめ](#まとめ)

---

## 全体フロー概要

ブラウザでURLを入力してからページが表示されるまでには、多くのステップが実行されます。

```mermaid
flowchart TB
    subgraph Browser["ブラウザ"]
        A[URL入力] --> B[URL解析]
    end

    subgraph DNS["DNS解決"]
        B --> C[キャッシュ確認]
        C --> D[DNSサーバー問い合わせ]
        D --> E[IPアドレス取得]
    end

    subgraph Connection["接続確立"]
        E --> F[TCP 3-way handshake]
        F --> G[TLSハンドシェイク]
    end

    subgraph Request["リクエスト/レスポンス"]
        G --> H[HTTPリクエスト送信]
        H --> I[サーバー処理]
        I --> J[HTTPレスポンス受信]
    end

    subgraph Render["レンダリング"]
        J --> K[HTML解析]
        K --> L[リソース取得]
        L --> M[ページ表示]
    end
```

---

## Step 1: URL解析

### URLの構造

ブラウザはまず入力されたURLを解析し、各コンポーネントを識別します。

```
https://www.example.com:443/path/to/page?query=value#section
└─┬──┘ └──────┬───────┘└┬┘ └────┬────┘ └────┬────┘ └──┬───┘
  │           │         │       │           │         │
scheme     host       port    path        query    fragment
```

| コンポーネント | 説明 | 例 |
|--------------|------|-----|
| scheme | プロトコル | `https`, `http` |
| host | ドメイン名またはIPアドレス | `www.example.com` |
| port | ポート番号（省略時はデフォルト） | `443`（HTTPS）, `80`（HTTP） |
| path | リソースへのパス | `/path/to/page` |
| query | クエリパラメータ | `query=value` |
| fragment | ページ内の位置 | `section` |

### ブラウザキャッシュの確認

URL解析後、ブラウザは以下の順序でキャッシュを確認します：

```mermaid
flowchart LR
    A[URLリクエスト] --> B{ブラウザ<br/>キャッシュ?}
    B -->|Hit| C[キャッシュから返却]
    B -->|Miss| D{Service Worker<br/>キャッシュ?}
    D -->|Hit| C
    D -->|Miss| E[ネットワークリクエスト]
```

---

## Step 2: DNS解決

### DNSとは

**DNS**（Domain Name System）は、人間が読めるドメイン名をIPアドレスに変換するシステムです。

### 解決すべき課題

| 課題 | DNSによる解決 |
|-----|-------------|
| IPアドレスは覚えにくい | ドメイン名（例: google.com）で接続可能に |
| サーバーIPの変更 | DNSレコードの更新で対応、ユーザーは意識不要 |
| 負荷分散 | 複数IPへの振り分けが可能 |

### DNS解決の流れ

```mermaid
sequenceDiagram
    participant Browser as ブラウザ
    participant Local as ローカルキャッシュ
    participant Resolver as DNSリゾルバ
    participant Root as ルートDNS
    participant TLD as TLD DNS (.com)
    participant Auth as 権威DNS

    Browser->>Local: example.comのIP?
    Local-->>Browser: キャッシュなし

    Browser->>Resolver: example.comのIP?
    Resolver->>Root: .comの管理サーバーは?
    Root-->>Resolver: TLD DNSのアドレス

    Resolver->>TLD: example.comの管理サーバーは?
    TLD-->>Resolver: 権威DNSのアドレス

    Resolver->>Auth: example.comのIP?
    Auth-->>Resolver: 93.184.216.34

    Resolver-->>Browser: 93.184.216.34
    Browser->>Local: キャッシュに保存
```

### DNSレコードの種類

| レコード | 用途 | 例 |
|---------|------|-----|
| A | ドメイン → IPv4 | `example.com → 93.184.216.34` |
| AAAA | ドメイン → IPv6 | `example.com → 2001:db8::1` |
| CNAME | ドメインの別名 | `www.example.com → example.com` |
| MX | メールサーバー | `example.com → mail.example.com` |
| NS | ネームサーバー | `example.com → ns1.example.com` |

### DNSキャッシュの階層

```mermaid
flowchart TB
    A[ブラウザキャッシュ] --> B[OSキャッシュ]
    B --> C[ルーターキャッシュ]
    C --> D[ISP DNSキャッシュ]
    D --> E[DNSリゾルバ]
    E --> F[権威DNSサーバー]

    style A fill:#e1f5fe
    style B fill:#e1f5fe
    style C fill:#e1f5fe
    style D fill:#b3e5fc
    style E fill:#81d4fa
    style F fill:#4fc3f7
```

---

## Step 3: TCP接続

### TCPとは

**TCP**（Transmission Control Protocol）は、信頼性のあるデータ転送を提供するプロトコルです。

### 解決すべき課題

| 課題 | TCPによる解決 |
|-----|-------------|
| パケットロス | 再送制御で確実に届ける |
| パケットの順序乱れ | シーケンス番号で順序を保証 |
| データ破損 | チェックサムで検証 |
| 輻輳（ネットワーク混雑） | フロー制御・輻輳制御 |

### 3-way handshake

TCPでは接続確立のために3つのパケットをやり取りします。

```mermaid
sequenceDiagram
    participant Client as クライアント
    participant Server as サーバー

    Note over Client,Server: 3-way handshake

    Client->>Server: SYN (seq=x)
    Note right of Client: 接続要求

    Server->>Client: SYN-ACK (seq=y, ack=x+1)
    Note left of Server: 接続要求 + 確認応答

    Client->>Server: ACK (ack=y+1)
    Note right of Client: 確認応答

    Note over Client,Server: 接続確立完了
```

### なぜ3-wayなのか

| 目的 | 説明 |
|-----|------|
| 双方向の通信確認 | クライアント→サーバー、サーバー→クライアント両方の疎通確認 |
| 初期シーケンス番号の交換 | データの順序制御のため |
| なりすまし防止 | シーケンス番号による検証 |

---

## Step 4: TLSハンドシェイク

### TLSとは

**TLS**（Transport Layer Security）は、通信を暗号化するプロトコルです。HTTPSではTCPの上にTLSを使用します。

### 解決すべき課題

| 課題 | TLSによる解決 |
|-----|-------------|
| 盗聴 | 通信の暗号化 |
| なりすまし | 証明書による認証 |
| 改ざん | メッセージ認証コード（MAC） |

### TLS 1.3 ハンドシェイク

TLS 1.3では往復回数が削減され、高速化されています。

```mermaid
sequenceDiagram
    participant Client as クライアント
    participant Server as サーバー

    Note over Client,Server: TLS 1.3 ハンドシェイク (1-RTT)

    Client->>Server: ClientHello<br/>+ 鍵交換パラメータ<br/>+ 対応暗号スイート

    Server->>Client: ServerHello<br/>+ 鍵交換パラメータ<br/>+ 選択暗号スイート<br/>+ 証明書<br/>+ Finished

    Note over Client: 証明書検証<br/>共通鍵生成

    Client->>Server: Finished

    Note over Client,Server: 暗号化通信開始
```

### プロトコルスタック

```mermaid
flowchart TB
    subgraph HTTP["アプリケーション層"]
        A[HTTP/HTTPS]
    end

    subgraph TLS["セキュリティ層"]
        B[TLS]
    end

    subgraph TCP["トランスポート層"]
        C[TCP]
    end

    subgraph IP["ネットワーク層"]
        D[IP]
    end

    subgraph Link["リンク層"]
        E[Ethernet/WiFi]
    end

    A --> B --> C --> D --> E
```

---

## Step 5: HTTPリクエスト送信

### HTTPリクエストの構造

```
GET /path/to/page HTTP/1.1
Host: www.example.com
User-Agent: Mozilla/5.0 ...
Accept: text/html,application/xhtml+xml
Accept-Language: ja,en;q=0.9
Accept-Encoding: gzip, deflate, br
Connection: keep-alive
Cookie: session_id=abc123
```

### リクエストの構成要素

| 要素 | 説明 | 例 |
|-----|------|-----|
| メソッド | 操作の種類 | `GET`, `POST`, `PUT`, `DELETE` |
| パス | リソースの場所 | `/path/to/page` |
| HTTPバージョン | プロトコルバージョン | `HTTP/1.1`, `HTTP/2` |
| ヘッダー | メタ情報 | `Host`, `User-Agent`, `Accept` |
| ボディ | リクエストデータ | POSTデータなど |

### HTTP/1.1 vs HTTP/2 vs HTTP/3

```mermaid
flowchart LR
    subgraph HTTP1["HTTP/1.1"]
        direction TB
        A1[リクエスト1] --> A2[レスポンス1]
        A2 --> A3[リクエスト2]
        A3 --> A4[レスポンス2]
    end

    subgraph HTTP2["HTTP/2"]
        direction TB
        B1[リクエスト1] --> B3[レスポンス1]
        B2[リクエスト2] --> B4[レスポンス2]
    end

    subgraph HTTP3["HTTP/3"]
        direction TB
        C1[リクエスト1] --> C3[レスポンス1]
        C2[リクエスト2] --> C4[レスポンス2]
        C5[UDP/QUIC<br/>低遅延]
    end
```

| バージョン | 特徴 |
|-----------|------|
| HTTP/1.1 | 1リクエストずつ順番に処理（Head-of-Line Blocking） |
| HTTP/2 | 1接続で複数リクエストを多重化（ストリーム） |
| HTTP/3 | UDP/QUICベースで接続確立が高速、パケットロス耐性向上 |

---

## Step 6: サーバー処理

### サーバーがリクエストを処理する流れ

```mermaid
flowchart TB
    A[リクエスト受信] --> B[ロードバランサー]
    B --> C[Webサーバー<br/>nginx/Apache]
    C --> D{静的コンテンツ?}
    D -->|Yes| E[ファイル返却]
    D -->|No| F[アプリケーションサーバー]
    F --> G[ルーティング]
    G --> H[ミドルウェア処理<br/>認証/ログ等]
    H --> I[ビジネスロジック]
    I --> J{データ必要?}
    J -->|Yes| K[(データベース)]
    K --> L[データ取得]
    L --> M[レスポンス生成]
    J -->|No| M
    M --> N[レスポンス送信]
```

### 主要コンポーネント

| コンポーネント | 役割 | 例 |
|--------------|------|-----|
| ロードバランサー | リクエストの分散 | AWS ALB, nginx |
| リバースプロキシ | キャッシュ、SSL終端 | nginx, Cloudflare |
| Webサーバー | HTTPリクエスト処理 | nginx, Apache |
| アプリケーションサーバー | ビジネスロジック | Go, Node.js, Python |
| データベース | データ永続化 | PostgreSQL, MySQL |
| キャッシュ | 高速データアクセス | Redis, Memcached |

---

## Step 7: HTTPレスポンス受信

### HTTPレスポンスの構造

```
HTTP/1.1 200 OK
Date: Wed, 05 Feb 2026 12:00:00 GMT
Content-Type: text/html; charset=utf-8
Content-Length: 1234
Content-Encoding: gzip
Cache-Control: max-age=3600
Connection: keep-alive

<!DOCTYPE html>
<html>
  <head>...</head>
  <body>...</body>
</html>
```

### 主要なステータスコード

| コード | 意味 | 説明 |
|-------|------|------|
| 200 | OK | 正常終了 |
| 301 | Moved Permanently | 恒久的リダイレクト |
| 302 | Found | 一時的リダイレクト |
| 304 | Not Modified | キャッシュ利用可 |
| 400 | Bad Request | リクエスト不正 |
| 401 | Unauthorized | 認証必要 |
| 403 | Forbidden | アクセス拒否 |
| 404 | Not Found | リソース不在 |
| 500 | Internal Server Error | サーバー内部エラー |
| 503 | Service Unavailable | サービス一時停止 |

### レスポンスヘッダーの重要なもの

| ヘッダー | 用途 |
|---------|------|
| Content-Type | コンテンツのMIMEタイプ |
| Content-Encoding | 圧縮方式（gzip等） |
| Cache-Control | キャッシュ制御 |
| Set-Cookie | Cookie設定 |
| Content-Security-Policy | XSS対策 |
| Strict-Transport-Security | HTTPS強制 |

---

## Step 8: ブラウザレンダリング

### レンダリングプロセス

```mermaid
flowchart TB
    subgraph Parse["パース処理"]
        A[HTML受信] --> B[HTMLパース]
        B --> C[DOM構築]
        B --> D[CSS取得・パース]
        D --> E[CSSOM構築]
    end

    subgraph Build["構築処理"]
        C --> F[レンダーツリー構築]
        E --> F
        F --> G[レイアウト計算]
    end

    subgraph Render["描画処理"]
        G --> H[ペイント]
        H --> I[コンポジット]
        I --> J[画面表示]
    end

    subgraph JS["JavaScript"]
        K[JS取得・実行] -.->|DOM操作| C
        K -.->|スタイル変更| D
    end
```

### レンダリングの各ステップ

| ステップ | 説明 |
|---------|------|
| DOM構築 | HTMLをパースしてDocument Object Modelを構築 |
| CSSOM構築 | CSSをパースしてCSS Object Modelを構築 |
| レンダーツリー | DOMとCSSOMを結合して表示用ツリーを生成 |
| レイアウト | 各要素の位置とサイズを計算 |
| ペイント | 実際のピクセルを描画 |
| コンポジット | レイヤーを合成して最終画像を生成 |

### Critical Rendering Path

ページ表示を高速化するために重要なリソースを優先的に読み込みます。

```mermaid
flowchart LR
    A[HTML] --> B[CSS<br/>レンダーブロック]
    A --> C[JS<br/>パーサーブロック]
    B --> D[初回描画]
    C --> D

    style B fill:#ffcdd2
    style C fill:#fff9c4
```

| リソース | 特性 | 最適化 |
|---------|------|--------|
| CSS | レンダーブロック | head内で読み込み、クリティカルCSSをインライン化 |
| JavaScript | パーサーブロック | async/deferを使用、body末尾に配置 |
| 画像 | 非ブロック | lazy loading、適切なサイズ指定 |
| フォント | レンダーに影響 | font-display: swap、preload |

---

## まとめ

### 全体の処理時間の目安

| ステップ | 通常の所要時間 |
|---------|---------------|
| DNS解決 | 20-120ms |
| TCP接続 | 20-100ms（往復時間） |
| TLSハンドシェイク | 30-150ms（1-2 RTT） |
| HTTPリクエスト/レスポンス | 50-500ms |
| ブラウザレンダリング | 100-1000ms |
| **合計** | **約200ms〜2秒以上** |

### パフォーマンス最適化のポイント

```mermaid
flowchart TB
    A[高速化] --> B[DNS]
    A --> C[接続]
    A --> D[転送]
    A --> E[レンダリング]

    B --> B1[DNS prefetch<br/>キャッシュ活用]
    C --> C1[Keep-Alive<br/>HTTP/2,3]
    D --> D1[gzip圧縮<br/>CDN活用]
    E --> E1[クリティカルCSS<br/>遅延読み込み]
```

### 各ステップと関連技術

| ステップ | 関連技術・仕様 |
|---------|---------------|
| URL解析 | URL Standard (WHATWG) |
| DNS解決 | DNS, DNS over HTTPS (DoH) |
| TCP接続 | TCP, RFC 793 |
| TLS | TLS 1.3 (RFC 8446) |
| HTTP | HTTP/1.1, HTTP/2, HTTP/3 |
| レンダリング | DOM, CSSOM, Layout, Paint |

---

## 技術の歴史と意義

### なぜこれらの技術が必要だったのか

| 技術 | 登場年 | 解決した課題 |
|-----|-------|-------------|
| DNS | 1983年 | hosts.txtの手動管理の限界、インターネットの急速な成長 |
| TCP | 1974年 | 信頼性のないネットワーク上での確実なデータ転送 |
| HTTP | 1991年 | ハイパーテキストの転送プロトコルの標準化 |
| SSL/TLS | 1995年 | インターネット商取引における通信の安全性 |
| HTTP/2 | 2015年 | HTTP/1.1のHead-of-Line Blocking問題、多重化の必要性 |
| HTTP/3 | 2022年 | TCP上のHoL Blocking解消、モバイル環境での接続切り替え |

これらの技術は、インターネットの成長と共に、スケーラビリティ、セキュリティ、パフォーマンスの課題を解決するために発展してきました。
