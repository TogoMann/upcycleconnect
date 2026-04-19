package financial

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
