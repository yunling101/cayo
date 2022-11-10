package chart

import (
	"strconv"

	"github.com/yunling101/cayo/pkg/model/alarm"
)

type H map[string]interface{}

var options = H{
	"chart": H{
		"zoom": H{
			"type":           "x",
			"enabled":        false,
			"autoScaleYaxis": true,
		},
		"toolbar": H{"show": false},
	},
	"fill":       H{"opacity": 1},
	"dataLabels": H{"enabled": false},
}

func Get(n *dataSets, rules []alarm.AlarmRule) H {
	newOptions := options
	chart := newOptions["chart"].(H)
	chart["sparkline"] = H{"enabled": true}
	newOptions["chart"] = chart
	newOptions["labels"] = n.Labels
	newOptions["colors"] = []string{"#5564BE"}
	newOptions["title"] = H{
		"text": "正常",
		"style": H{
			"fontSize": "13px",
			"color":    "#77aa60",
		},
	}
	newOptions["yaxis"] = H{"show": false}
	newOptions["subtitle"] = H{
		"text":    "可用率",
		"offsetY": 25,
	}
	data1 := make([]int, 0)
	data2 := make([]H, 0)

	if len(n.Keys) == len(n.Avail) {
		for _, v := range n.Keys {
			data1 = append(data1, int(n.Avail[v]/float64(n.Node)*100))
		}
	}

	if len(data1) == len(n.Labels) {
		for i := 0; i < len(data1); i++ {
			data2 = append(data2, H{"x": n.Labels[i], "y": data1[i]})
		}
	}
	series1 := []interface{}{
		H{"name": "可用率", "data": data1},
	}
	series2 := []interface{}{
		H{"name": "可用率", "data": data2},
	}
	var series3 []interface{}

	var (
		availablePercent int = 0
		responseTime     int = 0
		max              int = 0
	)
	if len(n.Data) == len(n.Labels) && len(n.Data) == len(n.Keys) {
		index := 0
		reset := make(map[string][]H)
		order := make([]string, 0)
		for _, v := range n.Keys {
			// for _, v := range n.Data {
			for c, x := range n.Data[v].Value {
				y := int(x / float64(n.Node))
				if max < y {
					max = y
				}
				if b, ok := reset[c]; ok {
					b = append(b, H{"x": n.Labels[index], "y": y})
					reset[c] = b
				} else {
					reset[c] = []H{{"x": n.Labels[index], "y": y}}
				}
				//
				if !isExist(order, c) {
					order = append(order, c)
				}
			}
			index++
		}
		for _, v := range order {
			// for k, v := range reset {
			series3 = append(series3, H{"name": v, "data": reset[v]})
		}
		if max != 0 {
			max = int((float64(max) * 0.1) + float64(max))
		}
	}
	for o := 0; o < len(rules); o++ {
		if rules[o].Metric == "ResponseTime" {
			for _, h := range rules[o].Alarm {
				if h.Level == "Critical" && h.Threshold != "" {
					if v, ok := strconv.Atoi(h.Threshold); ok != nil {
						responseTime = 400
					} else {
						responseTime = v
					}
				}
			}
		} else if rules[o].Metric == "AvailablePercent" {
			for _, h := range rules[o].Alarm {
				if h.Level == "Critical" && h.Threshold != "" {
					if v, ok := strconv.Atoi(h.Threshold); ok != nil {
						availablePercent = 90
					} else {
						availablePercent = v
					}
				}
			}
		}
	}
	if len(series3) == 0 {
		data3 := make([]int, 0)
		series3 = []interface{}{
			H{"name": "响应时间", "data": data3},
		}
	}
	if max < 5 {
		max = 5
	} else if max >= responseTime {
		max = responseTime
	}

	return H{
		"options": newOptions, "rules": rules,
		"usable": availablePercent, "response": responseTime,
		"series": series1, "line1": series2, "line2": series3, "results": n.Results, "max": max,
	}
}

func isExist(array []string, k string) bool {
	for i := 0; i < len(array); i++ {
		if array[i] == k {
			return true
		}
	}
	return false
}
