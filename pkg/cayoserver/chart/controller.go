package chart

import (
	"sync"

	"github.com/yunling101/toolkits/date"
)

type Sets struct {
	Value map[string]float64 `json:"value"`
}

type dataSets struct {
	mux     sync.Mutex
	Time    *date.Date
	Results []interface{}      `json:"results"`
	Avail   map[string]float64 `json:"avail"`
	Node    int                `json:"node"`
	Labels  []string           `json:"labels"`
	Keys    []string           `json:"keys"`
	Data    map[string]Sets    `json:"data"`
}

func NewVariable() *dataSets {
	avail := make(map[string]float64, 0)
	data := make(map[string]Sets, 0)
	labels := make([]string, 0)
	results := make([]interface{}, 0)
	return &dataSets{
		Avail: avail, Data: data, Labels: labels, Results: results,
		Time: date.New(),
	}
}

func (d *dataSets) AddAvail(k string, c int) {
	d.mux.Lock()
	defer d.mux.Unlock()

	if c == 1 {
		if _, ok := d.Avail[k]; ok {
			d.Avail[k] += 1
		} else {
			d.Avail[k] = 1
		}
	} else {
		if v, ok := d.Avail[k]; ok {
			// d.Avail[k] -= 1
			d.Avail[k] = v - float64(c)
		}
	}
}

func (d *dataSets) SetAvail(k string, v float64) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.Avail[k] = v
}

func (d *dataSets) GetAvail(k string) bool {
	d.mux.Lock()
	defer d.mux.Unlock()
	if _, ok := d.Avail[k]; ok {
		return true
	}
	return false
}

func (d *dataSets) AddData(k string, set Sets) {
	d.mux.Lock()
	defer d.mux.Unlock()

	if old, ok := d.Data[k]; ok {
		for n, v := range set.Value {
			if o, ok := old.Value[n]; ok {
				old.Value[n] = o + v
			} else {
				old.Value[n] = v
			}
		}
		d.Data[k] = old
	} else {
		d.Data[k] = set
	}
}
