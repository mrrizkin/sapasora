// Package sql provides utilities for sapasorag with SQL
package sql

import (
	"fmt"
	"strings"
)

// LogicBuilder helps build SQL-like conditions with proper parameter handling
type LogicBuilder struct {
	conditions []condition
	args       []any
	isNested   bool
}

type condition struct {
	logic string
	op    string // "AND" or "OR"
}

// NewLogicBuilder creates a new LogicBuilder instance
func NewLogicBuilder() *LogicBuilder {
	return &LogicBuilder{
		conditions: make([]condition, 0),
		args:       make([]any, 0),
		isNested:   false,
	}
}

// And adds a condition with AND logic
func (lb *LogicBuilder) And(logic any, args ...any) *LogicBuilder {
	return lb.addCondition("AND", logic, args...)
}

// Or adds a condition with OR logic
func (lb *LogicBuilder) Or(logic any, args ...any) *LogicBuilder {
	return lb.addCondition("OR", logic, args...)
}

// addCondition is a helper function to add conditions with the specified operator
func (lb *LogicBuilder) addCondition(
	op string,
	logic any,
	args ...any,
) *LogicBuilder {
	switch v := logic.(type) {
	case string:
		// Regular string condition
		lb.conditions = append(lb.conditions, condition{
			logic: v,
			op:    op,
		})
		lb.args = append(lb.args, args...)
	case func(*LogicBuilder):
		// Nested condition using a builder function
		nested := NewLogicBuilder()
		nested.isNested = true
		v(nested)

		nestedLogic, nestedArgs := nested.GetLogic()
		if nestedLogic != "" {
			lb.conditions = append(lb.conditions, condition{
				logic: fmt.Sprintf("(%s)", nestedLogic),
				op:    op,
			})
			lb.args = append(lb.args, nestedArgs...)
		}
	}

	return lb
}

// GetLogic returns the constructed logic string and arguments
func (lb *LogicBuilder) GetLogic() (string, []any) {
	if len(lb.conditions) == 0 {
		return "", nil
	}

	var builder strings.Builder

	for i, cond := range lb.conditions {
		if i > 0 {
			builder.WriteString(fmt.Sprintf(" %s ", cond.op))
		}
		builder.WriteString(cond.logic)
	}

	return builder.String(), lb.args
}

// Reset clears all conditions and arguments
func (lb *LogicBuilder) Reset() {
	lb.conditions = make([]condition, 0)
	lb.args = make([]any, 0)
}

// Count returns the number of conditions
func (lb *LogicBuilder) Count() int {
	return len(lb.conditions)
}

// HasConditions returns true if there are any conditions
func (lb *LogicBuilder) HasConditions() bool {
	return len(lb.conditions) > 0
}
