package main

import (
    "fmt"
)

// ロータをシミュレートする
var rotor = []string{
    "EKMFLGDQVZNTOWYHXUSPAIBRCJ", // ロータ1
    "AJDKSIRUXBLHWTMCQGZNPYFVOE", // ロータ2
    "BDFHJLCPRTXVZNYEIWGAKMUSQO", // ロータ3
}

// リフレクタをシミュレートする
var reflector = "YRUHQSLDPXNGOKMIEBFZCWVJAT"

// プラグボード設定
var plugboard = map[rune]rune{
    'A': 'B', 'B': 'A',
    'C': 'D', 'D': 'C',
    'E': 'F', 'F': 'E',
    // 他の設定を追加可能
}

func rotate(rotor string, steps int) string {
    return rotor[steps:] + rotor[:steps]
}

func encryptChar(c rune, rotors []string, reflector string) rune {
    index := int(c - 'A')

    // 前進処理
    for i := 0; i < len(rotors); i++ {
       index = int(rotors[i][index] - 'A')
    }

    // 反射
    index = int(reflector[index] - 'A')

    // 後退処理
    for i := len(rotors) - 1; i >= 0; i-- {
       for j := 0; j < 26; j++ {
          if int(rotors[i][j]-'A') == index {
             index = j
             break
          }
       }
    }

    return rune(index + 'A')
}

func plugboardSwap(c rune) rune {
    if swapped, exists := plugboard[c]; exists {
       return swapped
    }
    return c
}

func processMessage(message string) string {
    processed := make([]rune, len(message))
    rotorPositions := []int{0, 0, 0}
    for i, char := range message {
       pluggedChar := plugboardSwap(char)

       // 現在の rotor 位置で文字をエンコード
       rotatedRotors := make([]string, len(rotor))
       for j := range rotor {
          rotatedRotors[j] = rotate(rotor[j], rotorPositions[j])
       }

       // 暗号化処理
       encryptedChar := encryptChar(pluggedChar, rotatedRotors, reflector)

       // プラグボードの確認
       processed[i] = plugboardSwap(encryptedChar)

       // ロータを回転
       rotorPositions[0] = (rotorPositions[0] + 1) % 26
       if rotorPositions[0] == 0 {
          rotorPositions[1] = (rotorPositions[1] + 1) % 26
          if rotorPositions[1] == 0 {
             rotorPositions[2] = (rotorPositions[2] + 1) % 26
          }
       }
    }
    return string(processed)
}

func enigmaCipherExec() {
    message := "HELLOENIGMA"
    encryptedMessage := processMessage(message)

    // ロータの位置をリセットして復号
    decryptedMessage := processMessage(encryptedMessage)

    fmt.Println("元のメッセージ:", message)
    fmt.Println("暗号化メッセージ:", encryptedMessage)
    fmt.Println("復号化メッセージ:", decryptedMessage)
}