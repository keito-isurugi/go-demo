package bridge

import "fmt"

// Device インターフェース（家電の抽象化）
type Device interface {
	TurnOn()
	TurnOff()
	SetVolume(volume int)
}

// TV 構造体（具体的な家電: テレビ）
type TV struct {
	volume int
}

func (t *TV) TurnOn() {
	fmt.Println("テレビをつけました")
}

func (t *TV) TurnOff() {
	fmt.Println("テレビを消しました")
}

func (t *TV) SetVolume(volume int) {
	t.volume = volume
	fmt.Printf("テレビの音量を %d に設定しました\n", volume)
}

// Radio 構造体（具体的な家電: ラジオ）
type Radio struct {
	volume int
}

func (r *Radio) TurnOn() {
	fmt.Println("ラジオをつけました")
}

func (r *Radio) TurnOff() {
	fmt.Println("ラジオを消しました")
}

func (r *Radio) SetVolume(volume int) {
	r.volume = volume
	fmt.Printf("ラジオの音量を %d に設定しました\n", volume)
}

// RemoteControl 抽象構造体（リモコンの抽象化）
type RemoteControl struct {
	device Device
}

func (r *RemoteControl) TurnOn() {
	r.device.TurnOn()
}

func (r *RemoteControl) TurnOff() {
	r.device.TurnOff()
}

func (r *RemoteControl) VolumeUp() {
	r.device.SetVolume(10) // 仮に10をセット
}

// AdvancedRemote 構造体（拡張リモコン）
type AdvancedRemote struct {
	RemoteControl
}

func (a *AdvancedRemote) Mute() {
	a.device.SetVolume(0)
	fmt.Println("ミュートしました")
}

func BridgeExec() {
	tv := &TV{}
	radio := &Radio{}

	basicRemote := RemoteControl{device: tv}
	advancedRemote := AdvancedRemote{RemoteControl{device: radio}}

	basicRemote.TurnOn()
	basicRemote.VolumeUp()
	basicRemote.TurnOff()

	advancedRemote.TurnOn()
	advancedRemote.Mute()
	advancedRemote.TurnOff()
}

// 出力
// テレビをつけました
// テレビの音量を 10 に設定しました
// テレビを消しました
// ラジオをつけました
// ラジオの音量を 0 に設定しました
// ミュートしました
// ラジオを消しました
