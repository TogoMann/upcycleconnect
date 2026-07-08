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
	ActionListingCreated = PointAction{
		Points:      5,
		Description: "Annonce publiée",
	}
)

const PointsPerKg = 10

type Quest struct {
	ActionDescription string
	WindowDays        int
	Threshold         int
	BonusPoints       int32
	BonusDescription  string
}

var Quests = []Quest{
	{
		ActionDescription: ActionListingCreated.Description,
		WindowDays:        7,
		Threshold:         5,
		BonusPoints:       50,
		BonusDescription:  "Défi hebdomadaire : 5 annonces publiées",
	},
	{
		ActionDescription: ActionSaleCompleted.Description,
		WindowDays:        7,
		Threshold:         3,
		BonusPoints:       75,
		BonusDescription:  "Défi hebdomadaire : 3 ventes finalisées",
	},
	{
		ActionDescription: ActionDepositValidated.Description,
		WindowDays:        7,
		Threshold:         3,
		BonusPoints:       30,
		BonusDescription:  "Défi hebdomadaire : 3 dépôts validés",
	},
}

func LevelFromScore(score int32) int32 {
	if score < 0 {
		score = 0
	}
	return score/100 + 1
}
