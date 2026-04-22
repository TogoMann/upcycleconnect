package utils

type PointAction struct {
	Points      int32
	Description string
}

var (
	ActionRegistration = PointAction{
		Points:      10,
		Description: "Inscription plateforme",
	}
	ActionCommentUseful = PointAction{
		Points:      5,
		Description: "Commentaire utile",
	}
	ActionDepositValidated = PointAction{
		Points:      20,
		Description: "Dépôt validé",
	}
	ActionSaleCompleted = PointAction{
		Points:      50,
		Description: "Vente finalisée",
	}
	ActionAtelierParticipation = PointAction{
		Points:      15,
		Description: "Participation atelier",
	}
	ActionMaterialCollection = PointAction{
		Points:      25,
		Description: "Collecte matériaux",
	}
)
