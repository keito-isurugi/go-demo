package handler

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SecurityDemoHandler はセキュリティデモのメインハンドラ
type SecurityDemoHandler struct {
	DB *gorm.DB
}

// ==================== CORS Demo ====================

// CORSVulnerableHandler - 脆弱なCORS設定（全てのオリジンを許可）
func (h *SecurityDemoHandler) CORSVulnerableHandler(w http.ResponseWriter, r *http.Request) {
	// 脆弱性: すべてのオリジンからのアクセスを許可
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	response := map[string]interface{}{
		"message": "This endpoint has vulnerable CORS configuration (allows all origins)",
		"data":    "Sensitive user data",
		"status":  "vulnerable",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

// CORSSecureHandler - 安全なCORS設定（特定のオリジンのみ許可）
func (h *SecurityDemoHandler) CORSSecureHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	allowedOrigins := []string{"http://localhost:3000", "https://trusted-site.com"}

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
		"message": "This endpoint has secure CORS configuration (whitelist only)",
		"data":    "Sensitive user data",
		"status":  status,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

// ==================== CSRF Demo ====================

// CSRFトークンストア（本番環境ではRedisやセッションストアを使用）
var csrfTokens = make(map[string]time.Time)

// CSRFTokenHandler - CSRFトークンを発行
func (h *SecurityDemoHandler) CSRFTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := uuid.New().String()
	csrfTokens[token] = time.Now().Add(10 * time.Minute)

	response := map[string]interface{}{
		"message": "CSRF token generated",
		"token":   token,
		"expires": "10 minutes",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

// CSRFVulnerableHandler - CSRF保護なしのエンドポイント
func (h *SecurityDemoHandler) CSRFVulnerableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 脆弱性: CSRFトークンのチェックなし
	response := map[string]interface{}{
		"message": "Action executed without CSRF protection",
		"status":  "vulnerable",
		"action":  "Money transferred successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

// CSRFSecureHandler - CSRF保護ありのエンドポイント
func (h *SecurityDemoHandler) CSRFSecureHandler(w http.ResponseWriter, r *http.Request) {
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
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
		return
	}

	// トークンを使用済みにする
	delete(csrfTokens, token)

	response := map[string]interface{}{
		"message": "Action executed with CSRF protection",
		"status":  "secure",
		"action":  "Money transferred successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

// ==================== XSS Demo ====================

// XSSVulnerableHandler - XSS脆弱性のあるエンドポイント
func (h *SecurityDemoHandler) XSSVulnerableHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	// 脆弱性: ユーザー入力をエスケープせずに出力
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
	if _, err := w.Write([]byte(html)); err != nil {
		return
	}
}

// XSSSecureHandler - XSS対策済みのエンドポイント
func (h *SecurityDemoHandler) XSSSecureHandler(w http.ResponseWriter, r *http.Request) {
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
	if _, err := w.Write([]byte(htmlContent)); err != nil {
		return
	}
}

// ==================== SQL Injection Demo ====================

// SQLInjectionVulnerableHandler - SQLインジェクション脆弱性のあるエンドポイント
func (h *SecurityDemoHandler) SQLInjectionVulnerableHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		response := map[string]interface{}{
			"message": "Please provide username parameter",
			"example": "/api/security/sql-injection/vulnerable?username=admin",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
		return
	}

	// 脆弱性: SQLインジェクションが可能
	query := fmt.Sprintf("SELECT id, name, email FROM users WHERE name = '%s'", username)

	var results []map[string]interface{}
	rows, err := h.DB.Raw(query).Rows()
	if err != nil {
		response := map[string]interface{}{
			"message": "Query error (this might indicate SQL injection attempt was blocked by database)",
			"error":   err.Error(),
			"query":   query,
			"status":  "vulnerable",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			continue
		}
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
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

// SQLInjectionSecureHandler - SQLインジェクション対策済みのエンドポイント
func (h *SecurityDemoHandler) SQLInjectionSecureHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		response := map[string]interface{}{
			"message": "Please provide username parameter",
			"example": "/api/security/sql-injection/secure?username=admin",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
		return
	}

	// 入力値の検証
	if strings.Contains(username, "'") || strings.Contains(username, "\"") || strings.Contains(username, "--") {
		response := map[string]interface{}{
			"message": "Invalid characters detected in username",
			"status":  "secure - request blocked",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
		return
	}

	// 対策: プリペアドステートメント（パラメータ化クエリ）を使用
	var results []map[string]interface{}
	rows, err := h.DB.Raw("SELECT id, name, email FROM users WHERE name = ?", username).Rows()
	if err != nil {
		response := map[string]interface{}{
			"message": "Query error",
			"error":   err.Error(),
			"status":  "secure",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			continue
		}
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
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

// SecurityInfoHandler - セキュリティデモの情報を返す
func (h *SecurityDemoHandler) SecurityInfoHandler(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"message": "Security Demonstration APIs",
		"endpoints": map[string]interface{}{
			"CORS": map[string]string{
				"vulnerable": "/api/security/cors/vulnerable",
				"secure":     "/api/security/cors/secure",
			},
			"CSRF": map[string]string{
				"get_token":  "GET /api/security/csrf/token",
				"vulnerable": "POST /api/security/csrf/vulnerable",
				"secure":     "POST /api/security/csrf/secure (requires X-CSRF-Token header)",
			},
			"XSS": map[string]string{
				"vulnerable": "/api/security/xss/vulnerable?name=<script>alert('XSS')</script>",
				"secure":     "/api/security/xss/secure?name=<script>alert('XSS')</script>",
			},
			"SQL_Injection": map[string]string{
				"vulnerable": "/api/security/sql-injection/vulnerable?username=admin",
				"secure":     "/api/security/sql-injection/secure?username=admin",
				"test":       "Try: ?username=admin' OR '1'='1",
			},
		},
		"warning": "These endpoints are for educational purposes only. Do not use vulnerable patterns in production.",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		return
	}
}
