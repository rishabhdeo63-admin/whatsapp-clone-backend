package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/rishabhdeo60-dev/whatsapp-clone/internal/model"
	"github.com/rishabhdeo60-dev/whatsapp-clone/internal/repository"
)

// MessageService defines the interface for message-related operations
type MessageService interface {
	SendMessage(ctx context.Context, message *model.Message) (*model.Message, error)
	GetConversation(ctx context.Context, senderID, receiverID int64) ([]*model.Message, error)
	UpdateMessageStatus(ctx context.Context, messageID int64, status string) error
}

// messageService is the concrete implementation of MessageService
type messageService struct {
	messageRepo repository.MessageRepository
	userRepo    repository.UserRepository
}

// SendMessage sends a message using the repository
func (s *messageService) SendMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
	if message == nil {
		return nil, nil
	}
	return s.messageRepo.SendMessage(ctx, message)
}

// GetConversation retrieves the conversation between two users
func (s *messageService) GetConversation(ctx context.Context, senderID, receiverID int64) ([]*model.Message, error) {
	_, SenderErr := s.userRepo.FindByID(ctx, senderID)
	_, receiverErr := s.userRepo.FindByID(ctx, receiverID)
	if SenderErr != nil || receiverErr != nil {
		if SenderErr != nil {
			return nil, errors.New(fmt.Errorf("sender with senderID: %d not found", senderID).Error())
		} else {
			return nil, errors.New(fmt.Errorf("user with receiverID: %d not found", receiverID).Error())
		}
	}
	return s.messageRepo.GetConversation(ctx, senderID, receiverID)
}

// UpdateMessageStatus updates the status of a message
func (s *messageService) UpdateMessageStatus(ctx context.Context, messageID int64, status string) error {
	return s.messageRepo.UpdateMessageStatus(ctx, messageID, status)
}

// NewMessageService creates a new instance of MessageService
func NewMessageService(messageRepo repository.MessageRepository, userRepo repository.UserRepository) MessageService {
	return &messageService{messageRepo: messageRepo, userRepo: userRepo}
}
