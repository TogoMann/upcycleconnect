package chat

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type MessageType string

const (
	TypeText          MessageType = "text"
	TypePriceProposal MessageType = "price_proposal"
)

type ProposalStatus string

const (
	StatusPending  ProposalStatus = "pending"
	StatusAccepted ProposalStatus = "accepted"
	StatusDeclined ProposalStatus = "declined"
)

type Conversation struct {
	Id           int64            `db:"id" json:"id"`
	ListingId    int64            `db:"listing_id" json:"listing_id"`
	BuyerId      int64            `db:"buyer_id" json:"buyer_id"`
	SellerId     int64            `db:"seller_id" json:"seller_id"`
	IsClosed     bool             `db:"is_closed" json:"is_closed"`
	CreatedAt    pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt    pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	ListingTitle string           `json:"listing_title"`
}

type Message struct {
	Id             int64            `db:"id" json:"id"`
	ConversationId int64            `db:"conversation_id" json:"conversation_id"`
	SenderId       int64            `db:"sender_id" json:"sender_id"`
	Content        string           `db:"content" json:"content"`
	MessageType    MessageType      `db:"message_type" json:"message_type"`
	ProposedPrice  *float64         `db:"proposed_price" json:"proposed_price"`
	ProposalStatus *ProposalStatus  `db:"proposal_status" json:"proposal_status"`
	CreatedAt      pgtype.Timestamp `db:"created_at" json:"created_at"`
	IsEdited       bool             `db:"is_edited" json:"is_edited"`
}

type MessageEditHistory struct {
	Id               int64            `db:"id" json:"id"`
	MessageId        int64            `db:"message_id" json:"message_id"`
	OldContent       string           `db:"old_content" json:"old_content"`
	OldProposedPrice *float64         `db:"old_proposed_price" json:"old_proposed_price"`
	EditedAt         pgtype.Timestamp `db:"edited_at" json:"edited_at"`
}

type CreateMessageRequest struct {
	ListingId      int64       `json:"listing_id"`
	ConversationId int64       `json:"conversation_id"`
	Content        string      `json:"content"`
	MessageType    MessageType `json:"message_type"`
	ProposedPrice  *float64    `json:"proposed_price"`
}

type EditMessageRequest struct {
	Content       string   `json:"content"`
	ProposedPrice *float64 `json:"proposed_price"`
}

type AdminConversationReview struct {
	Conversation
	ListingTitle string    `json:"listing_title"`
	BuyerName    string    `json:"buyer_name"`
	SellerName   string    `json:"seller_name"`
	Messages     []Message `json:"messages"`
}
