package stats

type PublicStats struct {
	ActiveMembers   int64 `json:"active_members"`
	ItemsRenewed    int64 `json:"items_renewed"`
	PartnerArtisans int64 `json:"partner_artisans"`
	RegionsCovered  int64 `json:"regions_covered"`
}

type ProStats struct {
	Annonces int `json:"annonces"`
	Projets  int `json:"projets"`
	Vues     int `json:"vues"`
	Score    int `json:"score"`
}
