package models

import (
	"errors"

	"gorm.io/gorm"
)

type Plants struct {
	// Primary key with auto-increment
	ID uint `gorm:"primaryKey;autoIncrement"`

	// Required Fields
	Name              string `gorm:"uniqueIndex;not null;size:50"`
	OutdoorSowDate    string `gorm:"not null;type:char(5);check LENGTH(OutdoorSowDate) = 5 AND OutdoorSowDate ~ '^[0-9]{2}-[0-9]{2}'"`
	DaysToGermination int    `gorm:"not null;check DaysToGermination > 0"`
	DaysToHarvest     int    `gorm:"not null;check DaysToHarvest > 0"`
	CanStartIndoors   bool   `gorm:"not null"`

	// Fields required if CanStartIndoors is true and not allowed otherwise
	IndoorSowDate  string `gorm:"type:char(5);check LENGTH(IndoorSowDate) = 5 AND IndoorSowDate ~ '^[0-9]{2}-[0-9]{2}$'"`
	TransplantDate string `gorm:"type:char(5));check LENGTH(TransplantDate) = 5 AND TransplantDate ~ '^[0-9]{4}-[0-9]{2}-[0-9]{2}$'"`
}

// TODO - writes tests to verify this works on updates and inserts
func (p *Plants) BeforeSave(tx *gorm.DB) (err error) {
	if !p.CanStartIndoors && (p.IndoorSowDate != "" || p.TransplantDate != "") {
		return errors.New("cannot define either IndoorSowDate or TransplantDate when CanStartIndoors is false")
	}

	if p.CanStartIndoors && (p.IndoorSowDate == "" || p.TransplantDate == "") {
		return errors.New("IndoorSowDate and TransplantDate must be defined when CanStartIndoors is true")
	}

	return nil
}
