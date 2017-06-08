package models

import (
	"github.com/jinzhu/gorm"
	"golib/modules/logr"
	"time"
)

type DbBase struct {
	//times
	RecUpdTs time.Time
	RecCrtTs time.Time
}

func (t *DbBase) BeforeCreate(scope *gorm.Scope) (err error) {
	logr.Debug("get in BeforeCreate")
	scope.SetColumn("RecUpdTs", time.Now().UTC())
	scope.SetColumn("RecCrtTs", time.Now().UTC())
	return nil
}

func (t *DbBase) BeforeUpdate(scope *gorm.Scope) (err error) {
	logr.Debug("get in BeforeUpdate")
	scope.SetColumn("RecUpdTs", time.Now().UTC())
	return nil
}
