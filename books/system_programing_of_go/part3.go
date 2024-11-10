package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

func Part3() {
	for {
		buffer := make([]byte, 5)
		size, err := os.Stdin.Read(buffer)
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		fmt.Printf("size=%d input='%s'\n", size, string(buffer))
	}

	file, err := os.Open("part1.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.Copy(os.Stdout, file)

	conn, _ := net.Dial("tcp", "example.com:80")
	conn.Write([]byte("GET / HTTP/1.0\r\nHost: example.com\r\n\r\n"))
	res, _ := http.ReadResponse(bufio.NewReader(conn), nil)
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)

	reader := strings.NewReader("Exmaple of io.sectionReader\n")
	sectionReader := io.NewSectionReader(reader, 1, 7)
	io.Copy(os.Stdout, sectionReader)

	data := []byte{0x0, 0x0, 0x27, 0x10}
	var i int32
	binary.Read(bytes.NewReader(data), binary.BigEndian, &i)
	fmt.Printf("data: %d\n", i)

	file2, _ := os.Open("sample_png_1.png")
	defer file2.Close()
	newFile, _ := os.Create("sample_png_1_create.png")
	defer newFile.Close()
	chunks := readChunks(file2)
	// シグニチャ書き込み
	io.WriteString(newFile, "\x89PNG\r\n\x1a\n")
	// 先頭に必要なIHDRチャンクを書き込み
	io.Copy(newFile, chunks[0])
	// テキストチャンクを追加
	io.Copy(newFile, textChunk("ASCII PROGRAMMING++"))
	// 残りのチャンクを追加
	for _, chunk := range chunks[1:] {
		io.Copy(newFile, chunk)
	}

	chunks = readChunks(newFile)
	for _, chunk := range chunks {
		dumpChunk(chunk)
	}

	var source = `1行目
	2行目
	3行目`

	reader2 := bufio.NewReader(strings.NewReader(source))
	for {
		line, err := reader2.ReadString('\n')
		fmt.Printf("%#v\n", line)
		if err == io.EOF {
			break
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(source))
	for scanner.Scan() {
		fmt.Printf("%#v\n", scanner.Text())
	}

	var source2 = "123 1.234 1.0e4 test"
	reader3 := strings.NewReader(source2)
	var i2 int
	var f, g float64
	var s string
	fmt.Fscan(reader3, &i2, &f, &g, &s)
	fmt.Printf("i=%#v f=%#v g=%#v s=%#v\n", i2, f, g, s)

	var csvSource = `13101,"100 ","1000003"," ﾄｳｷｮｳﾄ "," ﾁﾖﾀﾞｸ "," ﾋﾄﾂﾊﾞｼ (1 ﾁｮｳﾒ )"," 東京都 "," 千代田区 "," 一ツ橋（１丁目）",1,0,1,0,0,0
13101,"101 ","1010003"," ﾄｳｷｮｳﾄ "," ﾁﾖﾀﾞｸ "," ﾋﾄﾂﾊﾞｼ (2 ﾁｮｳﾒ )"," 東京都 "," 千代田区 "," 一ツ橋（２丁目）",1,0,1,0,0,0
13101,"100 ","1000012"," ﾄｳｷｮｳﾄ "," ﾁﾖﾀﾞｸ "," ﾋﾋﾞﾔｺｳｴﾝ "," 東京都 "," 千代田区 "," 日比谷公園 ",0,0,0,0,0,0
13101,"102 ","1020093"," ﾄｳｷｮｳﾄ "," ﾁﾖﾀﾞｸ "," ﾋﾗｶﾜﾁｮｳ "," 東京都 "," 千代田区 "," 平河町 ",0,0,1,0,0,0
13101,"102 ","1020071"," ﾄｳｷｮｳﾄ "," ﾁﾖﾀﾞｸ "," ﾌｼﾞﾐ "," 東京都 "," 千代田区 "," 富士見 ",0,0,1,0,0,0`

	reader4 := strings.NewReader(csvSource)
	csvReader := csv.NewReader(reader4)
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(line[2], line[6:9])
	}

	header := bytes.NewBufferString("----------- HEADER ----------\n")
	content := bytes.NewBufferString("Example of io.MultiReader\n")
	footer := bytes.NewBufferString("----------- FOOTER ----------\n")
	reader5 := io.MultiReader(header, content, footer)
	io.Copy(os.Stdout, reader5)

	var buffer bytes.Buffer
	reader6 := bytes.NewBufferString("Example of io.TeeReader\n")
	teeReader := io.TeeReader(reader6, &buffer)
	_, _ = io.ReadAll(teeReader)
	fmt.Println(buffer.String())
}

func dumpChunk(chunk io.Reader) {
	var length int32
	binary.Read(chunk, binary.BigEndian, &length)
	buffer := make([]byte, 4)
	chunk.Read(buffer)
	fmt.Printf("chunk '%v' (%d bytes)\n", string(buffer), length)
	if bytes.Equal(buffer, []byte("tEXt")) {
		rawText := make([]byte, length)
		chunk.Read(rawText)
		fmt.Println(string(rawText))
	}
}

func readChunks(file *os.File) []io.Reader {
	// チャンクを格納する配列
	var chunks []io.Reader

	// 最初の8バイトを飛ばす
	file.Seek(8, 0)
	var offset int64 = 8

	for {
		var lenght int32
		err := binary.Read(file, binary.BigEndian, &lenght)
		if err == io.EOF {
			break
		}
		chunks = append(chunks, io.NewSectionReader(file, offset, int64(lenght)+12))

		// 次のチャンクの先頭に移動
		// 現在位置は長さを読み終わった箇所なので
		// チャンク名(4バイト) + データ長 + CRC(4バイト)先に移動
		offset, _ = file.Seek(int64(lenght+8), 1)
	}
	return chunks
}

func textChunk(text string) io.Reader {
	byteData := []byte(text)
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, int32(len(byteData)))
	buffer.WriteString("tEXt")
	buffer.Write(byteData)
	// CRCを計算して追加
	crc := crc32.NewIEEE()
	io.WriteString(crc, "tEXt")
	crc.Write(byteData)
	binary.Write(&buffer, binary.BigEndian, crc.Sum32())
	return &buffer
}
