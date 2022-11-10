package task

import (
	"encoding/json"

	"github.com/yunling101/cayo/pkg/model/q"
)

type Factory interface {
	Save() (int, error)
	Query() ([]int, error)
	Delete(selector interface{}) error
}

type modal struct {
	Relation Factory
	queue    []int
	column   string
}

func NewModal(val Factory) *modal {
	return &modal{Relation: val}
}

func (v *modal) Add(in []Factory, out interface{}) error {
	var result []int
	for i := 0; i < len(in); i++ {
		if i == 0 {
			if b, err := in[i].Query(); err != nil {
				return err
			} else {
				v.queue = b
			}
		}
		if id, err := in[i].Save(); err != nil {
			return err
		} else {
			result = append(result, id)
		}
	}
	if err := v.contrast(result); err != nil {
		return err
	}
	b, _ := json.Marshal(result)
	return json.Unmarshal(b, out)
}

func (v *modal) contrast(b []int) error {
	for i := 0; i < len(v.queue); i++ {
		if !q.Contains(b, v.queue[i]) {
			if err := v.Relation.Delete(q.M{"id": v.queue[i]}); err != nil {
				return err
			}
		}
	}
	return nil
}
