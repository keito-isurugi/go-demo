package main

type Notifier interface {
	SendNotification(message string) string
}

type EmailNotifier struct{}
func (n *EmailNotifier) SendNotification(message string) string {
	return "Email: " + message
}

type SMSNotifier struct{}
func (n *SMSNotifier) SendNotification(message string) string {
	return "SMS: " + message
}

type NotificationFactory interface {
	CreateNotifier() Notifier
}

type EmailNotificationFactory struct{}
func (f *EmailNotificationFactory) CreateNotifier() Notifier {
	return &EmailNotifier{}
}

type SMSNotificationFactory struct{}
func (f *SMSNotificationFactory) CreateNotifier() Notifier {
	return &SMSNotifier{}
}

func SendNotification(factory NotificationFactory, message string) string {
	notifier := factory.CreateNotifier()
	return notifier.SendNotification(message)
}

func abstractFactoryExec() {
	emailFactory := &EmailNotificationFactory{}
	emailMessage := SendNotification(emailFactory, "Hello, Email!")
	println(emailMessage)

	smsFactory := &SMSNotificationFactory{}
	smsMessage := SendNotification(smsFactory, "Hello, SMS!")
	println(smsMessage)
}