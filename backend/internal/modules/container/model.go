package container

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ConteneurFrontend struct {
	Id           int64  `db:"id" json:"id"`
	CodeBarres   string `db:"code_barres" json:"code_barres"`
	Localisation string `db:"localisation" json:"localisation"`
	Etat         string `db:"etat" json:"etat"`
	Capacite     int    `db:"capacite" json:"capacite"`
	Objets       int    `db:"objets" json:"objets"`
}

type Container struct {
	Id        pgtype.Int8      `db:"id" json:"id"`
	SiteId    pgtype.Int8      `db:"site_id" json:"site_id"`
	Status    string           `db:"status" json:"status"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
}

type Locker struct {
	Id          pgtype.Int8      `db:"id" json:"id"`
	ContainerId pgtype.Int8      `db:"container_id" json:"container_id"`
	Label       string           `db:"label" json:"label"`
	Status      string           `db:"status" json:"status"`
	Size        string           `db:"size" json:"size"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
}

type LockerAccess struct {
	Id         pgtype.Int8      `db:"id" json:"id"`
	LockerId   pgtype.Int8      `db:"locker_id" json:"locker_id"`
	ItemId     pgtype.Int8      `db:"item_id" json:"item_id"`
	UserId     pgtype.Int8      `db:"user_id" json:"user_id"`
	AccessCode string           `db:"access_code" json:"access_code"`
	ExpiresAt  pgtype.Timestamp `db:"expires_at" json:"expires_at"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
	UsedAt     pgtype.Timestamp `db:"used_at" json:"used_at"`
}

type LockerAccessDetails struct {
	Id         pgtype.Int8      `db:"id" json:"id"`
	LockerId   pgtype.Int8      `db:"locker_id" json:"locker_id"`
	ItemId     pgtype.Int8      `db:"item_id" json:"item_id"`
	UserId     pgtype.Int8      `db:"user_id" json:"user_id"`
	AccessCode string           `db:"access_code" json:"access_code"`
	ExpiresAt  pgtype.Timestamp `db:"expires_at" json:"expires_at"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
	UsedAt     pgtype.Timestamp `db:"used_at" json:"used_at"`

	LockerLabel      string `db:"locker_label" json:"locker_label"`
	ContainerAddress string `db:"container_address" json:"container_address"`
}

type CreateLockerAccessRequest struct {
	ItemId int64 `json:"item_id"`
	UserId int64 `json:"user_id"`
}

type LockerOption struct {
	Id          pgtype.Int8 `db:"id" json:"id"`
	Label       string      `db:"label" json:"label"`
	Size        string      `db:"size" json:"size"`
	Status      string      `db:"status" json:"status"`
	ContainerId pgtype.Int8 `db:"container_id" json:"container_id"`
	SiteId      pgtype.Int8 `db:"site_id" json:"site_id"`
	Address     string      `db:"address" json:"address"`
}

type SiteOption struct {
	SiteId   pgtype.Int8 `db:"site_id" json:"site_id"`
	Address  string      `db:"address" json:"address"`
	TypeSite string      `db:"type_site" json:"type_site"`
}

type SiteWithLockers struct {
	SiteId   pgtype.Int8    `json:"site_id"`
	Address  string         `json:"address"`
	TypeSite string         `json:"type_site"`
	Lockers  []LockerOption `json:"lockers"`
}
