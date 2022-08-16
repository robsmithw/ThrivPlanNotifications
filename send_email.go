package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/aws/aws-sdk-go/service/sesv2/sesv2iface"

	"fmt"
)

var sesSvc *sesv2.SESV2

const (
	REGION = "us-west-1"
)

type deps struct {
	ses sesv2iface.SESV2API
}

type MessageBody struct {
	Message string   `json:"message"`
	To      []string `json:"emailTo"`
}

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
	d := deps{
		ses: sesSvc,
	}

	lambda.Start(d.handler)
}

func (d *deps) handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {

		var body MessageBody

		if err := json.Unmarshal([]byte(message.Body), &body); err != nil {
			fmt.Println(err.Error())
			return err
		}

		emailRequest := createEmailRequest(body)

		result, err := sesSvc.SendEmail(&emailRequest)

		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		fmt.Print("Sent email with MessageId: ")
		fmt.Println(result.MessageId)
	}

	fmt.Println("Finished sending emails.")

	return nil
}

func createEmailRequest(body MessageBody) sesv2.SendEmailInput {
	fromEmail := "test@gmail.com"

	var toAddresses []*string

	for _, address := range body.To {
		toAddresses = append(toAddresses, &address)
	}

	rawMessage := sesv2.RawMessage{
		Data: []byte(body.Message),
	}

	emailContent := sesv2.EmailContent{
		Raw: &rawMessage,
	}

	emailRequest := sesv2.SendEmailInput{
		Content:          &emailContent,
		FromEmailAddress: &fromEmail,
		Destination: &sesv2.Destination{
			ToAddresses: toAddresses,
		},
	}

	return emailRequest
}
