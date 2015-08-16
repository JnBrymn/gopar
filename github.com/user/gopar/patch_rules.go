package gopar


import (
	"fmt"
)


type PatchError string
func (p PatchError) Error() string {
	return string(p)
}


func collectRules(
					rule Parser,
					name2Rule map[string]Parser,
					placeholderRules []*placeholderRule,
				) []*placeholderRule {
	switch rule := rule.(type) {
		default: 
			name2Rule[rule.GetName()] = rule
			for _, rule := range rule.GetSubRules() {
				placeholderRules = collectRules(rule, name2Rule, placeholderRules)
			}
		case *placeholderRule:
			placeholderRules = append(placeholderRules,rule)
	}
	return placeholderRules
}

func Patch(rules... Parser) error {
	name2Rule := map[string]Parser{}
	placeholderRules := []*placeholderRule{}
	for _,rule := range rules {
		placeholderRules = collectRules(rule, name2Rule, placeholderRules)
	}
	for _,rule := range placeholderRules {
		if patchRule, ok := name2Rule[rule.GetName()]; ok {
			rule.patchRule = patchRule	
		} else {
			return PatchError(fmt.Sprintf("couldn't find patch rule '%s'",rule.GetName()))
		}
	}
	return nil
}