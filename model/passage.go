package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/Eraac/train-sniffer/utils"
)

type (
	Passage struct {
		gorm.Model
		Time      time.Time
		IsWeek    bool
		Station   Station
		StationID uint
		Train     Train
		TrainID   uint
	}

	passageRepository struct{}
)

var PassageRepository *passageRepository

func (r *passageRepository) IsExist(code string, s *Station) bool {
	key := cache.buildKeyPassage(*s, code)

	exist := cache.IsKeyExist(key)

	if exist {
		return true
	}

	count := 0

	err := db.Model(&Passage{}).Joins(
		"LEFT JOIN train ON train.id = train_id",
			).Where("passage.station_id = ? AND train.code = ?", s.ID, code).Count(&count).Error

	if err != nil {
		utils.Error(err.Error())
		return false
	}

	if count > 0 {
		err = cache.set(key, true)

		if err != nil {
			utils.Error(err.Error())
		}
	}

	return count > 0
}

func (p *Passage) Persist() error {
	return db.Save(p).Error
}

func (p *Passage) String() string {
	return fmt.Sprintf("[passage] id: %d - station_id: %d - train_id: %d", p.ID, p.StationID, p.TrainID)
}
