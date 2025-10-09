package algorithm

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"math/rand"
	"sync"
	"time"
)

// ============================================================
// Server - バックエンドサーバーの表現
// ============================================================

// Server はバックエンドサーバーを表す
type Server struct {
	ID       string
	Name     string
	Load     int // 現在の負荷（接続数）
	mu       sync.Mutex
}

// NewServer は新しいサーバーを作成
func NewServer(id, name string) *Server {
	return &Server{
		ID:   id,
		Name: name,
		Load: 0,
	}
}

// AddLoad はサーバーの負荷を増加
func (s *Server) AddLoad() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Load++
}

// GetLoad は現在の負荷を取得
func (s *Server) GetLoad() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Load
}

// ============================================================
// 1. Cookie-Based Sticky Session
// ============================================================

// CookieBasedStickySession はCookieベースのSticky Sessionを実装
type CookieBasedStickySession struct {
	servers   []*Server
	sessions  map[string]*Server // sessionID -> server
	mu        sync.RWMutex
}

// NewCookieBasedStickySession は新しいCookie-Based Sticky Sessionを作成
func NewCookieBasedStickySession(servers []*Server) *CookieBasedStickySession {
	return &CookieBasedStickySession{
		servers:  servers,
		sessions: make(map[string]*Server),
	}
}

// GetServer はセッションIDに基づいてサーバーを取得
func (lb *CookieBasedStickySession) GetServer(sessionID string) *Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// 既存のセッションがあればそのサーバーを返す
	if server, exists := lb.sessions[sessionID]; exists {
		fmt.Printf("  [Cookie] SessionID=%s -> Server=%s (sticky)\n", sessionID, server.Name)
		return server
	}

	// 新規セッション: ラウンドロビンで選択
	server := lb.servers[len(lb.sessions)%len(lb.servers)]
	lb.sessions[sessionID] = server
	server.AddLoad()
	fmt.Printf("  [Cookie] SessionID=%s -> Server=%s (new)\n", sessionID, server.Name)

	return server
}

// GetSessionCount はセッション数を取得
func (lb *CookieBasedStickySession) GetSessionCount() int {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	return len(lb.sessions)
}

// ============================================================
// 2. IP Hash-Based Sticky Session
// ============================================================

// IPHashStickySession はIPハッシュベースのSticky Sessionを実装
type IPHashStickySession struct {
	servers []*Server
	mu      sync.RWMutex
}

// NewIPHashStickySession は新しいIP Hash-Based Sticky Sessionを作成
func NewIPHashStickySession(servers []*Server) *IPHashStickySession {
	return &IPHashStickySession{
		servers: servers,
	}
}

// GetServer はクライアントIPに基づいてサーバーを取得
func (lb *IPHashStickySession) GetServer(clientIP string) *Server {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	// IPアドレスをハッシュ化してサーバーを決定
	hash := lb.hashIP(clientIP)
	serverIndex := hash % uint32(len(lb.servers))
	server := lb.servers[serverIndex]

	fmt.Printf("  [IP Hash] ClientIP=%s -> Server=%s (hash=%d)\n", clientIP, server.Name, serverIndex)

	return server
}

// hashIP はIPアドレスをハッシュ化
func (lb *IPHashStickySession) hashIP(ip string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(ip))
	return h.Sum32()
}

// ============================================================
// 3. Consistent Hashing Sticky Session
// ============================================================

// ConsistentHashRing はコンシステントハッシュリングを実装
type ConsistentHashRing struct {
	ring         map[uint32]*Server // ハッシュ値 -> サーバー
	sortedKeys   []uint32           // ソート済みハッシュ値
	virtualNodes int                // 仮想ノード数
	mu           sync.RWMutex
}

// NewConsistentHashRing は新しいコンシステントハッシュリングを作成
func NewConsistentHashRing(servers []*Server, virtualNodes int) *ConsistentHashRing {
	ring := &ConsistentHashRing{
		ring:         make(map[uint32]*Server),
		virtualNodes: virtualNodes,
	}

	// 各サーバーに対して仮想ノードを作成
	for _, server := range servers {
		for i := 0; i < virtualNodes; i++ {
			hash := ring.hash(fmt.Sprintf("%s:%d", server.ID, i))
			ring.ring[hash] = server
			ring.sortedKeys = append(ring.sortedKeys, hash)
		}
	}

	// ハッシュ値をソート
	ring.sortKeys()

	return ring
}

// GetServer はキーに基づいてサーバーを取得
func (ring *ConsistentHashRing) GetServer(key string) *Server {
	ring.mu.RLock()
	defer ring.mu.RUnlock()

	if len(ring.ring) == 0 {
		return nil
	}

	hash := ring.hash(key)

	// バイナリサーチで最も近いノードを探す
	idx := ring.search(hash)
	server := ring.ring[ring.sortedKeys[idx]]

	fmt.Printf("  [Consistent Hash] Key=%s -> Server=%s (hash=%d)\n", key, server.Name, hash)

	return server
}

// hash はキーをハッシュ化
func (ring *ConsistentHashRing) hash(key string) uint32 {
	h := md5.New()
	h.Write([]byte(key))
	hashBytes := h.Sum(nil)
	return uint32(hashBytes[0])<<24 | uint32(hashBytes[1])<<16 | uint32(hashBytes[2])<<8 | uint32(hashBytes[3])
}

// search はハッシュ値に最も近いノードのインデックスを検索
func (ring *ConsistentHashRing) search(hash uint32) int {
	// バイナリサーチ
	left, right := 0, len(ring.sortedKeys)-1

	for left <= right {
		mid := (left + right) / 2
		if ring.sortedKeys[mid] == hash {
			return mid
		}
		if ring.sortedKeys[mid] < hash {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// 見つからない場合は最初のノード
	if left >= len(ring.sortedKeys) {
		return 0
	}
	return left
}

// sortKeys はハッシュ値をソート（バブルソート）
func (ring *ConsistentHashRing) sortKeys() {
	n := len(ring.sortedKeys)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if ring.sortedKeys[j] > ring.sortedKeys[j+1] {
				ring.sortedKeys[j], ring.sortedKeys[j+1] = ring.sortedKeys[j+1], ring.sortedKeys[j]
			}
		}
	}
}

// ============================================================
// 4. Session Replication (Shared Session)
// ============================================================

// SessionStore はセッションストアを表す
type SessionStore struct {
	sessions map[string]map[string]interface{} // sessionID -> data
	mu       sync.RWMutex
}

// NewSessionStore は新しいセッションストアを作成
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]map[string]interface{}),
	}
}

// Set はセッションデータを保存
func (store *SessionStore) Set(sessionID, key string, value interface{}) {
	store.mu.Lock()
	defer store.mu.Unlock()

	if store.sessions[sessionID] == nil {
		store.sessions[sessionID] = make(map[string]interface{})
	}
	store.sessions[sessionID][key] = value
}

// Get はセッションデータを取得
func (store *SessionStore) Get(sessionID, key string) (interface{}, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if session, exists := store.sessions[sessionID]; exists {
		if value, exists := session[key]; exists {
			return value, true
		}
	}
	return nil, false
}

// SessionReplicationLoadBalancer はセッションレプリケーション型のロードバランサー
type SessionReplicationLoadBalancer struct {
	servers      []*Server
	sessionStore *SessionStore // 共有セッションストア
	mu           sync.RWMutex
}

// NewSessionReplicationLoadBalancer は新しいセッションレプリケーション型LBを作成
func NewSessionReplicationLoadBalancer(servers []*Server) *SessionReplicationLoadBalancer {
	return &SessionReplicationLoadBalancer{
		servers:      servers,
		sessionStore: NewSessionStore(),
	}
}

// GetServer はラウンドロビンでサーバーを取得（セッションは共有）
func (lb *SessionReplicationLoadBalancer) GetServer(requestCount int) *Server {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	server := lb.servers[requestCount%len(lb.servers)]
	fmt.Printf("  [Session Replication] Request=%d -> Server=%s (session shared)\n", requestCount, server.Name)

	return server
}

// SetSessionData はセッションデータを保存
func (lb *SessionReplicationLoadBalancer) SetSessionData(sessionID, key string, value interface{}) {
	lb.sessionStore.Set(sessionID, key, value)
	fmt.Printf("  [Session Store] Set: SessionID=%s, Key=%s, Value=%v\n", sessionID, key, value)
}

// GetSessionData はセッションデータを取得
func (lb *SessionReplicationLoadBalancer) GetSessionData(sessionID, key string) (interface{}, bool) {
	value, exists := lb.sessionStore.Get(sessionID, key)
	if exists {
		fmt.Printf("  [Session Store] Get: SessionID=%s, Key=%s, Value=%v\n", sessionID, key, value)
	}
	return value, exists
}

// ============================================================
// デモ関数
// ============================================================

// DemoStickySession はSticky Sessionのデモを実行
func DemoStickySession() {
	fmt.Println("==============================================")
	fmt.Println("Sticky Sessionの実装と比較")
	fmt.Println("==============================================\n")

	// サーバーの準備
	servers := []*Server{
		NewServer("s1", "Server-1"),
		NewServer("s2", "Server-2"),
		NewServer("s3", "Server-3"),
	}

	demoCookieBased(servers)
	fmt.Println()

	demoIPHashBased(servers)
	fmt.Println()

	demoConsistentHashing(servers)
	fmt.Println()

	demoSessionReplication(servers)
}

// demoCookieBased はCookie-Based Sticky Sessionのデモ
func demoCookieBased(servers []*Server) {
	fmt.Println("1. Cookie-Based Sticky Session")
	fmt.Println("----------------------------------------------")
	fmt.Println("セッションIDをCookieに保存し、同じサーバーに振り分ける")
	fmt.Println()

	lb := NewCookieBasedStickySession(servers)

	// シミュレーション
	fmt.Println("【シナリオ】3つのセッションIDから5回のリクエスト:")
	sessionIDs := []string{"session-alice", "session-bob", "session-alice", "session-charlie", "session-alice"}

	for i, sessionID := range sessionIDs {
		fmt.Printf("リクエスト %d:\n", i+1)
		lb.GetServer(sessionID)
	}

	fmt.Println()
	fmt.Printf("総セッション数: %d\n", lb.GetSessionCount())
	fmt.Println()
	printServerLoad(servers)

	fmt.Println()
	fmt.Println("特徴:")
	fmt.Println("  ✓ セッションごとに同じサーバーにルーティング")
	fmt.Println("  ✓ セッション情報をCookieで管理")
	fmt.Println("  ✗ Cookieが削除されると新しいサーバーに振り分けられる")
	fmt.Println("  ✗ ロードバランサーがセッションテーブルを保持")
}

// demoIPHashBased はIP Hash-Based Sticky Sessionのデモ
func demoIPHashBased(servers []*Server) {
	fmt.Println("2. IP Hash-Based Sticky Session")
	fmt.Println("----------------------------------------------")
	fmt.Println("クライアントIPをハッシュ化して同じサーバーに振り分ける")
	fmt.Println()

	lb := NewIPHashStickySession(servers)

	// サーバーの負荷をリセット
	resetServerLoad(servers)

	// シミュレーション
	fmt.Println("【シナリオ】3つのクライアントIPから5回のリクエスト:")
	clientIPs := []string{"192.168.1.10", "192.168.1.20", "192.168.1.10", "192.168.1.30", "192.168.1.10"}

	for i, ip := range clientIPs {
		fmt.Printf("リクエスト %d:\n", i+1)
		server := lb.GetServer(ip)
		server.AddLoad()
	}

	fmt.Println()
	printServerLoad(servers)

	fmt.Println()
	fmt.Println("特徴:")
	fmt.Println("  ✓ 同じIPからのリクエストは同じサーバーに")
	fmt.Println("  ✓ セッション管理が不要")
	fmt.Println("  ✗ NAT環境では複数ユーザーが同じサーバーに集中")
	fmt.Println("  ✗ サーバー追加/削除時に振り分けが変わる")
}

// demoConsistentHashing はConsistent Hashingのデモ
func demoConsistentHashing(servers []*Server) {
	fmt.Println("3. Consistent Hashing")
	fmt.Println("----------------------------------------------")
	fmt.Println("仮想ノードを使ったハッシュリングで分散")
	fmt.Println()

	virtualNodes := 3
	ring := NewConsistentHashRing(servers, virtualNodes)

	// サーバーの負荷をリセット
	resetServerLoad(servers)

	// シミュレーション
	fmt.Println("【シナリオ】5つのセッションIDをコンシステントハッシュで振り分け:")
	sessionIDs := []string{"user-1", "user-2", "user-3", "user-4", "user-5"}

	for i, sessionID := range sessionIDs {
		fmt.Printf("リクエスト %d:\n", i+1)
		server := ring.GetServer(sessionID)
		server.AddLoad()
	}

	fmt.Println()
	printServerLoad(servers)

	fmt.Println()
	fmt.Printf("仮想ノード数: %d (各サーバー)\n", virtualNodes)
	fmt.Printf("ハッシュリングの総ノード数: %d\n", len(servers)*virtualNodes)

	fmt.Println()
	fmt.Println("特徴:")
	fmt.Println("  ✓ サーバー追加/削除時の影響が最小限")
	fmt.Println("  ✓ 仮想ノードで負荷を均等分散")
	fmt.Println("  ✓ スケーラビリティが高い")
	fmt.Println("  ✗ 実装が複雑")
}

// demoSessionReplication はSession Replicationのデモ
func demoSessionReplication(servers []*Server) {
	fmt.Println("4. Session Replication (Shared Session)")
	fmt.Println("----------------------------------------------")
	fmt.Println("セッションを外部ストア（Redis等）で共有")
	fmt.Println()

	lb := NewSessionReplicationLoadBalancer(servers)

	// サーバーの負荷をリセット
	resetServerLoad(servers)

	// シミュレーション
	fmt.Println("【シナリオ】セッションを共有してラウンドロビン:")

	sessionID := "session-alice"

	fmt.Println("\n1. 最初のリクエスト（セッション作成）:")
	server1 := lb.GetServer(0)
	server1.AddLoad()
	lb.SetSessionData(sessionID, "username", "Alice")
	lb.SetSessionData(sessionID, "cart_items", 3)

	fmt.Println("\n2. 2回目のリクエスト（別のサーバーでもセッションアクセス可能）:")
	server2 := lb.GetServer(1)
	server2.AddLoad()
	username, _ := lb.GetSessionData(sessionID, "username")
	cartItems, _ := lb.GetSessionData(sessionID, "cart_items")
	fmt.Printf("  Retrieved: username=%v, cart_items=%v\n", username, cartItems)

	fmt.Println("\n3. 3回目のリクエスト（さらに別のサーバー）:")
	server3 := lb.GetServer(2)
	server3.AddLoad()
	username, _ = lb.GetSessionData(sessionID, "username")
	fmt.Printf("  Retrieved: username=%v\n", username)

	fmt.Println()
	printServerLoad(servers)

	fmt.Println()
	fmt.Println("特徴:")
	fmt.Println("  ✓ どのサーバーでもセッションにアクセス可能")
	fmt.Println("  ✓ サーバー障害時もセッションが失われない")
	fmt.Println("  ✓ 柔軟なロードバランシング（ラウンドロビン等）")
	fmt.Println("  ✗ 外部ストア（Redis等）への依存")
	fmt.Println("  ✗ ネットワークオーバーヘッド")
}

// printServerLoad はサーバーの負荷を表示
func printServerLoad(servers []*Server) {
	fmt.Println("サーバー負荷:")
	for _, server := range servers {
		load := server.GetLoad()
		fmt.Printf("  %s: %d リクエスト\n", server.Name, load)
	}
}

// resetServerLoad はサーバーの負荷をリセット
func resetServerLoad(servers []*Server) {
	for _, server := range servers {
		server.mu.Lock()
		server.Load = 0
		server.mu.Unlock()
	}
}

// ============================================================
// パフォーマンス比較デモ
// ============================================================

// DemoStickySessionPerformance はパフォーマンス比較デモ
func DemoStickySessionPerformance() {
	fmt.Println("\n==============================================")
	fmt.Println("パフォーマンス比較")
	fmt.Println("==============================================\n")

	servers := []*Server{
		NewServer("s1", "Server-1"),
		NewServer("s2", "Server-2"),
		NewServer("s3", "Server-3"),
	}

	requestCount := 10000

	fmt.Printf("各方式で %d リクエストを処理\n\n", requestCount)

	// Cookie-Based
	fmt.Println("1. Cookie-Based:")
	start := time.Now()
	cookieLB := NewCookieBasedStickySession(servers)
	for i := 0; i < requestCount; i++ {
		sessionID := generateSessionID(i)
		cookieLB.GetServer(sessionID)
	}
	fmt.Printf("  所要時間: %v\n", time.Since(start))
	fmt.Printf("  セッション数: %d\n", cookieLB.GetSessionCount())

	// IP Hash
	fmt.Println("\n2. IP Hash:")
	start = time.Now()
	ipHashLB := NewIPHashStickySession(servers)
	for i := 0; i < requestCount; i++ {
		clientIP := generateClientIP(i)
		ipHashLB.GetServer(clientIP)
	}
	fmt.Printf("  所要時間: %v\n", time.Since(start))

	// Consistent Hashing
	fmt.Println("\n3. Consistent Hashing:")
	start = time.Now()
	ring := NewConsistentHashRing(servers, 100)
	for i := 0; i < requestCount; i++ {
		key := generateSessionID(i)
		ring.GetServer(key)
	}
	fmt.Printf("  所要時間: %v\n", time.Since(start))

	// Session Replication
	fmt.Println("\n4. Session Replication:")
	start = time.Now()
	sessionLB := NewSessionReplicationLoadBalancer(servers)
	for i := 0; i < requestCount; i++ {
		sessionLB.GetServer(i)
	}
	fmt.Printf("  所要時間: %v\n", time.Since(start))

	fmt.Println("\n結論:")
	fmt.Println("  - IP Hash: 最速（ハッシュ計算のみ）")
	fmt.Println("  - Consistent Hash: やや遅い（リング検索）")
	fmt.Println("  - Cookie-Based: セッション管理のオーバーヘッド")
	fmt.Println("  - Session Replication: 外部ストアアクセスで最も遅い")
}

// generateSessionID はランダムなセッションIDを生成
func generateSessionID(seed int) string {
	rand.Seed(int64(seed))
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateClientIP はランダムなクライアントIPを生成
func generateClientIP(seed int) string {
	rand.Seed(int64(seed))
	return fmt.Sprintf("192.168.%d.%d", rand.Intn(256), rand.Intn(256))
}
