package item

import (
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type ItemStatus string

const (
	Deposited ItemStatus = "deposited"
	Validated ItemStatus = "validated"
	Collected ItemStatus = "collected"
)

type ItemState string

const (
	Neuf    ItemState = "neuf"
	BonEtat ItemState = "bon etat"
	Abime   ItemState = "abime"
	Casse   ItemState = "casse"
)

func NormalizeState(input string) ItemState {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "neuf":
		return Neuf
	case "bon état", "bon etat":
		return BonEtat
	case "abimé", "abime":
		return Abime
	case "cassé", "casse":
		return Casse
	default:
		return BonEtat
	}
}

type ItemSize string

const (
	SizeS ItemSize = "S"
	SizeM ItemSize = "M"
	SizeL ItemSize = "L"
)

var SizeRank = map[ItemSize]int{
	SizeS: 1,
	SizeM: 2,
	SizeL: 3,
}

type Item struct {
	Id            pgtype.Int8      `db:"id" json:"id"`
	OwnerId       pgtype.Int8      `db:"owner_id" json:"owner_id"`
	LockerId      pgtype.Int8      `db:"locker_id" json:"locker_id"`
	SiteId        pgtype.Int8      `db:"site_id" json:"site_id"`
	MaterialType  string           `db:"material_type" json:"material_type"`
	PhysicalState ItemState        `db:"physical_state" json:"physical_state"`
	Size          ItemSize         `db:"size" json:"size"`
	Status        ItemStatus       `db:"status" json:"status"`
	Weight        pgtype.Numeric   `db:"weight" json:"weight"`
	Name          string           `db:"name" json:"name"`
	Description   string           `db:"description" json:"description"`
	EntryId       pgtype.Int8      `db:"entry_id" json:"entry_id"`
	ScheduleDate  string           `db:"schedule_date" json:"schedule_date"`
	ScheduleTime  string           `db:"schedule_time" json:"schedule_time"`
	CreatedAt     pgtype.Timestamp `db:"created_at" json:"created_at"`
	SiteType      string           `db:"site_type" json:"site_type"`
}

type AdminDepot struct {
	Id          int64  `db:"id" json:"id"`
	Utilisateur string `db:"utilisateur" json:"utilisateur"`
	Objet       string `db:"objet" json:"objet"`
	Date        string `db:"date" json:"date"`
	Statut      string `db:"statut" json:"statut"`
	CodeEnvoye  bool   `db:"code_envoye" json:"code_envoye"`
}
