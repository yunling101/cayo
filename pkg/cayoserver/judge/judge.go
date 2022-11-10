package judge

import (
	"fmt"
	"log"
	"strings"

	"github.com/yunling101/cayo/pkg/cayoserver/sender"
	"github.com/yunling101/cayo/pkg/model/alarm"
	"github.com/yunling101/cayo/pkg/model/q"
	"github.com/yunling101/cayo/pkg/model/settings"
	"github.com/yunling101/cayo/pkg/types"
	"github.com/yunling101/toolkits/date"
	"github.com/yunling101/toolkits/text"
)

type Cond struct {
	NodeID  int             `json:"node_id"`
	Rule    alarm.AlarmRule `json:"rule"`
	Alarm   types.Alarm     `json:"alarm"`
	alert   alarm.AlarmTask
	channel chan []string
}

// Judge 判读模块，报警数据入库与分派
func (c *Cond) Judge(current float64) {
	c.alert = alarm.AlarmTask{
		Team:      c.Rule.Team,
		NodeID:    c.NodeID,
		RuleID:    c.Rule.ID,
		ID:        text.GenerateRandomToInt(7),
		TaskID:    c.Rule.Task.ID,
		TaskName:  c.Rule.Task.Name,
		Threshold: c.Alarm.Threshold,
		CurValue:  current,
		Level:     c.Alarm.Level,
		Unit:      alarm.GetMetricUnit(c.Rule.Metric),
		Title: fmt.Sprintf("%s 连续%v次, %s %v %s%s.",
			c.Rule.Task.Name, c.Alarm.Continuity, c.Rule.Name, c.Rule.Condition, c.Alarm.Threshold,
			alarm.GetMetricUnit(c.Rule.Metric),
		),
		Status:      "trigger",
		Archive:     false,
		ArchiveTime: date.New().Now,
		Counter:     1, // 计数器默认为1。每触发一次加1
		Metric:      alarm.GetMetricName(c.Rule.Metric),
		CreateTime:  date.New().Now,
	}
	c.channel = make(chan []string, 1)
	for idx, id := range c.Rule.Notify {
		var contact alarm.AlarmContact
		err := q.Table(contact.TableName()).QueryOne(q.M{"id": id}, &contact)
		if err == nil {
			// 使用local方式，所有联系人通知渠道是一样的
			// 每个联系人的通知渠道可不一样
			c.alert.Users = append(c.alert.Users, types.Users{
				ID: contact.ID, Username: contact.Username, Channel: contact.Channel,
			})
			// 异步分派给每个联系人
			go c.judgeChannel(idx, contact)
		}
	}
	// 修改local方式，把分派渠道赋值给用户
	if n, ok := <-c.channel; ok {
		if len(n) != 0 {
			for index, u := range c.alert.Users {
				if u.Channel == "local" {
					c.alert.Users[index].Channel = strings.Join(n, ",")
				}
			}
		}
		close(c.channel)
	}
	go func() {
		if err := c.alert.Save(); err != nil {
			log.Printf("task.id %s rule.id %v alarm save error: %s", c.Rule.Task.ID, c.Rule.ID, err.Error())
		}
	}()
}

// judgeChannel
func (c *Cond) judgeChannel(index int, contact alarm.AlarmContact) {
	send := sender.Channel{Alert: c.alert, Contact: contact, Type: contact.Channel}
	notifyChannel := make([]string, 1)
	switch contact.Channel {
	case "local":
		cfg := settings.Settings{}.GetUserSettings()
		if cfg.UserLocal {
			if err := c.workerTime(cfg); err != nil {
				log.Printf("task.id %s work time error: %s", c.alert.TaskID, err.Error())
				return
			}
			// 遍历用户配置策略与规则策略一致时，分发到开启的通道
			for _, strategy := range cfg.NotifyStrategy {
				if strategy.Name == c.Alarm.Level {
					// 遍历开启的通道
					for _, notify := range strategy.Notify {
						if notify.Enable {
							send.Type = notify.Alias
							notifyChannel = append(notifyChannel, notify.Alias)
							service := send.NewChannel()
							if r := service.Delivery(); r.Err != nil {
								log.Printf("delivery(%s) task.id %s error: %s",
									notify.Alias, c.alert.TaskID, r.Err.Error())
							}
						}
					}
				}
			}
		}
	default:
		service := send.NewChannel()
		if r := service.Delivery(); r.Err != nil {
			log.Printf("delivery(%s) task.id %s error: %s", contact.Channel, c.alert.TaskID, r.Err.Error())
		}
	}
	if index == 0 {
		c.channel <- notifyChannel
	}
}

// workerTime
func (c *Cond) workerTime(cfg settings.Settings) error {
	if cfg.SendTime {
		n := date.New()
		start, err := n.StringToDate(cfg.StartTime, "2006-01-02 15:04:05")
		if err != nil {
			return err
		}
		end, err := n.StringToDate(cfg.EndTime, "2006-01-02 15:04:05")
		if err != nil {
			return err
		}
		if n.Now.Unix() >= start.Unix() && n.Now.Unix() <= end.Unix() {
			return nil
		}
		return fmt.Errorf("%s", "out of working hours")
	}
	return nil
}
