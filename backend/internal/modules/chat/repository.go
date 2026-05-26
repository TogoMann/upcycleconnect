package chat

import (
	db "backend/internal/database"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetOrCreateConversation(listingId, buyerId, sellerId int64) (*Conversation, error) {
	var conv Conversation
	err := r.db.QueryRow(db.Ctx, `
		WITH ins AS (
			INSERT INTO chat_conversation (listing_id, buyer_id, seller_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (listing_id, buyer_id, seller_id) DO UPDATE 
			SET updated_at = NOW()
			RETURNING id, listing_id, buyer_id, seller_id, is_closed, created_at, updated_at
		)
		SELECT ins.id, ins.listing_id, ins.buyer_id, ins.seller_id, ins.is_closed, ins.created_at, ins.updated_at, COALESCE(l.name, '')
		FROM ins
		LEFT JOIN listing l ON ins.listing_id = l.id
	`, listingId, buyerId, sellerId).Scan(&conv.Id, &conv.ListingId, &conv.BuyerId, &conv.SellerId, &conv.IsClosed, &conv.CreatedAt, &conv.UpdatedAt, &conv.ListingTitle)

	if err != nil {
		return nil, fmt.Errorf("package chat/repo GetOrCreateConversation: %w", err)
	}

	return &conv, nil
}

func (r *Repository) GetConversationByListingAndUsers(listingId, userId1, userId2 int64) (*Conversation, error) {
	var conv Conversation
	err := r.db.QueryRow(db.Ctx, `
		SELECT id, listing_id, buyer_id, seller_id, is_closed, created_at, updated_at
		FROM chat_conversation
		WHERE listing_id = $1 AND (
			(buyer_id = $2 AND seller_id = $3) OR
			(buyer_id = $3 AND seller_id = $2)
		)
	`, listingId, userId1, userId2).Scan(&conv.Id, &conv.ListingId, &conv.BuyerId, &conv.SellerId, &conv.IsClosed, &conv.CreatedAt, &conv.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("package chat/repo GetConversationByListingAndUsers: %w", err)
	}

	return &conv, nil
}

func (r *Repository) GetConversationById(id int64) (*Conversation, error) {
	var conv Conversation
	err := r.db.QueryRow(db.Ctx, `
		SELECT id, listing_id, buyer_id, seller_id, is_closed, created_at, updated_at
		FROM chat_conversation
		WHERE id = $1
	`, id).Scan(&conv.Id, &conv.ListingId, &conv.BuyerId, &conv.SellerId, &conv.IsClosed, &conv.CreatedAt, &conv.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("package chat/repo GetConversationById: %w", err)
	}

	return &conv, nil
}

func (r *Repository) CreateMessage(msg Message) (*Message, error) {
	var newMsg Message
	err := r.db.QueryRow(db.Ctx, `
		INSERT INTO chat_message (conversation_id, sender_id, content, message_type, proposed_price, proposal_status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, conversation_id, sender_id, content, message_type, proposed_price, proposal_status, created_at
	`, msg.ConversationId, msg.SenderId, msg.Content, msg.MessageType, msg.ProposedPrice, msg.ProposalStatus).Scan(
		&newMsg.Id, &newMsg.ConversationId, &newMsg.SenderId, &newMsg.Content, &newMsg.MessageType, &newMsg.ProposedPrice, &newMsg.ProposalStatus, &newMsg.CreatedAt,
	)
	newMsg.IsEdited = false

	if err != nil {
		return nil, fmt.Errorf("package chat/repo CreateMessage: %w", err)
	}

	_, _ = r.db.Exec(db.Ctx, "UPDATE chat_conversation SET updated_at = NOW() WHERE id = $1", msg.ConversationId)

	return &newMsg, nil
}

func (r *Repository) GetMessagesByConversationId(convId int64) ([]Message, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			m.id, m.conversation_id, m.sender_id, m.content, m.message_type, 
			m.proposed_price, m.proposal_status, m.created_at,
			EXISTS (SELECT 1 FROM chat_message_edit_history h WHERE h.message_id = m.id) as is_edited
		FROM chat_message m
		WHERE m.conversation_id = $1
		ORDER BY m.created_at ASC
	`, convId)
	if err != nil {
		return nil, fmt.Errorf("package chat/repo GetMessagesByConversationId query: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Message])
}

func (r *Repository) GetMessageById(id int64) (*Message, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			m.id, m.conversation_id, m.sender_id, m.content, m.message_type, 
			m.proposed_price, m.proposal_status, m.created_at,
			EXISTS (SELECT 1 FROM chat_message_edit_history h WHERE h.message_id = m.id) as is_edited
		FROM chat_message m
		WHERE m.id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("package chat/repo GetMessageById query: %w", err)
	}

	msg, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Message])
	if err != nil {
		return nil, fmt.Errorf("package chat/repo GetMessageById: %w", err)
	}
	return &msg, nil
}

func (r *Repository) UpdateMessage(id int64, oldContent string, oldPrice *float64, newContent string, newPrice *float64) error {
	tx, err := r.db.Begin(db.Ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(db.Ctx)

	_, err = tx.Exec(db.Ctx, `
		INSERT INTO chat_message_edit_history (message_id, old_content, old_proposed_price)
		VALUES ($1, $2, $3)
	`, id, oldContent, oldPrice)
	if err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}

	_, err = tx.Exec(db.Ctx, `
		UPDATE chat_message
		SET content = $1, proposed_price = $2
		WHERE id = $3
	`, newContent, newPrice, id)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	return tx.Commit(db.Ctx)
}

func (r *Repository) GetConversationsByUserId(userId int64) ([]Conversation, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			c.id, c.listing_id, c.buyer_id, c.seller_id, c.is_closed, c.created_at, c.updated_at,
			COALESCE(l.name, '') as listing_title
		FROM chat_conversation c
		LEFT JOIN listing l ON c.listing_id = l.id
		WHERE c.buyer_id = $1 OR c.seller_id = $1
		ORDER BY c.updated_at DESC
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("package chat/repo GetConversationsByUserId query: %w", err)
	}
	defer rows.Close()

	conversations := []Conversation{}
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(
			&conv.Id, &conv.ListingId, &conv.BuyerId, &conv.SellerId, &conv.IsClosed, &conv.CreatedAt, &conv.UpdatedAt,
			&conv.ListingTitle,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func (r *Repository) CloseConversationsByListingId(listingId int64) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE chat_conversation SET is_closed = true WHERE listing_id = $1", listingId)
	return err
}

func (r *Repository) GetAdminConversationReviews() ([]AdminConversationReview, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			c.id, c.listing_id, c.buyer_id, c.seller_id, c.created_at, c.updated_at,
			l.name as listing_title,
			u1.username as buyer_name,
			u2.username as seller_name
		FROM chat_conversation c
		JOIN listing l ON c.listing_id = l.id
		JOIN users u1 ON c.buyer_id = u1.id
		JOIN users u2 ON c.seller_id = u2.id
		ORDER BY c.updated_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("package chat/repo GetAdminConversationReviews query: %w", err)
	}
	defer rows.Close()

	var reviews []AdminConversationReview
	for rows.Next() {
		var rev AdminConversationReview
		err := rows.Scan(
			&rev.Id, &rev.ListingId, &rev.BuyerId, &rev.SellerId, &rev.CreatedAt, &rev.UpdatedAt,
			&rev.ListingTitle, &rev.BuyerName, &rev.SellerName,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}

	return reviews, nil
}

func (r *Repository) GetMessageEditHistory(messageId int64) ([]MessageEditHistory, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT id, message_id, old_content, old_proposed_price, edited_at
		FROM chat_message_edit_history
		WHERE message_id = $1
		ORDER BY edited_at DESC
	`, messageId)
	if err != nil {
		return nil, fmt.Errorf("package chat/repo GetMessageEditHistory query: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[MessageEditHistory])
}

func (r *Repository) UpdateProposalStatus(messageId int64, status ProposalStatus) error {
	_, err := r.db.Exec(db.Ctx, `
		UPDATE chat_message
		SET proposal_status = $1
		WHERE id = $2 AND message_type = 'price_proposal'
	`, status, messageId)
	return err
}
