package main

import "github.com/appleboy/go-fcm"

func main() {
	sender, _ := fcm.NewClient("") // api key
	sender.Send(&fcm.Message{
		RegistrationIDs: []string{"cyeBz3DDOuE:APA91bE9HGc4PAX01K8ztaoSpaLvlzdIwVJSk7mExUolyf2joddxxO-kLvwaWnDmClDMAQMeIwwFUBs-M-6tH1EuE6QXjtoE7gwnMslLYfHkdPuBq6uyfROhnbejWuvz5JYEH9iT-utm"},
		Notification: &fcm.Notification{
			Title: "Test Notification",
			Body:  "With Body",
		},
		Data: map[string]interface{}{
			"Test": "Hello World",
		},
	})
}
