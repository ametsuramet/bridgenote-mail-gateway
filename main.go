package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	sqsApi "mail_gateway/sqs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

// MailRequest represents the structure of the JSON request
type MailRequest struct {
	ProcessType string `json:"process_type"`
	MessageType string `json:"message_type"`
	Mail        struct {
		To      string `json:"to"`
		CC      string `json:"cc"`
		BCC     string `json:"bcc"`
		ReplyTo string `json:"reply-to"`
		Title   string `json:"title"`
		Body    string `json:"body"`
	} `json:"mail"`
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create an Echo instance
	e := echo.New()

	// Define a route to handle the POST request
	e.POST("/send-mail", handlePostRequest)

	// Start the server
	e.Start(":2340")
}

func handlePostRequest(c echo.Context) error {
	// Bind JSON request to struct
	var request MailRequest
	if err := c.Bind(&request); err != nil {
		fmt.Println("ERROR", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to parse JSON request"})
	}

	// Print the received struct for demonstration
	fmt.Printf("Received request: %+v\n", request)

	// Forward the request to AWS SQS
	if err := sendToSQS(request, os.Getenv("AWS_REGION")); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send message to SQS"})
	}

	// Example response
	response := map[string]string{"status": "success"}

	// Return the response
	return c.JSON(http.StatusOK, response)
}

func sendToSQS(request MailRequest, region string) error {
	queueName := os.Getenv("AWS_SQS_QUEUE")
	// Create an SQS client
	client := Client(context.TODO(), os.Getenv("AWS_SQS_URL"), region)

	// Generate SQS queue URL based on some criteria (you can customize this function)
	// queueURL := getSQSQueueURL("your-queue-name", region)
	input := &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	}
	resultGet, err := sqsApi.GetQueueURL(context.TODO(), client, input)
	if err != nil {
		log.Printf("error getting the queue URL: %v", err)
	}
	queueURL := resultGet.QueueUrl
	fmt.Println("queueURL", *queueURL)

	// Convert the request to JSON
	jsonData, err := json.Marshal(SqsPayload{
		Type: "TYPE_SEND_MAIL_WITH_TEMPLATE",
		Data: MailData{
			Subject: request.Mail.Title,
			Body:    request.Mail.Body,
			Email:   request.Mail.To,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	stringData := string(jsonData)
	fmt.Printf("Send request: %+v\n", stringData)

	// Send the message to SQS
	_, err = client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody:  &stringData,
		QueueUrl:     queueURL,
		DelaySeconds: 0,
	})
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %v", err)
	}

	return nil
}

func Client(ctx context.Context, awsURL, region string) *sqs.Client {
	// customResolver is required here since we use localstack and need to point the aws url to localhost.
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           awsURL,
			SigningRegion: region,
		}, nil

	})

	// load the default aws config along with custom resolver.
	cfg, err := config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	return sqs.NewFromConfig(cfg)
}

type SqsPayload struct {
	Type  string   `json:"Type"`
	RefID string   `json:"RefID"`
	Data  MailData `json:"Data"`
}

type MailData struct {
	Subject string `json:"Subject"`
	Body    string `json:"Body"`
	Email   string `json:"Email"`
}
