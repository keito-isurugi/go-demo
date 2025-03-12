package main

// Notifier は通知を送るためのインターフェース
type Notifier interface {
	// SendNotification は通知メッセージを送信する
	SendNotification(message string) string
}

// EmailNotifier は Email 通知の実装
type EmailNotifier struct{}
func (n *EmailNotifier) SendNotification(message string) string {
	return "Email: " + message
}

// SMSNotifier は SMS 通知の実装
type SMSNotifier struct{}
func (n *SMSNotifier) SendNotification(message string) string {
	return "SMS: " + message
}

// NotificationFactory は Notifier を作成するファクトリインターフェース
type NotificationFactory interface {
	// CreateNotifier は Notifier を作成する
	CreateNotifier() Notifier
}

// EmailNotificationFactory は Email 通知のファクトリ
type EmailNotificationFactory struct{}
func (f *EmailNotificationFactory) CreateNotifier() Notifier {
	return &EmailNotifier{}
}

// SMSNotificationFactory は SMS 通知のファクトリ
type SMSNotificationFactory struct{}
func (f *SMSNotificationFactory) CreateNotifier() Notifier {
	return &SMSNotifier{}
}

// SendNotification はファクトリを使って通知を送信
func SendNotification(factory NotificationFactory, message string) string {
	notifier := factory.CreateNotifier()
	return notifier.SendNotification(message)
}

// abstractFactoryExec は Abstract Factory パターンの使用例
func abstractFactoryExec() {
	// Email 通知の送信
	emailFactory := &EmailNotificationFactory{}
	emailMessage := SendNotification(emailFactory, "Hello, Email!")
	println(emailMessage)

	// SMS 通知の送信
	smsFactory := &SMSNotificationFactory{}
	smsMessage := SendNotification(smsFactory, "Hello, SMS!")
	println(smsMessage)
}


// 出力
// Email: Hello, Email!
// SMS: Hello, SMS!