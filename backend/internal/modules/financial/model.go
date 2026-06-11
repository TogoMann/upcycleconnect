package financial

import "github.com/jackc/pgx/v5/pgtype"

type Commission struct {
	Id             int64   `db:"id" json:"id"`
	Type           string  `db:"type" json:"type"`
	Taux           float64 `db:"taux" json:"taux"`
	MontantTotal   float64 `db:"montant_total" json:"montant_total"`
	NbTransactions int     `db:"nb_transactions" json:"nb_transactions"`
	Periode        string  `db:"periode" json:"periode"`
}

type Evolution struct {
	Mois    string  `json:"mois"`
	Ca      float64 `json:"ca"`
	Charges float64 `json:"charges"`
}

type FinancierData struct {
	CaTotal   float64     `json:"ca_total"`
	CaMois    float64     `json:"ca_mois"`
	Charges   float64     `json:"charges"`
	Marge     float64     `json:"marge"`
	Evolution []Evolution `json:"evolution"`
}

type FinancialReport struct {
	TotalRevenue        float64 `json:"total_revenue"`
	ListingRevenue      float64 `json:"listing_revenue"`
	CourseRevenue       float64 `json:"course_revenue"`
	SubscriptionRevenue float64 `json:"subscription_revenue"`
	AdRevenue           float64 `json:"ad_revenue"`
	TotalCommissions    float64 `json:"total_commissions"`
}

type Invoice struct {
	Id            int64       `db:"id" json:"id"`
	UserId        int64       `db:"user_id" json:"user_id"`
	SellerId      pgtype.Int8 `db:"seller_id" json:"seller_id"`
	OrderId       int64       `db:"order_id" json:"order_id"`
	OrderType     string      `db:"order_type" json:"order_type"`
	InvoiceNumber string      `db:"invoice_number" json:"invoice_number"`
	Amount        float64     `db:"amount" json:"amount"`
	VatAmount     float64     `db:"vat_amount" json:"vat_amount"`
	TotalAmount   float64     `db:"total_amount" json:"total_amount"`
	Status        string      `db:"status" json:"status"`
	CreatedAt     string      `db:"created_at" json:"created_at"`
}

type InvoiceDetail struct {
	Invoice
	UserName    string `json:"user_name"`
	UserEmail   string `json:"user_email"`
	UserAddress string `json:"user_address"`
	Description string `json:"description"`
}
