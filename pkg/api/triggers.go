package api

import (
	"fmt"
	"time"
)

type ListTriggersResponse struct {
	Triggers []Trigger `json:"triggers"`
}

type Trigger struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Slug        *string           `json:"slug"`
	Kind        TriggerKind       `json:"kind"`
	KindConfig  TriggerKindConfig `json:"kindConfig"`
	DisabledAt  *time.Time        `json:"disabledAt"`
	ArchivedAt  *time.Time        `json:"archivedAt"`
}

type TriggerKind string

const (
	TriggerKindUnknown  TriggerKind = ""
	TriggerKindForm     TriggerKind = "form"
	TriggerKindSchedule TriggerKind = "schedule"
)

type TriggerKindConfig struct {
	Form     *TriggerKindConfigForm     `json:"form,omitempty"`
	Schedule *TriggerKindConfigSchedule `json:"schedule,omitempty"`
}

type TriggerKindConfigForm struct {
	Parameters Parameters `json:"parameters"`
}

type TriggerKindConfigSchedule struct {
	ParamValues map[string]interface{} `json:"paramValues"`
	CronExpr    CronExpr               `json:"cronExpr"`
}

type CronExpr struct {
	Minute     string `json:"minute,omitempty"`
	Hour       string `json:"hour,omitempty"`
	DayOfMonth string `json:"dayOfMonth,omitempty"`
	Month      string `json:"month,omitempty"`
	DayOfWeek  string `json:"dayOfWeek,omitempty"`
}

func (ce CronExpr) String() string {
	return fmt.Sprintf("%s %s %s %s %s", ce.Minute, ce.Hour, ce.DayOfMonth, ce.Month, ce.DayOfWeek)
}
