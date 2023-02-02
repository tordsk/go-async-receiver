package main

import (
	"encoding/json"
	"fmt"
	"github.com/appleboy/go-fcm"
	pr "github.com/crow-misia/go-push-receiver"
	"os"
)

func main() {
	creds, _ := loadCreds()
	sender, _ := fcm.NewClient("") // api key
	send, _ := sender.Send(&fcm.Message{
		To: "",
		Notification: &fcm.Notification{
			Title: "Test Notification",
			Body:  "With Body",
		},
		Data: map[string]interface{}{
			"Test": "Hello World",
		},
	})
	fmt.Println(creds.Token)
	fmt.Println(send.Results[0])
	/*
	 {"multicast_id":4827136729016708860,"success":0,"failure":1,"canonical_ids":0,"results":[{"error":"InvalidParameters: unsupported token type"}]}
	*/
}

func loadCreds() (*pr.FCMCredentials, error) {
	content, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}
	creds := pr.FCMCredentials{}
	return &creds, json.Unmarshal(content, &creds)
}
