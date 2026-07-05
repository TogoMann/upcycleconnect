package chat

import (
	"backend/internal/modules/course"
	"backend/internal/modules/listing"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo           *Repository
	listingService *listing.Service
	courseService  *course.Service
}

func NewService(repo *Repository, listingService *listing.Service, courseService *course.Service) *Service {
	return &Service{repo: repo, listingService: listingService, courseService: courseService}
}

func (s *Service) SendMessage(senderId int64, req CreateMessageRequest) (*Message, error) {
	var conv *Conversation
	var err error

	if req.ConversationId > 0 {
		conv, err = s.repo.GetConversationById(req.ConversationId)
		if err != nil {
			return nil, fmt.Errorf("conversation not found: %w", err)
		}
		if conv == nil {
			return nil, fmt.Errorf("conversation not found")
		}
		if conv.BuyerId != senderId && conv.SellerId != senderId {
			return nil, fmt.Errorf("unauthorized to send message to this conversation")
		}
	} else if req.CourseId > 0 {
		c, err := s.courseService.GetById(pgtype.Int8{Int64: req.CourseId, Valid: true})
		if err != nil {
			return nil, fmt.Errorf("course not found: %w", err)
		}
		sellerId := c.CreatedBy.Int64

		if senderId != sellerId {
			conv, err = s.repo.GetOrCreateCourseConversation(req.CourseId, senderId, sellerId)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("formateur ne peut pas initier une conversation sans destinataire précis")
		}
	} else {
		l, err := s.listingService.GetById(pgtype.Int8{Int64: req.ListingId, Valid: true})
		if err != nil {
			return nil, fmt.Errorf("listing not found: %w", err)
		}
		sellerId := l.CreatedBy.Int64

		if senderId != sellerId {
			conv, err = s.repo.GetOrCreateConversation(req.ListingId, senderId, sellerId)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("vendeur ne peut pas initier une conversation sans destinataire précis")
		}
	}

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
