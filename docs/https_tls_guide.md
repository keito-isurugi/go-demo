# HTTPS/TLS 完全ガイド

Webセキュリティの基盤となるHTTPS/TLSの仕組みを解説します。

## 目次
1. [HTTPSとは](#httpsとは)
2. [TLSの概要](#tlsの概要)
3. [TLSハンドシェイク](#tlsハンドシェイク)
4. [証明書の仕組み](#証明書の仕組み)
5. [暗号化の種類](#暗号化の種類)
6. [TLS 1.3の改善点](#tls-13の改善点)
7. [Goでの実装例](#goでの実装例)

---

## HTTPSとは

**HTTPS**（HyperText Transfer Protocol Secure）は、HTTPにTLS/SSLによる暗号化を加えたプロトコルです。

```mermaid
flowchart LR
    subgraph HTTP["HTTP（平文）"]
        A[クライアント] -->|"GET /login<br/>password=secret"| B[サーバー]
    end

    subgraph HTTPS["HTTPS（暗号化）"]
        C[クライアント] -->|"🔒 暗号化されたデータ"| D[サーバー]
    end
```

### HTTPとHTTPSの違い

| 項目 | HTTP | HTTPS |
|------|------|-------|
| ポート | 80 | 443 |
| 暗号化 | なし | TLS/SSL |
| データの可視性 | 平文（盗聴可能） | 暗号化（保護） |
| 認証 | なし | サーバー証明書 |
| 完全性 | なし | 改ざん検知 |

### HTTPSが提供する3つの保護

```mermaid
flowchart TB
    HTTPS[HTTPS] --> E[暗号化<br/>Encryption]
    HTTPS --> A[認証<br/>Authentication]
    HTTPS --> I[完全性<br/>Integrity]

    E --> E1["通信内容を<br/>第三者から保護"]
    A --> A1["通信相手が<br/>本物であることを確認"]
    I --> I1["データが<br/>改ざんされていないことを確認"]
```

---

## TLSの概要

**TLS**（Transport Layer Security）は、通信を暗号化するためのプロトコルです。

### TLSのバージョン履歴

```mermaid
timeline
    title TLS/SSLの歴史
    1995 : SSL 2.0（脆弱性あり、非推奨）
    1996 : SSL 3.0（脆弱性あり、非推奨）
    1999 : TLS 1.0（非推奨）
    2006 : TLS 1.1（非推奨）
    2008 : TLS 1.2（現在も広く使用）
    2018 : TLS 1.3（最新、推奨）
```

### TLSのレイヤー構造

```mermaid
flowchart TB
    subgraph Application["アプリケーション層"]
        HTTP[HTTP/HTTPS]
    end

    subgraph TLS["TLS層"]
        direction TB
        HS[Handshake Protocol<br/>鍵交換・認証]
        REC[Record Protocol<br/>暗号化・復号]
        ALERT[Alert Protocol<br/>エラー通知]
        CCS[Change Cipher Spec<br/>暗号切替]
    end

    subgraph Transport["トランスポート層"]
        TCP[TCP]
    end

    Application --> TLS
    TLS --> Transport
```

---

## TLSハンドシェイク

TLSハンドシェイクは、暗号化通信を確立するための手続きです。

### TLS 1.2 ハンドシェイク（フルハンドシェイク）

```mermaid
sequenceDiagram
    participant C as クライアント
    participant S as サーバー

    Note over C,S: 1. ネゴシエーション
    C->>S: ClientHello<br/>（対応するTLSバージョン、暗号スイート一覧、乱数）
    S->>C: ServerHello<br/>（選択したTLSバージョン、暗号スイート、乱数）

    Note over C,S: 2. サーバー認証
    S->>C: Certificate<br/>（サーバー証明書）
    S->>C: ServerKeyExchange<br/>（鍵交換パラメータ）
    S->>C: ServerHelloDone

    Note over C,S: 3. 鍵交換
    C->>C: 証明書を検証
    C->>S: ClientKeyExchange<br/>（プリマスターシークレット）
    C->>S: ChangeCipherSpec
    C->>S: Finished（暗号化）

    Note over C,S: 4. 完了
    S->>C: ChangeCipherSpec
    S->>C: Finished（暗号化）

    Note over C,S: 🔒 暗号化通信開始
    C->>S: Application Data（暗号化）
    S->>C: Application Data（暗号化）
```

### ハンドシェイクの各ステップ詳細

#### Step 1: ClientHello

クライアントが送信する情報:

```
ClientHello {
    version: TLS 1.2
    random: 32バイトの乱数
    session_id: セッション再開用ID
    cipher_suites: [
        TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        ...
    ]
    compression_methods: [null]
    extensions: [
        server_name: "example.com",
        supported_groups: [x25519, secp256r1],
        ...
    ]
}
```

#### Step 2: ServerHello + Certificate

サーバーが送信する情報:

```
ServerHello {
    version: TLS 1.2
    random: 32バイトの乱数
    session_id: セッションID
    cipher_suite: TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
}

Certificate {
    certificate_list: [
        サーバー証明書,
        中間CA証明書,
        ...
    ]
}
```

#### Step 3: 鍵交換とマスターシークレット生成

```mermaid
flowchart TB
    subgraph Client["クライアント側"]
        CR[Client Random]
        PMS1[Pre-Master Secret]
    end

    subgraph Server["サーバー側"]
        SR[Server Random]
        PMS2[Pre-Master Secret]
    end

    CR --> MS1[Master Secret]
    SR --> MS1
    PMS1 --> MS1

    CR --> MS2[Master Secret]
    SR --> MS2
    PMS2 --> MS2

    MS1 --> KEY1[Session Keys]
    MS2 --> KEY2[Session Keys]

    KEY1 -.->|"同じ鍵"| KEY2
```

**Master Secret の計算:**
```
master_secret = PRF(pre_master_secret,
                    "master secret",
                    ClientHello.random + ServerHello.random)
```

---

## 証明書の仕組み

### 証明書チェーン

```mermaid
flowchart TB
    subgraph Root["ルート認証局（Root CA）"]
        RCA[ルート証明書<br/>自己署名]
    end

    subgraph Intermediate["中間認証局（Intermediate CA）"]
        ICA[中間証明書<br/>ルートCAが署名]
    end

    subgraph Server["サーバー"]
        SC[サーバー証明書<br/>中間CAが署名]
    end

    RCA -->|署名| ICA
    ICA -->|署名| SC

    subgraph Browser["ブラウザ/OS"]
        TS[トラストストア<br/>ルート証明書一覧]
    end

    TS -.->|"検証"| RCA
```

### 証明書の検証プロセス

```mermaid
flowchart TD
    START[証明書を受信] --> CHECK1{有効期限内？}
    CHECK1 -->|No| FAIL1[❌ 期限切れ]
    CHECK1 -->|Yes| CHECK2{ドメイン名が一致？}

    CHECK2 -->|No| FAIL2[❌ ドメイン不一致]
    CHECK2 -->|Yes| CHECK3{署名は有効？}

    CHECK3 -->|No| FAIL3[❌ 署名無効]
    CHECK3 -->|Yes| CHECK4{発行者は信頼できる？}

    CHECK4 -->|No| CHECK5{中間証明書あり？}
    CHECK5 -->|No| FAIL4[❌ 信頼されていない]
    CHECK5 -->|Yes| CHECK3

    CHECK4 -->|Yes| CHECK6{失効していない？}

    CHECK6 -->|No| FAIL5[❌ 失効済み]
    CHECK6 -->|Yes| SUCCESS[✅ 検証成功]
```

### X.509証明書の構造

```
Certificate {
    Version: 3
    Serial Number: 123456789...
    Signature Algorithm: sha256WithRSAEncryption
    Issuer: CN=Example CA, O=Example Inc
    Validity:
        Not Before: Jan 1 00:00:00 2024 GMT
        Not After:  Jan 1 00:00:00 2025 GMT
    Subject: CN=www.example.com, O=Example Inc
    Subject Public Key Info:
        Algorithm: rsaEncryption
        Public Key: (2048 bit)
    Extensions:
        Subject Alternative Name:
            DNS: www.example.com
            DNS: example.com
        Key Usage: Digital Signature, Key Encipherment
        Extended Key Usage: TLS Web Server Authentication
}
```

---

## 暗号化の種類

### 対称鍵暗号と公開鍵暗号

```mermaid
flowchart TB
    subgraph Symmetric["対称鍵暗号（共通鍵暗号）"]
        direction LR
        SK[共通鍵 🔑]
        P1[平文] --> E1[暗号化] --> C1[暗号文]
        C1 --> D1[復号] --> P2[平文]
        SK --> E1
        SK --> D1
    end

    subgraph Asymmetric["公開鍵暗号（非対称鍵暗号）"]
        direction LR
        PK[公開鍵 🔓]
        PRK[秘密鍵 🔐]
        P3[平文] --> E2[暗号化] --> C2[暗号文]
        C2 --> D2[復号] --> P4[平文]
        PK --> E2
        PRK --> D2
    end
```

### TLSで使用される暗号

| 用途 | アルゴリズム | 例 |
|------|-------------|-----|
| 鍵交換 | DH, ECDH | ECDHE (Elliptic Curve Diffie-Hellman Ephemeral) |
| 認証 | RSA, ECDSA | RSA-2048, ECDSA P-256 |
| 暗号化 | AES, ChaCha20 | AES-256-GCM, ChaCha20-Poly1305 |
| ハッシュ | SHA | SHA-256, SHA-384 |

### 暗号スイートの読み方

```
TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
│    │     │        │   │   │    │
│    │     │        │   │   │    └── ハッシュアルゴリズム
│    │     │        │   │   └────── 認証タグ（AEAD）
│    │     │        │   └────────── 鍵長
│    │     │        └────────────── 暗号化アルゴリズム
│    │     └─────────────────────── 認証アルゴリズム
│    └───────────────────────────── 鍵交換アルゴリズム
└────────────────────────────────── プロトコル
```

### Forward Secrecy（前方秘匿性）

```mermaid
sequenceDiagram
    participant C as クライアント
    participant S as サーバー
    participant A as 攻撃者

    Note over C,S: ECDHE鍵交換
    C->>S: 一時的な公開鍵（セッションごとに生成）
    S->>C: 一時的な公開鍵（セッションごとに生成）

    Note over C,S: 🔒 暗号化通信
    C->>S: 暗号化データ

    Note over A: 通信を記録
    A-->>A: 暗号化データを保存

    Note over A: 後日、サーバーの秘密鍵を入手
    A-->>A: ❌ 過去の通信は復号できない
    Note over A: セッション鍵は一時的で<br/>サーバー秘密鍵から導出不可
```

---

## TLS 1.3の改善点

### ハンドシェイクの高速化

```mermaid
flowchart LR
    subgraph TLS12["TLS 1.2"]
        direction TB
        R12_1[RTT 1: ClientHello/ServerHello]
        R12_2[RTT 2: 証明書/鍵交換]
        R12_3[RTT 3: Finished]
        R12_1 --> R12_2 --> R12_3
    end

    subgraph TLS13["TLS 1.3"]
        direction TB
        R13_1[RTT 1: ClientHello/ServerHello<br/>+ 鍵交換 + Finished]
        R13_2[Application Data]
        R13_1 --> R13_2
    end
```

### TLS 1.3 ハンドシェイク

```mermaid
sequenceDiagram
    participant C as クライアント
    participant S as サーバー

    Note over C,S: 1-RTT ハンドシェイク
    C->>S: ClientHello<br/>+ key_share（鍵交換パラメータ）<br/>+ supported_versions

    S->>C: ServerHello<br/>+ key_share
    S->>C: {EncryptedExtensions}
    S->>C: {Certificate}
    S->>C: {CertificateVerify}
    S->>C: {Finished}

    Note over C,S: この時点で暗号化開始
    C->>S: {Finished}
    C->>S: [Application Data]
    S->>C: [Application Data]
```

### 0-RTT（ゼロラウンドトリップ）

```mermaid
sequenceDiagram
    participant C as クライアント
    participant S as サーバー

    Note over C,S: 初回接続
    C->>S: 通常のハンドシェイク
    S->>C: NewSessionTicket（再開用チケット）

    Note over C,S: 2回目以降（0-RTT）
    C->>S: ClientHello<br/>+ early_data<br/>+ [Application Data]（暗号化済み）

    Note over C,S: 最初のリクエストと同時に<br/>データを送信可能！
    S->>C: ServerHello...
    S->>C: [Application Data]
```

### TLS 1.2 vs TLS 1.3 比較

| 項目 | TLS 1.2 | TLS 1.3 |
|------|---------|---------|
| ハンドシェイクRTT | 2 RTT | 1 RTT (0-RTTも可能) |
| 暗号スイート | 多数（レガシー含む） | 5つのみ（安全なもの） |
| 鍵交換 | RSA, DH, ECDH | ECDHE, DHEのみ |
| Forward Secrecy | オプション | 必須 |
| 暗号化開始 | Finishedメッセージ後 | ServerHello直後 |

### TLS 1.3で削除された機能

```mermaid
flowchart TB
    subgraph Removed["❌ 削除された機能"]
        RSA_KE[RSA鍵交換]
        DES[DES/3DES]
        RC4[RC4]
        MD5[MD5]
        COMP[圧縮]
        RENEG[再ネゴシエーション]
    end

    subgraph Reason["削除理由"]
        RSA_KE --> R1[Forward Secrecyなし]
        DES --> R2[暗号強度不足]
        RC4 --> R2
        MD5 --> R3[衝突攻撃に脆弱]
        COMP --> R4[CRIME攻撃]
        RENEG --> R5[複雑性・脆弱性]
    end
```

---

## Goでの実装例

### HTTPSサーバー（基本）

```go
package main

import (
    "crypto/tls"
    "log"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, HTTPS!"))
    })

    server := &http.Server{
        Addr:    ":443",
        Handler: mux,
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
            CipherSuites: []uint16{
                tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
            },
        },
    }

    log.Println("Starting HTTPS server on :443")
    log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}
```

### TLS 1.3専用サーバー

```go
package main

import (
    "crypto/tls"
    "log"
    "net/http"
)

func main() {
    tlsConfig := &tls.Config{
        MinVersion: tls.VersionTLS13, // TLS 1.3以上を強制
        // TLS 1.3ではCipherSuitesは自動選択される
    }

    server := &http.Server{
        Addr:      ":443",
        TLSConfig: tlsConfig,
    }

    http.HandleFunc("/", handler)
    log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}
```

### 相互TLS認証（mTLS）

```mermaid
sequenceDiagram
    participant C as クライアント
    participant S as サーバー

    Note over C,S: 通常のTLS
    C->>S: ClientHello
    S->>C: ServerHello + 証明書
    C->>C: サーバー証明書を検証 ✅

    Note over C,S: 相互TLS（mTLS）
    S->>C: CertificateRequest
    C->>S: クライアント証明書
    S->>S: クライアント証明書を検証 ✅

    Note over C,S: 双方が認証済み 🔒
```

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "log"
    "net/http"
    "os"
)

func main() {
    // クライアント証明書を検証するためのCAプール
    caCert, err := os.ReadFile("ca.crt")
    if err != nil {
        log.Fatal(err)
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    tlsConfig := &tls.Config{
        ClientCAs:  caCertPool,
        ClientAuth: tls.RequireAndVerifyClientCert, // クライアント証明書を必須に
        MinVersion: tls.VersionTLS12,
    }

    server := &http.Server{
        Addr:      ":443",
        TLSConfig: tlsConfig,
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // クライアント証明書の情報を取得
        if len(r.TLS.PeerCertificates) > 0 {
            cert := r.TLS.PeerCertificates[0]
            log.Printf("Client: %s", cert.Subject.CommonName)
        }
        w.Write([]byte("Hello, mTLS!"))
    })

    log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}
```

### HTTPSクライアント

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

func main() {
    // カスタムCA証明書を読み込む場合
    caCert, err := os.ReadFile("ca.crt")
    if err != nil {
        log.Fatal(err)
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                RootCAs:    caCertPool,
                MinVersion: tls.VersionTLS12,
            },
        },
    }

    resp, err := client.Get("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### 自己署名証明書の生成

```go
package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "log"
    "math/big"
    "net"
    "os"
    "time"
)

func main() {
    // 秘密鍵を生成
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        log.Fatal(err)
    }

    // 証明書テンプレート
    template := x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject: pkix.Name{
            Organization: []string{"Example Inc"},
            CommonName:   "localhost",
        },
        NotBefore:             time.Now(),
        NotAfter:              time.Now().Add(365 * 24 * time.Hour),
        KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        BasicConstraintsValid: true,
        DNSNames:              []string{"localhost"},
        IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
    }

    // 証明書を生成（自己署名）
    certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
    if err != nil {
        log.Fatal(err)
    }

    // 証明書をPEM形式で保存
    certFile, _ := os.Create("server.crt")
    pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
    certFile.Close()

    // 秘密鍵をPEM形式で保存
    keyFile, _ := os.Create("server.key")
    keyBytes, _ := x509.MarshalECPrivateKey(privateKey)
    pem.Encode(keyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
    keyFile.Close()

    log.Println("Generated server.crt and server.key")
}
```

---

## よくある攻撃と対策

### 中間者攻撃（MITM）

```mermaid
flowchart LR
    subgraph NoTLS["HTTPの場合"]
        C1[クライアント] --> A1[攻撃者] --> S1[サーバー]
        A1 -->|盗聴・改ざん可能| A1
    end

    subgraph WithTLS["HTTPSの場合"]
        C2[クライアント] -->|🔒| S2[サーバー]
        A2[攻撃者]
        A2 -.->|❌ 復号不可| C2
    end
```

### ダウングレード攻撃

**攻撃:** 古い脆弱なプロトコルへの強制ダウングレード

**対策:**
```go
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12, // TLS 1.2未満を拒否
}
```

### 証明書の検証をスキップしない

```go
// ❌ 絶対にやってはいけない（開発時のみ）
tlsConfig := &tls.Config{
    InsecureSkipVerify: true, // 証明書の検証をスキップ
}

// ✅ 正しい方法
tlsConfig := &tls.Config{
    RootCAs:    caCertPool, // 信頼するCAを指定
    MinVersion: tls.VersionTLS12,
}
```

---

## ベストプラクティス

### サーバー設定のチェックリスト

- [ ] TLS 1.2以上を使用（TLS 1.3推奨）
- [ ] 強力な暗号スイートのみを有効化
- [ ] Forward Secrecyを有効化（ECDHE）
- [ ] 有効な証明書を使用（期限切れに注意）
- [ ] HSTSヘッダーを設定
- [ ] OCSP Staplingを有効化

### 推奨TLS設定（Go）

```go
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
    CurvePreferences: []tls.CurveID{
        tls.X25519,
        tls.CurveP256,
    },
    CipherSuites: []uint16{
        // TLS 1.3の暗号スイート（自動選択）
        // TLS 1.2の暗号スイート
        tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
        tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
    },
}
```

---

## 参考資料

- [RFC 8446 - TLS 1.3](https://tools.ietf.org/html/rfc8446)
- [Go crypto/tls パッケージ](https://pkg.go.dev/crypto/tls)
- [Mozilla SSL Configuration Generator](https://ssl-config.mozilla.org/)
- [SSL Labs Server Test](https://www.ssllabs.com/ssltest/)
