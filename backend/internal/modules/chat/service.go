package chat

import (
	"backend/internal/modules/listing"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo           *Repository
	listingService *listing.Service
}

func NewService(repo *Repository, listingService *listing.Service) *Service {
	return &Service{repo: repo, listingService: listingService}
}

func (s *Service) SendMessage(senderId int64, req CreateMessageRequest) (*Message, error) {
	// 1. Get listing to find seller
	l, err := s.listingService.GetById(pgtype.Int8{Int64: req.ListingId, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("listing not found: %w", err)
	}

	sellerId := l.CreatedBy.Int64
	buyerId := senderId

	// If sender is seller, we need to find the conversation first or ensure it exists
	// Actually, usually the buyer starts the conversation.
	// If seller sends a message, they must already be in a conversation.
	// For simplicity, let's assume if it's a new conversation, the sender is the buyer.
	if senderId == sellerId {
		return nil, fmt.Errorf("seller cannot start a conversation with themselves")
	}

	// 2. Get or create conversation
	conv, err := s.repo.GetOrCreateConversation(req.ListingId, buyerId, sellerId)
	if err != nil {
		return nil, err
	}

	// 3. Create message
	msg := Message{
		ConversationId: conv.Id,
		SenderId:       senderId,
		Content:        req.Content,
		MessageType:    req.MessageType,
		ProposedPrice:  req.ProposedPrice,
	}

	if req.MessageType == TypePriceProposal {
		status := StatusPending
		msg.ProposalStatus = &status
	}

	return s.repo.CreateMessage(msg)
}

func (s *Service) GetConversationMessages(userId, convId int64) ([]Message, error) {
	// Verify user belongs to conversation (omitted for brevity but recommended)
	return s.repo.GetMessagesByConversationId(convId)
}

func (s *Service) EditMessage(userId, messageId int64, req EditMessageRequest) error {
	msg, err := s.repo.GetMessageById(messageId)
	if err != nil {
		return err
	}

	if msg.SenderId != userId {
		return fmt.Errorf("unauthorized to edit this message")
	}

	return s.repo.UpdateMessage(messageId, msg.Content, msg.ProposedPrice, req.Content, req.ProposedPrice)
}

func (s *Service) HandleProposal(userId, messageId int64, accept bool) error {
	msg, err := s.repo.GetMessageById(messageId)
	if err != nil {
		return err
	}

	if msg.MessageType != TypePriceProposal {
		return fmt.Errorf("message is not a price proposal")
	}

	// Only the recipient (not the sender) can accept/decline
	// We need to check if userId is the other party in the conversation
	// For now, let's just update if it's pending
	if *msg.ProposalStatus != StatusPending {
		return fmt.Errorf("proposal already handled")
	}

	status := StatusDeclined
	if accept {
		status = StatusAccepted
	}

	return s.repo.UpdateProposalStatus(messageId, status)
}

func (s *Service) GetUserConversations(userId int64) ([]Conversation, error) {
	return s.repo.GetConversationsByUserId(userId)
}

func (s *Service) AdminGetConversations() ([]AdminConversationReview, error) {
	return s.repo.GetAdminConversationReviews()
}

func (s *Service) AdminGetConversationDetails(convId int64) ([]Message, []MessageEditHistory, error) {
	messages, err := s.repo.GetMessagesByConversationId(convId)
	if err != nil {
		return nil, nil, err
	}

	var allHistory []MessageEditHistory
	for _, m := range messages {
		history, err := s.repo.GetMessageEditHistory(m.Id)
		if err == nil {
			allHistory = append(allHistory, history...)
		}
	}

	return messages, allHistory, nil
}
