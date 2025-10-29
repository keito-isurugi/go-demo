## 概要

- Goで実装するWebセキュリティの脆弱性と対策
- CORS、CSRF、XSS、SQLインジェクションの4つの主要な脆弱性を解説
- 脆弱な実装と安全な実装を比較し、攻撃手法と防御方法を理解する

### 特徴

- **CORS (Cross-Origin Resource Sharing)**: オリジン間リソース共有の適切な設定
- **CSRF (Cross-Site Request Forgery)**: クロスサイトリクエストフォージェリ対策
- **XSS (Cross-Site Scripting)**: クロスサイトスクリプティング対策
- **SQL Injection**: SQLインジェクション対策
- **実践的な例**: 脆弱なコードと安全なコードを並べて比較

## Webセキュリティの重要性

### なぜセキュリティが重要か

1. **ユーザー保護**: 個人情報やデータの漏洩を防止
2. **サービスの信頼性**: セキュリティインシデントは信頼を失う
3. **法的責任**: GDPRなど法規制への対応
4. **経済的損失**: 攻撃による損害は甚大
5. **ブランドイメージ**: セキュリティ問題は企業イメージを損なう

### OWASP Top 10

OWASP（Open Web Application Security Project）が定める主要な脆弱性：

1. **Broken Access Control** - アクセス制御の不備
2. **Cryptographic Failures** - 暗号化の失敗
3. **Injection** - インジェクション攻撃（SQL、XSS等）
4. **Insecure Design** - 不安全な設計
5. **Security Misconfiguration** - セキュリティ設定ミス
6. **Vulnerable Components** - 脆弱なコンポーネント
7. **Authentication Failures** - 認証の不備
8. **Data Integrity Failures** - データ整合性の失敗
9. **Logging Failures** - ログ・監視の不備
10. **SSRF** - サーバーサイドリクエストフォージェリ

今回解説する4つの脆弱性は、これらの中でも特に頻繁に発生する問題。

## 1. CORS (Cross-Origin Resource Sharing)

### 概要

CORS（Cross-Origin Resource Sharing）は、異なるオリジン間でのリソース共有を制御するメカニズム。

**オリジンとは**:
```
https://example.com:443
 ↓       ↓          ↓
スキーム  ドメイン    ポート

異なるオリジンの例:
- http://example.com (スキームが違う)
- https://api.example.com (ドメインが違う)
- https://example.com:8080 (ポートが違う)
```

### 脆弱な実装

```go
// 脆弱性: すべてのオリジンからのアクセスを許可
func CORSVulnerableHandler(w http.ResponseWriter, r *http.Request) {
	// 危険: ワイルドカード使用
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	response := map[string]interface{}{
		"message": "This endpoint has vulnerable CORS configuration",
		"data":    "Sensitive user data",
		"status":  "vulnerable",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

**問題点**:
- `Access-Control-Allow-Origin: *` は全てのオリジンを許可
- 悪意のあるサイトからAPIを呼び出し可能
- ユーザーの認証情報を使って不正リクエスト

### 安全な実装

```go
// 安全: 特定のオリジンのみ許可（ホワイトリスト方式）
func CORSSecureHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	allowedOrigins := []string{
		"http://localhost:3000",
		"https://trusted-site.com",
	}

	// 許可されたオリジンかチェック
	isAllowed := false
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			isAllowed = true
			w.Header().Set("Access-Control-Allow-Origin", origin)
			break
		}
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var status string
	if isAllowed {
		status = "secure - access granted"
	} else {
		status = "secure - access denied"
	}

	response := map[string]interface{}{
		"message": "This endpoint has secure CORS configuration",
		"data":    "Sensitive user data",
		"status":  status,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

**対策のポイント**:
1. **ホワイトリスト方式**: 信頼できるオリジンのみ許可
2. **動的な設定**: リクエストのOriginヘッダーを検証
3. **認証情報の扱い**: Cookieを含む場合は`*`使用不可
4. **環境変数で管理**: 本番環境と開発環境で設定を分ける

### CORSミドルウェアの実装例

```go
// 再利用可能なCORSミドルウェア
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// ホワイトリストチェック
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					break
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")

			// Preflightリクエスト
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// 使用例
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", dataHandler)

	allowedOrigins := []string{
		"https://myapp.com",
		"https://www.myapp.com",
	}
	handler := CORSMiddleware(allowedOrigins)(mux)

	http.ListenAndServe(":8080", handler)
}
```

## 2. CSRF (Cross-Site Request Forgery)

### 概要

CSRF（クロスサイトリクエストフォージェリ）は、ユーザーが意図しないリクエストを強制的に実行させる攻撃。

**攻撃シナリオ**:
```
1. ユーザーが銀行サイトにログイン（認証済み）
2. 別のタブで悪意のあるサイトを訪問
3. 悪意のあるサイトが隠しフォームで銀行APIへPOST
4. ブラウザは自動的にCookieを送信
5. 銀行は正規のリクエストと判断して処理
```

### トークン発行

```go
// CSRFトークンの発行
var csrfTokens = make(map[string]time.Time)

func CSRFTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := uuid.New().String()
	csrfTokens[token] = time.Now().Add(10 * time.Minute)

	response := map[string]interface{}{
		"message": "CSRF token generated",
		"token":   token,
		"expires": "10 minutes",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

### 脆弱な実装

```go
// 脆弱性: CSRFトークンのチェックなし
func CSRFVulnerableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 危険: トークン検証なしで重要な処理を実行
	response := map[string]interface{}{
		"message": "Action executed without CSRF protection",
		"status":  "vulnerable",
		"action":  "Money transferred successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

**問題点**:
- 認証済みセッションがあれば誰でも実行可能
- 悪意のあるサイトから攻撃可能
- ユーザーの意図しない操作が実行される

### 安全な実装

```go
// 安全: CSRFトークンを検証
func CSRFSecureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("X-CSRF-Token")

	// トークンの検証
	expiry, exists := csrfTokens[token]
	if !exists || time.Now().After(expiry) {
		response := map[string]interface{}{
			"message": "Invalid or expired CSRF token",
			"status":  "secure - request blocked",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
		return
	}

	// トークンを使用済みにする（ワンタイムトークン）
	delete(csrfTokens, token)

	response := map[string]interface{}{
		"message": "Action executed with CSRF protection",
		"status":  "secure",
		"action":  "Money transferred successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

**対策のポイント**:
1. **トークン検証**: すべての状態変更操作でトークンをチェック
2. **有効期限**: トークンに期限を設定（10分など）
3. **ワンタイムトークン**: 使用後は即座に無効化
4. **セキュアなトークン**: UUIDなど推測不可能な値を使用

### CSRFミドルウェアの実装例

```go
type CSRFMiddleware struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

func NewCSRFMiddleware() *CSRFMiddleware {
	csrf := &CSRFMiddleware{
		tokens: make(map[string]time.Time),
	}
	// 定期的に期限切れトークンを削除
	go csrf.cleanup()
	return csrf
}

func (c *CSRFMiddleware) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for token, expiry := range c.tokens {
			if now.After(expiry) {
				delete(c.tokens, token)
			}
		}
		c.mu.Unlock()
	}
}

func (c *CSRFMiddleware) GenerateToken() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	token := uuid.New().String()
	c.tokens[token] = time.Now().Add(10 * time.Minute)
	return token
}

func (c *CSRFMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// GET, HEAD, OPTIONS はCSRF対象外
		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("X-CSRF-Token")

		c.mu.Lock()
		expiry, exists := c.tokens[token]
		if exists && time.Now().Before(expiry) {
			delete(c.tokens, token) // ワンタイムトークン
			c.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}
		c.mu.Unlock()

		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
	})
}
```

## 3. XSS (Cross-Site Scripting)

### 概要

XSS（クロスサイトスクリプティング）は、Webページに悪意のあるスクリプトを挿入して実行させる攻撃。

**XSSの種類**:
1. **Reflected XSS**: URLパラメータからスクリプトを反射
2. **Stored XSS**: データベースに保存されたスクリプト
3. **DOM-based XSS**: クライアント側でDOMを操作

### 脆弱な実装

```go
// 脆弱性: ユーザー入力をエスケープせずに出力
func XSSVulnerableHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	// 危険: ユーザー入力をそのままHTML出力
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>XSS Vulnerable Page</title>
</head>
<body>
    <h1>Welcome, %s!</h1>
    <p>This page is vulnerable to XSS attacks.</p>
    <p>Try: <code>?name=&lt;script&gt;alert('XSS')&lt;/script&gt;</code></p>
</body>
</html>
`, name)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
```

**攻撃例**:
```
URL: /page?name=<script>alert('XSS')</script>
URL: /page?name=<img src=x onerror="alert('XSS')">
URL: /page?name=<iframe src="javascript:alert('XSS')">
```

**影響**:
- Cookieの盗取（セッションハイジャック）
- キーロガーの埋め込み
- フィッシングページの表示
- ユーザー操作の乗っ取り

### 安全な実装

```go
// 安全: HTMLエスケープ処理
func XSSSecureHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	// 対策: HTMLエスケープ処理
	safeName := html.EscapeString(name)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>XSS Secure Page</title>
</head>
<body>
    <h1>Welcome, %s!</h1>
    <p>This page is protected against XSS attacks.</p>
    <p>Try: <code>?name=&lt;script&gt;alert('XSS')&lt;/script&gt;</code></p>
    <p>The script tag will be displayed as text, not executed.</p>
</body>
</html>
`, safeName)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlContent))
}
```

**エスケープの結果**:
```
入力:  <script>alert('XSS')</script>
出力:  &lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;
表示:  <script>alert('XSS')</script> （テキストとして表示）
```

### 対策のポイント

1. **入力検証**: ユーザー入力を厳格にバリデーション
2. **出力エスケープ**: HTML、JavaScript、URL、CSSそれぞれに適したエスケープ
3. **Content-Security-Policy**: CSPヘッダーでスクリプト実行を制限
4. **テンプレートエンジン**: 自動エスケープ機能を持つテンプレートを使用

### テンプレートエンジンの使用例

```go
import "html/template"

// html/template は自動的にエスケープ
func XSSSecureTemplateHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	tmpl := template.Must(template.New("page").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>XSS Secure Page</title>
</head>
<body>
    <h1>Welcome, {{.Name}}!</h1>
    <p>This page automatically escapes HTML.</p>
</body>
</html>
`))

	data := struct {
		Name string
	}{
		Name: name, // 自動的にエスケープされる
	}

	tmpl.Execute(w, data)
}
```

### Content-Security-Policy の設定

```go
func SecureHandler(w http.ResponseWriter, r *http.Request) {
	// CSPヘッダーでスクリプトの実行元を制限
	w.Header().Set("Content-Security-Policy",
		"default-src 'self'; "+
		"script-src 'self' https://trusted-cdn.com; "+
		"style-src 'self' 'unsafe-inline'; "+
		"img-src 'self' data: https:;")

	// その他のセキュリティヘッダー
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	// コンテンツ出力
	// ...
}
```

## 4. SQL Injection

### 概要

SQLインジェクションは、ユーザー入力をSQL文に直接埋め込むことで、意図しないSQL文を実行させる攻撃。

**攻撃の影響**:
- データベース全体の閲覧
- データの改ざん・削除
- 管理者権限の取得
- システム乗っ取り

### 脆弱な実装

```go
// 脆弱性: SQLインジェクションが可能
func SQLInjectionVulnerableHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		response := map[string]interface{}{
			"message": "Please provide username parameter",
			"example": "/api/security/sql-injection/vulnerable?username=admin",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 危険: ユーザー入力を直接SQL文に埋め込み
	query := fmt.Sprintf("SELECT id, name, email FROM users WHERE name = '%s'", username)

	var results []map[string]interface{}
	rows, err := db.Raw(query).Rows()
	if err != nil {
		response := map[string]interface{}{
			"message": "Query error",
			"error":   err.Error(),
			"query":   query,
			"status":  "vulnerable",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email string
		rows.Scan(&id, &name, &email)
		results = append(results, map[string]interface{}{
			"id":    id,
			"name":  name,
			"email": email,
		})
	}

	response := map[string]interface{}{
		"message": "Query executed without parameterization",
		"query":   query,
		"results": results,
		"status":  "vulnerable",
		"warning": "Try: ?username=admin' OR '1'='1",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

**攻撃例**:
```sql
-- 正常なクエリ
SELECT id, name, email FROM users WHERE name = 'admin'

-- 攻撃パターン1: すべてのユーザーを取得
?username=admin' OR '1'='1
→ SELECT id, name, email FROM users WHERE name = 'admin' OR '1'='1'

-- 攻撃パターン2: コメントアウトで条件を無効化
?username=admin'--
→ SELECT id, name, email FROM users WHERE name = 'admin'--'

-- 攻撃パターン3: UNION でパスワードテーブルを結合
?username=admin' UNION SELECT id, password, salt FROM passwords--
→ SELECT id, name, email FROM users WHERE name = 'admin'
   UNION SELECT id, password, salt FROM passwords--'

-- 攻撃パターン4: データベース削除
?username=admin'; DROP TABLE users;--
→ SELECT id, name, email FROM users WHERE name = 'admin';
   DROP TABLE users;--'
```

### 安全な実装

```go
// 安全: プリペアドステートメントを使用
func SQLInjectionSecureHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		response := map[string]interface{}{
			"message": "Please provide username parameter",
			"example": "/api/security/sql-injection/secure?username=admin",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 入力値の検証（追加の防御層）
	if strings.Contains(username, "'") ||
	   strings.Contains(username, "\"") ||
	   strings.Contains(username, "--") {
		response := map[string]interface{}{
			"message": "Invalid characters detected in username",
			"status":  "secure - request blocked",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 対策: プリペアドステートメント（パラメータ化クエリ）
	var results []map[string]interface{}
	rows, err := db.Raw("SELECT id, name, email FROM users WHERE name = ?", username).Rows()
	if err != nil {
		response := map[string]interface{}{
			"message": "Query error",
			"error":   err.Error(),
			"status":  "secure",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email string
		rows.Scan(&id, &name, &email)
		results = append(results, map[string]interface{}{
			"id":    id,
			"name":  name,
			"email": email,
		})
	}

	response := map[string]interface{}{
		"message": "Query executed with parameterization",
		"results": results,
		"status":  "secure",
		"note":    "SQL injection attempts are safely escaped",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

### プリペアドステートメントの仕組み

```
プレースホルダー使用:
SQL: SELECT * FROM users WHERE name = ?
値:  admin' OR '1'='1

データベースの処理:
1. SQLを解析してクエリプランを作成
2. パラメータを「データ」として扱う
3. SQLコマンドとして解釈しない

結果: "admin' OR '1'='1" という名前のユーザーを検索
　　　（攻撃は失敗）
```

### GORMを使った安全な実装

```go
type User struct {
	ID    uint
	Name  string
	Email string
}

// GORMのクエリビルダーは自動的にエスケープ
func GetUserByName(db *gorm.DB, name string) ([]User, error) {
	var users []User

	// 安全: GORMは自動的にパラメータ化
	result := db.Where("name = ?", name).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// 動的な検索条件も安全に処理
func SearchUsers(db *gorm.DB, filters map[string]string) ([]User, error) {
	var users []User
	query := db.Model(&User{})

	for field, value := range filters {
		// ホワイトリストで許可されたフィールドのみ
		switch field {
		case "name":
			query = query.Where("name LIKE ?", "%"+value+"%")
		case "email":
			query = query.Where("email = ?", value)
		}
	}

	result := query.Find(&users)
	return users, result.Error
}
```

### 対策のポイント

1. **プリペアドステートメント**: 常にパラメータ化クエリを使用
2. **入力検証**: ホワイトリスト方式で許可する文字を制限
3. **最小権限の原則**: データベースユーザーに最小限の権限のみ付与
4. **ORM使用**: GORMなどのORMは自動的にエスケープ
5. **エラーメッセージ**: 詳細なSQL情報を返さない

## 総合的なセキュリティ対策

### セキュリティヘッダーの設定

```go
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// XSS対策
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// クリックジャッキング対策
		w.Header().Set("X-Frame-Options", "DENY")

		// Content Security Policy
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
			"style-src 'self' 'unsafe-inline';")

		// HTTPS強制（本番環境）
		w.Header().Set("Strict-Transport-Security",
			"max-age=31536000; includeSubDomains")

		// リファラー制御
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// 機能制御
		w.Header().Set("Permissions-Policy",
			"geolocation=(), microphone=(), camera=()")

		next.ServeHTTP(w, r)
	})
}
```

### 入力検証のベストプラクティス

```go
import "github.com/go-playground/validator/v10"

type UserInput struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,gte=0,lte=150"`
	Website  string `json:"website" validate:"omitempty,url"`
}

func ValidateInput(w http.ResponseWriter, r *http.Request) {
	var input UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		response := map[string]interface{}{
			"error":   "Validation failed",
			"details": err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// バリデーション成功
	// ...
}
```

### レート制限との組み合わせ

```go
func SecureAPIHandler(rateLimiter *RateLimiter, csrfMiddleware *CSRFMiddleware) http.Handler {
	mux := http.NewServeMux()

	// エンドポイント登録
	mux.HandleFunc("/api/data", dataHandler)
	mux.HandleFunc("/api/user", userHandler)

	// ミドルウェアチェーン
	handler := SecurityHeadersMiddleware(
		rateLimiter.Middleware(
			csrfMiddleware.Middleware(mux)))

	return handler
}
```

### ログと監視

```go
import "log"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 疑わしいパターンの検出
		username := r.URL.Query().Get("username")
		if strings.Contains(username, "'") ||
		   strings.Contains(username, "OR") ||
		   strings.Contains(username, "--") {
			log.Printf("SECURITY WARNING: Potential SQL injection attempt from %s: %s",
				r.RemoteAddr, username)
		}

		// XSS攻撃の検出
		if strings.Contains(r.URL.String(), "<script>") ||
		   strings.Contains(r.URL.String(), "javascript:") {
			log.Printf("SECURITY WARNING: Potential XSS attempt from %s: %s",
				r.RemoteAddr, r.URL.String())
		}

		next.ServeHTTP(w, r)

		log.Printf("%s %s %s %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start))
	})
}
```

## 使いどころ

### 実装すべき場面

1. **公開API**: すべてのセキュリティ対策が必須
2. **ユーザー認証**: CSRF、セッション管理が重要
3. **データベース操作**: SQLインジェクション対策は必須
4. **ユーザー生成コンテンツ**: XSS対策が必須
5. **クロスオリジンAPI**: CORS設定を適切に

### セキュリティチェックリスト

- [ ] すべてのユーザー入力を検証
- [ ] プリペアドステートメントを使用
- [ ] HTMLエスケープを実施
- [ ] CSRFトークンを実装
- [ ] CORS設定を適切に
- [ ] セキュリティヘッダーを設定
- [ ] HTTPSを使用
- [ ] レート制限を実装
- [ ] ログと監視を実装
- [ ] 定期的なセキュリティ監査

## まとめ

### 重要なポイント

| 脆弱性 | 主な対策 | 難易度 |
| --- | --- | --- |
| **CORS** | ホワイトリスト方式 | 低 |
| **CSRF** | トークン検証 | 中 |
| **XSS** | エスケープ処理 | 中 |
| **SQL Injection** | プリペアドステートメント | 低 |

### セキュリティの原則

1. **多層防御**: 複数の対策を組み合わせる
2. **最小権限**: 必要最小限の権限のみ付与
3. **ホワイトリスト**: 許可するものを明示的に定義
4. **失敗時の安全**: エラー時は安全側に倒す
5. **継続的な改善**: 定期的な監査とアップデート

### メリット

- **ユーザー保護**: データ漏洩や不正アクセスを防止
- **信頼性向上**: セキュリティ対策は信頼の証
- **法的コンプライアンス**: 規制への対応
- **コスト削減**: インシデント対応コストを削減

### 実装の優先順位

1. **最優先**: SQLインジェクション対策（プリペアドステートメント）
2. **高優先**: XSS対策（エスケープ処理）
3. **高優先**: HTTPS化、セキュリティヘッダー
4. **中優先**: CSRF対策（状態変更操作）
5. **中優先**: CORS設定（公開API）
6. **推奨**: レート制限、ログ監視

Webセキュリティは一度実装すれば終わりではなく、継続的な改善が必要。新しい攻撃手法が日々発見されるため、最新のセキュリティ情報をキャッチアップし、定期的な監査を実施することが重要。Goの標準ライブラリは基本的なセキュリティ機能を提供しているため、適切に使用すれば安全なアプリケーションを構築できる。
