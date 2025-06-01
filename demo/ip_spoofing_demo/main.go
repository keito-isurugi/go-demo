package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

func main() {
	bootType := flag.String("boottype", "receive", "send or receive")
	flag.Parse()
	fmt.Println("BootType: " + *bootType)
	if *bootType == "receive" {
		// 受信側の処理
		fmt.Println("Receiving ICMP packets...")
		receive()
	}
	if *bootType == "send" {
		// 送信側の処理
		fmt.Println("Sending spoofed ICMP packet...")
		send()
	}
}

// ＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝送信側＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝
func send() {
	// 送信先アドレス
	dstIP := net.ParseIP("127.0.0.1") // ローカル
	// 偽装する送信元アドレス
	srcIP := net.ParseIP("192.0.2.123") // 任意の偽装アドレス

	// raw socket作成
	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("ListenPacket error: %v", err)
	}
	defer conn.Close()

	// ICMPエコーリクエストパケット作成
	icmp := []byte{
		8, 0, 0, 0, // Type, Code, Checksum(仮)
		0, 13, 0, 37, // ID, Sequence
	}
	// チェックサム計算
	sum := checkSum(icmp)
	icmp[2] = byte(sum >> 8)
	icmp[3] = byte(sum & 0xff)

	rawConn, err := ipv4.NewRawConn(conn)
	if err != nil {
		log.Fatalf("NewRawConn error: %v", err)
	}

	// IPヘッダ手動で作成
	header := &ipv4.Header{
		Version:  4,
		Len:      20,
		TotalLen: 20 + len(icmp),
		TTL:      64,
		Protocol: 1, // ICMP
		Src:      srcIP,
		Dst:      dstIP,
	}

	// 送信
	if err := rawConn.WriteTo(header, icmp, nil); err != nil {
		log.Fatalf("WriteTo error: %v", err)
	}
	log.Println("偽装パケット送信完了")
}

// チェックサム計算関数
func checkSum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for (sum >> 16) != 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return ^uint16(sum)
}

// ＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝受信側＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝
func receive() {
	// ICMPパケットを全て受信するためのraw socket作成
	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("ListenPacket error: %v", err)
	}
	defer conn.Close()

	buf := make([]byte, 1500)
	fmt.Println("ICMP受信待ち...")

	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Printf("ReadFrom error: %v", err)
			continue
		}

		// IPヘッダをパースして送信元IPを取得
		packet := buf[:n]
		// Raw ICMPなのでIPヘッダは取得できないが、ReadFromでaddrとしてIPが取れる
		ipAddr, ok := addr.(*net.IPAddr)
		if ok {
			fmt.Printf("受信: 送信元IP = %s, データ長 = %d\n", ipAddr.IP.String(), n)
			// ICMPデータ内容も表示（必要なら下記コメント解除）
			fmt.Printf("ICMPデータ: %x\n", packet)
		}
	}
}
