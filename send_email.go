package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sesv2"

	"fmt"
)

var sesSvc *sesv2.SESV2

const (
	REGION = "us-west-1"
)

func init() {
	// Initialize a new session that the SDK uses to load
	newSession, _ := session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	})

	// Create an Amazon SES Client
	sesSvc = sesv2.New(newSession)
}

func main() {
	// Tell the lambda what to run
	lambda.Start(handler)
}

func handler() {
	emailRequest := createEmailRequest()

	result, err := sesSvc.SendEmail(&emailRequest)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print("Sent email with MessageId: ")
	fmt.Println(result.MessageId)

	fmt.Println("Finished sending emails.")
}

func createEmailRequest() sesv2.SendEmailInput {
	rawMessage := sesv2.RawMessage{
		Data: []byte("some text"),
	}

	emailContent := sesv2.EmailContent{
		Raw: &rawMessage,
	}

	emailRequest := sesv2.SendEmailInput{
		Content: &emailContent,
	}

	return emailRequest
}
