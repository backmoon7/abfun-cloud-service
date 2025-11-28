package model

import "gorm.io/gorm"

type Danmaku struct {
	gorm.Model
	VideoID uint   `gorm:"index"`
	Content string `gorm:"size:255"`
	Time    float64
	Color   string `gorm:"size:10"`
	Type    int
	AuthorID uint
}
