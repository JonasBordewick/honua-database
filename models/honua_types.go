package models

import (
	"time"
)

type Identity struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SensorType int32

const (
	NONE SensorType = iota
	ACLOADS
	TOTALPV
	GRID
	SOC
	BATTERYVALUE
	BATTERYSTATE
)

type Entity struct {
	Id              int        `json:"id"`
	IdentityId      string     `json:"identity"`
	EntityId        string     `json:"entity_id"`
	Name            string     `json:"name"`
	IsDevice        bool       `json:"is_device"`
	AllowRules      bool       `json:"allow_rules"`
	HasAttribute    bool       `json:"has_attribute"`
	Attribute       string     `json:"attribute"`
	IsVictronSensor bool       `json:"is_victron_sensor"`
	SensorType      SensorType `json:"sensor_type"`
	HasNumericState bool       `json:"has_numeric_state"`
	RulesEnabled    bool       `json:"rules_enabled"`
}

func (e *Entity) Equals(o *Entity) bool {
	return (e.IdentityId == o.IdentityId) && (e.EntityId == o.EntityId) && (e.Name == o.Name) && (e.IsDevice == o.IsDevice) && (e.AllowRules == o.AllowRules) && (e.HasAttribute == o.HasAttribute) && (e.Attribute == o.Attribute) && (e.IsVictronSensor == o.IsVictronSensor) && (e.HasNumericState == o.HasNumericState)
}

type State struct {
	Id         int
	EntityId   int
	State      string
	RecordTime *time.Time
}

type Rule struct {
	Id                   int
	Enabled              bool
	EventBasedEvaluation bool
	PeriodicTrigger      PeriodicTriggerType
	Name                 string
	Target               *Entity
	Condition            *Condition
	ThenActions          []*Action
	ElseActions          []*Action
}

type PeriodicTriggerType int

const (
	OneMin PeriodicTriggerType = iota
	TwoMin
	FiveMin
	TenMin
	FifteenMin
	TwentyMin
	TwentyFiveMin
	FortyFiveMin
	OneH
	TwoH
	SixH
)

type Condition struct {
	Id              int
	Type            ConditionType
	Sensor          *Entity
	ComparisonState string
	After           string
	Before          string
	Above           *ConditionValue
	Below           *ConditionValue
	SubConditions   []*Condition
}

type ConditionType int

const (
	AND ConditionType = iota
	OR
	NAND
	NOR
	NUMERICSTATE
	STATE
	TIME
)

type ConditionValue struct {
	Valid bool
	Value int
}

type ActionType int

const (
	SERVICE ActionType = iota
	DELAY
)

type Action struct {
	Id      int
	Type    ActionType
	Service string
	Delay   *Delay
}

type Service struct {
	Domain      string
	Name        string
	Description string
}

type Delay struct {
	Id      int
	Hours   int32
	Minutes int32
	Seconds int32
}
