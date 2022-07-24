package gormutil

import "gorm.io/gorm"

const DefaultLimit = 1000

type LimitAndOffset struct {
	Offset int
	Limit  int
}

func Unpointer(offset, limit *int64) *LimitAndOffset {
	o, l := 0, DefaultLimit
	if offset != nil {
		o = int(*offset)
	}

	if limit != nil {
		l = int(*limit)
	}
	return &LimitAndOffset{Offset: o, Limit: l}
}

// Paginate db.Scopes(Paginate(page,pageSize))..
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 500:
			pageSize = 500
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
