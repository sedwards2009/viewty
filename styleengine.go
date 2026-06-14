package viewty

import (
	"encoding/json"
	"maps"
)

type styleRule struct {
	Selector    string            `json:"selector"`
	Styles      StyleMap          `json:"styles"`
	InheritsFrom string            `json:"from,omitempty"`
}

type StyleBuilder struct {
	defaultStyles StyleMap
	widgetRules   []styleRule
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
		Styles:     styles,
		InheritsFrom: inheritsFrom[0],
	}
	b.widgetRules = append(b.widgetRules, rule)
	return b
}

func (b *StyleBuilder) LoadJSON(config string) error {
	var rules []styleRule
	if err := json.Unmarshal([]byte(config), &rules); err != nil {
		return err
	}

	for _, rule := range rules {
		if rule.Selector == "*" {
			b.defaultStyles = rule.Styles
		} else if rule.Selector != "" {
			b.widgetRules = append(b.widgetRules, rule)
		}
	}

	return nil
}

func (b *StyleBuilder) Build() (StyleFunc, error) {
	var rules []styleRule
	for _, rule := range b.widgetRules {
		rules = append(rules, rule)
	}

	return func(base StyleMap, widgeType string, class []string) StyleMap {
		result := make(StyleMap)
		maps.Copy(result, base)
		maps.Copy(result, b.defaultStyles)

		for _, rule := range rules {
			if rule.InheritsFrom != "" {
				parent := result[rule.InheritsFrom]
				if parentMap, ok := parent.(StyleMap); ok {
					maps.Copy(result, parentMap)
				}
			}
			if rule.Selector == widgeType {
				maps.Copy(result, rule.Styles)
			}
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
