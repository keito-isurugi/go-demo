package main

import (
    "crypto/rand"
    "fmt"
    "strings"
    "unicode"
)

func vigenereExec() {
    textToEncrypt := "Hello World"
    cipher := NewVigenereCipher(textToEncrypt)

    encrypted := cipher.EncryptText(textToEncrypt)
    decrypted := cipher.DecryptText(encrypted)

    fmt.Printf("Original Text: %s\n", textToEncrypt)
    fmt.Printf("Random Key: %s\n", cipher.Key)
    fmt.Printf("Encrypted: %s\n", encrypted)
    fmt.Printf("Decrypted: %s\n", decrypted)
}

type VigenereCipher struct {
    Key string
}

func NewVigenereCipher(text string) VigenereCipher {
    return VigenereCipher{Key: generateRandomAlphaKey(len(text))}
}

func generateRandomAlphaKey(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    key := make([]byte, length)
    _, err := rand.Read(key)
    if err != nil {
       panic(err)
    }
    for i := 0; i < length; i++ {
       key[i] = charset[key[i]%byte(len(charset))]
    }
    return string(key)
}

func (vc VigenereCipher) EncryptText(text string) string {
    return vc.transformText(text, vc.encryptShift)
}

func (vc VigenereCipher) DecryptText(text string) string {
    return vc.transformText(text, vc.decryptShift)
}

func (vc VigenereCipher) transformText(text string, shiftFunc func(rune, rune, rune) rune) string {
    var result strings.Builder
    keyIndex := 0
    keyLength := len(vc.Key)

    for _, char := range text {
       if unicode.IsLetter(char) { // Unicode上で文字と分類される場合のみを暗号化の対象とする
          base := rune('A')
          if unicode.IsLower(char) { // 暗号化対象の文字が小文字の場合
             base = rune('a')
          }

          // 暗号化対象の文字に対するkeyを取得
          keyChar := rune(vc.Key[keyIndex])
          if unicode.IsLower(char) { // 暗号化対象の文字が小文字の場合、keyを小文字にする
             keyChar = unicode.ToLower(keyChar)
          } else { // 暗号化対象の文字が大文字の場合、keyを大文字にする
             keyChar = unicode.ToUpper(keyChar)
          }

          // 暗号化(復号化)
          newChar := shiftFunc(char, keyChar, base)
          result.WriteRune(newChar)

          keyIndex = (keyIndex + 1) % keyLength
       } else { // 文字以外は暗号化しない
          result.WriteRune(char)
       }
    }
    return result.String()
}

func (vc VigenereCipher) encryptShift(char, keyChar, base rune) rune {
    // ((暗号化対象の文字のオフセット) + (keyの文字のオフセット) + 負の値防止) % 0~25の範囲に収める + 実際の文字コードに戻す
    return ((char-base)+(keyChar-base)+26)%26 + base
}

func (vc VigenereCipher) decryptShift(char, keyChar, base rune) rune {
    return ((char-base)-(keyChar-base)+26)%26 + base
}