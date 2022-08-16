package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/aws/aws-sdk-go/service/sesv2/sesv2iface"
	"github.com/google/uuid"
)

type CreateEmailRequestTests struct {
	name             string
	input            MessageBody
	expectedResponse sesv2.SendEmailInput
}

type HandlerTests struct {
	name             string
	input            events.SQSEvent
	expectedResponse error
}

type mockedSendEmail struct {
	sesv2iface.SESV2API
	Response sesv2.SendEmailOutput
}

func TestCreateEmailRequest(t *testing.T) {
	//Arrange
	fromEmail := "test@gmail.com"
	body := MessageBody{
		Message: "Some Text",
		To:      []string{"test@gmail.com"},
	}
	var toAddresses []*string

	for _, address := range body.To {
		toAddresses = append(toAddresses, &address)
	}

	tests := []CreateEmailRequestTests{
		{
			name:  "Given Message Should Return SendEmailInput",
			input: body,
			expectedResponse: sesv2.SendEmailInput{
				Content: &sesv2.EmailContent{
					Raw: &sesv2.RawMessage{
						Data: []byte("Some Text"),
					},
				},
				FromEmailAddress: &fromEmail,
				Destination: &sesv2.Destination{
					ToAddresses: toAddresses,
				},
			},
		},
	}

	//Act
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := createEmailRequest(test.input)

			if !reflect.DeepEqual(resp, test.expectedResponse) {
				t.Errorf("FAILED: expected %v, got %v\n", test.expectedResponse, resp)
			}
		})
	}
}

func (d mockedSendEmail) SendEmail(input *sesv2.SendEmailInput) (*sesv2.SendEmailOutput, error) {
	return &d.Response, nil
}

func TestHandler(t *testing.T) {
	//Arrange
	m := mockedSendEmail{
		Response: sesv2.SendEmailOutput{},
	}

	d := deps{
		ses: m.SESV2API,
	}

	tests := []HandlerTests{
		{
			name: "Given SQSEvent with Message Should Return No Error",
			input: events.SQSEvent{
				Records: []events.SQSMessage{
					{
						MessageId: uuid.New().String(),
						Body:      `{ "message": "some text", "emailTo": ["test@gmail.com"] }`,
					},
				},
			},
			expectedResponse: nil,
		},
	}

	//Act
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := d.handler(context.TODO(), test.input)

			if resp != test.expectedResponse {
				t.Errorf("FAILED: expected %v, got %v\n", test.expectedResponse, resp)
			}
		})
	}
}
