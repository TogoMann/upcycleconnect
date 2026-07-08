package payments

import (
	"backend/internal/modules/advertisement"
	"backend/internal/modules/listing"
	listingorder "backend/internal/modules/listing_order"
	"backend/internal/modules/plans"
	"backend/internal/modules/subscriptions"
	"backend/internal/utils"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/webhook"
)

type Service struct {
	listingService       *listing.Service
	listingOrderService  *listingorder.Service
	advertisementService *advertisement.Service
	planService          *plans.Service
	subscriptionService  *subscriptions.Service
	frontendURL          string
	webhookSecret        string
}

func NewService(listingService *listing.Service, listingOrderService *listingorder.Service, advertisementService *advertisement.Service, planService *plans.Service, subscriptionService *subscriptions.Service) *Service {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	return &Service{
		listingService:       listingService,
		listingOrderService:  listingOrderService,
		advertisementService: advertisementService,
		planService:          planService,
		subscriptionService:  subscriptionService,
		frontendURL:          frontendURL,
		webhookSecret:        os.Getenv("STRIPE_WEBHOOK_SECRET"),
	}
}

func (s *Service) CreateSubscriptionCheckout(userId int64, planId int64, siret string, returnPath string) (string, error) {
	siret = utils.CleanSiret(siret)
	p, err := s.planService.GetById(pgtype.Int8{Int64: planId, Valid: true})
	if err != nil {
		return "", fmt.Errorf("plan introuvable: %w", err)
	}

	price, err := p.Price.Float64Value()
	if err != nil || price.Float64 <= 0 {
		return "", fmt.Errorf("ce plan ne nécessite pas de paiement")
	}

	if returnPath == "" {
		returnPath = "/particulier/plans"
	}

	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(s.frontendURL + returnPath + "?checkout=success&session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(s.frontendURL + returnPath + "?checkout=cancel"),
		Metadata: map[string]string{
			"type":     "subscription",
			"user_id":  fmt.Sprintf("%d", userId),
			"plan_id":  fmt.Sprintf("%d", planId),
			"siret":    siret,
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String("eur"),
					UnitAmount: stripe.Int64(int64(price.Float64 * 100)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Abonnement " + p.Name),
					},
				},
			},
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return "", err
	}
	return sess.URL, nil
}

func (s *Service) CreateAdvertisementCheckout(userId int64, adId int64) (string, error) {
	ad, err := s.advertisementService.GetById(pgtype.Int8{Int64: adId, Valid: true})
	if err != nil {
		return "", fmt.Errorf("publicité introuvable: %w", err)
	}

	if ad.AnnouncerId.Int64 != userId {
		return "", fmt.Errorf("vous n'êtes pas le propriétaire de cette publicité")
	}

	budget, err := ad.Budget.Float64Value()
	if err != nil || budget.Float64 <= 0 {
		return "", fmt.Errorf("budget de campagne invalide")
	}

	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(s.frontendURL + "/pro/publicites?checkout=success&session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(s.frontendURL + "/pro/publicites?checkout=cancel"),
		Metadata: map[string]string{
			"type": "advertisement",
			"ad_id": fmt.Sprintf("%d", adId),
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String("eur"),
					UnitAmount: stripe.Int64(int64(budget.Float64 * 100)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Campagne publicitaire #" + fmt.Sprintf("%d", adId)),
					},
				},
			},
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return "", err
	}
	return sess.URL, nil
}

type ListingCheckoutResult struct {
	Free    bool
	OrderId int64
	URL     string
}

func (s *Service) CreateListingOrderCheckout(userId int64, listingId int64) (*ListingCheckoutResult, error) {
	l, err := s.listingService.GetById(pgtype.Int8{Int64: listingId, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("annonce introuvable: %w", err)
	}

	priceVal, err := l.Price.Float64Value()
	if err != nil {
		return nil, fmt.Errorf("prix invalide")
	}
	price := priceVal.Float64

	negotiatedPrice, hasNegotiation, err := s.listingService.GetNegotiatedPrice(listingId, userId)
	if err == nil && hasNegotiation {
		price = negotiatedPrice
	}

	if price <= 0 {
		id, err := s.listingOrderService.CreateFromRequest(userId, listingorder.CreateListingOrderRequest{
			ListingId: listingId,
			Price:     0,
		})
		if err != nil {
			return nil, err
		}
		return &ListingCheckoutResult{Free: true, OrderId: id.Int64}, nil
	}

	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(s.frontendURL + "/particulier/paiement/confirmation?type=listing&listing_id=" + fmt.Sprintf("%d", listingId) + "&session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(s.frontendURL + "/annonces/" + fmt.Sprintf("%d", listingId)),
		Metadata: map[string]string{
			"type":       "listing_order",
			"user_id":    fmt.Sprintf("%d", userId),
			"listing_id": fmt.Sprintf("%d", listingId),
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String("eur"),
					UnitAmount: stripe.Int64(int64(price * 100)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(l.Name),
					},
				},
			},
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, err
	}
	return &ListingCheckoutResult{Free: false, URL: sess.URL}, nil
}

type VerifyResult struct {
	Paid bool
	Type string
}

func (s *Service) VerifySession(sessionId string) (*VerifyResult, error) {
	sess, err := session.Get(sessionId, nil)
	if err != nil {
		return nil, err
	}

	return &VerifyResult{
		Paid: sess.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid,
		Type: sess.Metadata["type"],
	}, nil
}

func (s *Service) HandleWebhook(payload []byte, signatureHeader string) error {
	event, err := webhook.ConstructEvent(payload, signatureHeader, s.webhookSecret)
	if err != nil {
		return fmt.Errorf("signature webhook invalide: %w", err)
	}

	if event.Type != "checkout.session.completed" {
		return nil
	}

	var sess stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
		return fmt.Errorf("payload webhook invalide: %w", err)
	}

	if sess.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
		return nil
	}

	switch sess.Metadata["type"] {
	case "subscription":
		return s.handleSubscriptionPaid(sess)
	case "advertisement":
		return s.handleAdvertisementPaid(sess)
	case "listing_order":
		return s.handleListingOrderPaid(sess)
	}

	return nil
}

func (s *Service) handleSubscriptionPaid(sess stripe.CheckoutSession) error {
	var userId, planId int64
	fmt.Sscanf(sess.Metadata["user_id"], "%d", &userId)
	fmt.Sscanf(sess.Metadata["plan_id"], "%d", &planId)
	siret := sess.Metadata["siret"]

	return s.subscriptionService.ChoosePlan(userId, planId, siret)
}

func (s *Service) handleAdvertisementPaid(sess stripe.CheckoutSession) error {
	var adId int64
	fmt.Sscanf(sess.Metadata["ad_id"], "%d", &adId)

	return s.advertisementService.SetStripePaymentIntentId(pgtype.Int8{Int64: adId, Valid: true}, sess.ID)
}

func (s *Service) handleListingOrderPaid(sess stripe.CheckoutSession) error {
	var userId, listingId int64
	fmt.Sscanf(sess.Metadata["user_id"], "%d", &userId)
	fmt.Sscanf(sess.Metadata["listing_id"], "%d", &listingId)

	l, err := s.listingService.GetById(pgtype.Int8{Int64: listingId, Valid: true})
	if err != nil {
		return err
	}

	if l.Status == listing.Sold {
		return nil
	}

	_, err = s.listingOrderService.CreatePaidOrder(userId, listingId)
	return err
}
