package q

import (
	"fmt"

	"github.com/yunling101/cayo/pkg/global"
)

type M map[string]interface{}

type model struct {
	tableName string
}

func Table(tableName string) *model {
	return &model{tableName: tableName}
}

// QueryCount
func (m *model) QueryExist(selector interface{}, desc string) error {
	if count, err := global.DB.C(m.tableName).Find(selector).Count(); err != nil {
		return err
	} else {
		if count != 0 {
			return fmt.Errorf("%s!", desc+"已存在")
		}
	}
	return nil
}

// QueryCount
func (m *model) QueryCount(selector interface{}) (int, error) {
	return global.DB.C(m.tableName).Find(selector).Count()
}

// QueryMany
func (m *model) QueryMany(selector, result interface{}) error {
	return global.DB.C(m.tableName).Find(selector).All(result)
}

// QueryPage
func (m *model) QueryPage(limit int, offset int, selector, result interface{}) error {
	return global.DB.C(m.tableName).Find(selector).Limit(limit).Skip(offset).All(result)
}

// DeleteOne
func (m *model) DeleteOne(selector interface{}) error {
	return global.DB.C(m.tableName).Remove(selector)
}

// QueryOne
func (m *model) QueryOne(selector, result interface{}) error {
	return global.DB.C(m.tableName).Find(selector).One(result)
}

// InsertOne
func (m *model) InsertOne(docs interface{}) error {
	return global.DB.C(m.tableName).Insert(&docs)
}

// UpdateOne
func (m *model) UpdateOne(selector, update interface{}) error {
	return global.DB.C(m.tableName).Update(selector, M{"$set": update})
}

// PullOne
func (m *model) PullOne(selector, update interface{}) error {
	return global.DB.C(m.tableName).Update(selector, M{"$pull": update})
}

// PushOne
func (m *model) PushOne(selector, update interface{}) error {
	return global.DB.C(m.tableName).Update(selector, M{"$push": update})
}

// IncOne
func (m *model) IncOne(selector, update interface{}) error {
	return global.DB.C(m.tableName).Update(selector, M{"$inc": update})
}

// Contains
func Contains(array []int, id int) bool {
	for i := 0; i < len(array); i++ {
		if array[i] == id {
			return true
		}
	}
	return false
}
