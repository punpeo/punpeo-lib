package db

import (
	"database/sql"
	"time"
)

type LogErr struct {
	Id         uint           `db:"id"`
	Caller     string         `db:"caller" gorm:"not null;default:'';type:varchar(200)"`
	Content    sql.NullString `db:"content" gorm:"type:text"`
	Level      string         `db:"level" gorm:"not null;default:'';type:varchar(50)"`
	Span       string         `db:"span" gorm:"index:span;:not null;default:'';type:varchar(200)"`
	Trace      string         `db:"trace" gorm:"index:trace;not null;default:'';type:varchar(200)"`
	CreateTime time.Time      `db:"create_time" gorm:"index:create_time"`
}
