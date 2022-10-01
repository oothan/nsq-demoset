package criteria

import "gorm.io/gorm"

type Criteria interface {
	Build(db *gorm.DB) (*gorm.DB, error)
}

type set []Criteria

func (c set) Build(db *gorm.DB) (*gorm.DB, error) {
	for _, it := range c {
		db, err := it.Build(db)
		if err != nil {
			return db, err
		}
	}
	return db, nil
}

func Set(crits ...Criteria) Criteria {
	return set(crits)
}

type ID struct {
	ID uint64
}

func (c *ID) Build(db *gorm.DB) (*gorm.DB, error) {
	if c.ID <= 0 {
		return db, nil
	}
	return db.Where("id = ?", c.ID), nil
}

type IDs struct {
	IDs []uint64
}

func (c *IDs) Build(db *gorm.DB) (*gorm.DB, error) {
	if len(c.IDs) == 0 {
		return db, nil
	}
	return db.Where("id IN ?", c.IDs), nil
}

type UserID struct {
	UserID uint64
}

func (c *UserID) Build(db *gorm.DB) (*gorm.DB, error) {
	if c.UserID == 0 {
		return db, nil
	}
	return db.Where("user_id = ?", c.UserID), nil
}

type FromId struct {
	FromId uint64
}

func (c *FromId) Build(db *gorm.DB) (*gorm.DB, error) {
	if c.FromId == 0 {
		return db, nil
	}
	return db.Where("from_id = ?", c.FromId), nil
}

type Paging struct {
	Page     int
	PageSize int
}

func (c *Paging) Build(db *gorm.DB) (*gorm.DB, error) {
	if c.Page != 0 && c.PageSize != 0 {
		offset := (c.Page - 1) * c.PageSize
		return db.Offset(offset).Limit(c.PageSize), nil
	}
	return db, nil
}
