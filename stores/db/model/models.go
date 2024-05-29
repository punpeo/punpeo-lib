package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Model struct {
	Id         int `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	CreateTime int `gorm:"autoCreateTime;column:create_time;type:int(11)" json:"create_time"`
	UpdateTime int `gorm:"autoUpdateTime;column:update_time;type:int(11)" json:"update_time"`
}

type Preload struct {
	QueryItem string
	Where     string
}

// 分表（通过取指定数字取模运算）
func TableByNumberRemainder(tableName string, number, tableNumber int64) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tableName = TableByNumberRemainderName(tableName, number, tableNumber)
		return tx.Table(tableName)
	}
}
func TableByNumberRemainderName(tableName string, number, tableNumber int64) string {
	if number > 0 && tableNumber > 0 {
		n := number % tableNumber
		tableName = fmt.Sprintf("%s_%d", tableName, n)
		return tableName
	}
	return tableName
}

// 预加载关联处理
func PreloadListHandle(preloadList []Preload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(preloadList) > 0 {
			for _, p := range preloadList {
				db.Preload(p.QueryItem, p.Where)
			}
		}

		return db
	}
}

// 数据权限限制
func DataPermissions(creators []int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("create_user in (?)", creators)
	}
}

// 分页
func Paginate(page int64, pageSize int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize

		// 不限制分页和页数
		if offset <= 0 || page <= 0 {
			offset = -1
		}

		if pageSize <= 0 {
			pageSize = -1
		}

		return db.Limit(int(pageSize)).Offset(int(offset))
	}
}
