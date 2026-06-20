package viewty

import (
	"encoding/json"
	"fmt"
	"maps"
)

type styleRule struct {
	Selector       string   `json:"selector"`
	InputStyles    StyleMap `json:"styles"`
	InheritsFrom   string   `json:"from,omitempty"`
	ExpandedStyles StyleMap
}

type StyleBuilder struct {
	defaultStyles StyleMap
	rules   []*styleRule
}

func NewStyleBuilder() *StyleBuilder {
	return &StyleBuilder{}
}

func (b *StyleBuilder) SetDefaultStyles(styles StyleMap) *StyleBuilder {
	b.defaultStyles = styles
	return b
}

func (b *StyleBuilder) AddWidgetRule(selector string, styles StyleMap, inheritsFrom ...string) *StyleBuilder {
	rule := styleRule{
		Selector:   selector,
		InputStyles:     styles,
	}
	if len(inheritsFrom) != 0 {
		rule.InheritsFrom = inheritsFrom[0]
	}
	b.rules = append(b.rules, &rule)
	return b
}

func (b *StyleBuilder) LoadJSON(config string) error {
	var ruleMap map[string]any
	if err := json.Unmarshal([]byte(config), &ruleMap); err != nil {
		return err
	}

	for selector, value := range ruleMap {
		if valueMap, ok := value.(map[string]any); ok {
			var inheritsFrom string
			if inherits, ok := valueMap["from"].(string); ok {
				inheritsFrom = inherits
			}

			rule := &styleRule{
				Selector:   selector,
				InputStyles:     make(StyleMap),
				InheritsFrom: inheritsFrom,
			}
			rule.InputStyles = valueMap

			if selector == "*" {
				b.defaultStyles = rule.InputStyles
			} else if selector != "" {
				b.rules = append(b.rules, rule)
			}
		}
	}

	return nil
}

func (b *StyleBuilder) computeAllStyles() {
	for _, rule := range b.rules {
		rule.ExpandedStyles = make(StyleMap)
	}
	for _, rule := range b.rules {
		err := b.computeStyleRule(rule, len(b.rules)+1)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}

func (b *StyleBuilder) computeStyleRule(rule *styleRule, recursionLimit int) error {
	if recursionLimit == 0 {
		return fmt.Errorf("Recursion loop detected")
	}
	if len(rule.ExpandedStyles) >= len(rule.InputStyles) {
		return nil	// Already updated
	}

	if rule.InheritsFrom == "" {
		maps.Copy(rule.ExpandedStyles, b.defaultStyles)
	} else {
		parentRule := b.getRuleBySelector(rule.InheritsFrom)
		if parentRule == nil {
			return fmt.Errorf("Couldn't find style rule with name '%s'.", rule.InheritsFrom)
		}
		err := b.computeStyleRule(parentRule, recursionLimit-1)
		if err != nil {
			return err
		}
		maps.Copy(rule.ExpandedStyles, parentRule.ExpandedStyles)
	}
	maps.Copy(rule.ExpandedStyles, rule.InputStyles)

	return nil
}

func (b *StyleBuilder) getRuleBySelector(selector string) *styleRule {
  for _, rule := range b.rules {
    if rule.Selector == selector {
      return rule
    }
  }
  return nil
}

func (b *StyleBuilder) Build() (StyleFunc, error) {
	b.computeAllStyles()

	return func(base StyleMap, widgetType string, class []string) StyleMap {
		result := make(StyleMap)
		maps.Copy(result, base)

		rule := b.getRuleBySelector(widgetType)
		if rule != nil {
			maps.Copy(result, rule.ExpandedStyles)
		}
		return result
	}, nil
}

func LoadStyleRules(config string) (StyleFunc, error) {
	builder := NewStyleBuilder()
	if err := builder.LoadJSON(config); err != nil {
		return nil, err
	}
	return builder.Build()
}
