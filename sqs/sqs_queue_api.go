package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// QueueAPI Interface SQS queue API
type QueueAPI interface {
	CreateQueue(ctx context.Context,
		params *sqs.CreateQueueInput,
		optFns ...func(*sqs.Options)) (*sqs.CreateQueueOutput, error)
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)
	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	DeleteQueue(ctx context.Context,
		params *sqs.DeleteQueueInput,
		optFns ...func(*sqs.Options)) (*sqs.DeleteQueueOutput, error)
	DeleteMessage(ctx context.Context,
		params *sqs.DeleteMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

// CreateQueue Create New Queue
func CreateQueue(c context.Context, api QueueAPI, input *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
	return api.CreateQueue(c, input)
}

// GetQueueURL Get Queue URL
func GetQueueURL(c context.Context, api QueueAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

// ReceiveMessage Receive message from Queue
func ReceiveMessage(c context.Context, api QueueAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
}

// SendMessage Send message to Queue
func SendMessage(c context.Context, api QueueAPI, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return api.SendMessage(c, input)
}

// DeleteQueue delete queue message when acknowledge
func DeleteQueue(c context.Context, api QueueAPI, input *sqs.DeleteQueueInput) (*sqs.DeleteQueueOutput, error) {

	return api.DeleteQueue(c, input)
}

func RemoveMessage(c context.Context, api QueueAPI, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return api.DeleteMessage(c, input)
}
