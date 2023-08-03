package models

import (
	"time"
)

type Identity struct {
	Id   string
	Name string
}

type Entity struct {
	Id              int
	IdentityId      string
	EntityId        string
	Name            string
	IsDevice        bool
	AllowRules      bool
	HasAttribute    bool
	Attribute       string
	IsVictronSensor bool
	HasNumericState bool
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
	PeriodicTrigger      int
	Name                 string
	Target               Entity
	Condition            *Condition
	ThenActions          *[]Action
	ElseActions          *[]Action
}

type Condition struct {
	Id              int
	Type            int
	Sensor          Entity
	ComparisonState string
	After           string
	Before          string
	Above           *ConditionValue
	Below           *ConditionValue
	SubConditions   *[]Condition
}

type ConditionValue struct {
	Valid bool
	Value int
}

type Action struct {
	Id      int
	Service string
}

type Service struct {
	Domain      string
	Name        string
	Description string
}
