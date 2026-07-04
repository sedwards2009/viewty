package viewty

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/bits-and-blooms/bitset"
	"github.com/gdamore/tcell/v2"
)

/* CSS-like styling rules engine
 * =============================
 *
 * Matching
 * --------
 * Like CSS, rules are matched and applied in an order governed by specificity.
 *
 * Each widget requests its style rules based on its widget type (string) and
 * the list of classes assigned to the widget ([]string). These values will be
 * referred to as the "widget's properties".
 *
 * The ranking is in two phases:
 *
 * * Empty widget type selectors
 * * Specific widget type selectors
 *
 * In each phase, rules are matched based on the widget's properties.
 *
 * In the empty widget type selector phase, all the rules which don't have a
 * widget selector are considered. They match every widget.
 *
 * A widget's classes match a selector if the set of classes in the selector
 * are a subset of those on the widget. The rank of the match is determined
 * by the number of classes involved. A matching selector with many classes
 * is given high priority compared to a selector with fewer classes.
 *
 * In each phase every rule which matches the widget is collected, then matched
 * based on the class selector, then scored based on the quality of the class
 * match. The matching rules are sorted by the quality of the class match.
 * Finally, the rules are applied.
 *
 *
 */

type styleRule struct {
	Selector       string   `json:"selector"`
	InputStyles    StyleMap `json:"styles"`
	InheritsFrom   string   `json:"from,omitempty"`
	ExpandedStyles StyleMap
	widgetSelector string
	classes        []string
	classBitset    bitset.BitSet

	// Tmp value used when scoring rules against some widget properties.
	tmpScore       uint
}

type StyleBuilder struct {
	rules           []*styleRule
	classBitMapping map[string]uint
}

func NewStyleBuilder() *StyleBuilder {
	return &StyleBuilder{
		classBitMapping: make(map[string]uint),
	}
}

func (b *StyleBuilder) AddWidgetRule(selector string, styles StyleMap, inheritsFrom ...string) *StyleBuilder {
	rule := &styleRule{
		Selector:   selector,
		InputStyles:     styles,
	}
	if len(inheritsFrom) != 0 {
		rule.InheritsFrom = inheritsFrom[0]
	}
	b.rules = append(b.rules, rule)
	b.parseRuleStylesInPlace(rule)
	return b
}

func (b *StyleBuilder) parseRuleStylesInPlace(rule *styleRule) {
	for key, value := range rule.InputStyles {
		if strValue, ok := value.(string); ok {
			if strings.HasSuffix(key, "Color") {
				color := tcell.GetColor(strValue)
				rule.InputStyles[key] = color
			}
		}
	}
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
			b.parseRuleStylesInPlace(rule)
			b.rules = append(b.rules, rule)
		}
	}

	return nil
}

func (b *StyleBuilder) computeAllStyles() {
	for _, rule := range b.rules {
		b.parseRuleSelector(rule)
		rule.ExpandedStyles = make(StyleMap)
	}
	for _, rule := range b.rules {
		err := b.expandStyleRule(rule, len(b.rules)+1)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	b.indexClasses()
}

func (b *StyleBuilder) indexClasses() {
	// Compute the bitset indexes

	// Collect all of the different classes used.
	for _, rule := range b.rules {
		for _, class := range rule.classes {
			b.classBitMapping[class] = 0
		}
	}

	// Assign bit indexes to each class.
	var i uint = 0
	for key := range maps.Keys(b.classBitMapping) {
		b.classBitMapping[key] = i
		i++
	}

	// Compute the bitset index for each rule
    for _, rule := range b.rules {
       rule.classBitset.ClearAll()
       for _, class := range rule.classes {
         rule.classBitset.Set(b.classBitMapping[class])
       }
    }
}

func (b *StyleBuilder) parseRuleSelector(rule *styleRule) {
	selector := rule.Selector
	widgetSelector := ""
	classes := make([]string, 0)
	parts := strings.Split(selector, ".")
	if parts[0] != "" {
		widgetSelector = parts[0]
		parts = parts[1:]
	}
	for _, part := range parts {
		if part != "" {
			classes = append(classes, part)
		}
	}
	rule.widgetSelector = widgetSelector
	rule.classes = classes
}

func (b *StyleBuilder) expandStyleRule(rule *styleRule, recursionLimit int) error {
	if recursionLimit == 0 {
		return fmt.Errorf("Recursion loop detected")
	}
	if len(rule.ExpandedStyles) >= len(rule.InputStyles) {
		return nil	// Already updated
	}

	if rule.InheritsFrom != "" {
		parentRule := b.getRuleBySelector(rule.InheritsFrom)
		if parentRule == nil {
			return fmt.Errorf("Couldn't find style rule with name '%s'.", rule.InheritsFrom)
		}
		err := b.expandStyleRule(parentRule, recursionLimit-1)
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

	return func(base StyleMap, widgetType string, classes []string) StyleMap {
		result := make(StyleMap)
		maps.Copy(result, base)

		matchingEmptyRules := b.getAllRulesBySelector("", classes)
		for _, rule := range matchingEmptyRules {
			maps.Copy(result, rule.ExpandedStyles)
		}

		matchingRules := b.getAllRulesBySelector(widgetType, classes)
		for _, rule := range matchingRules {
			maps.Copy(result, rule.ExpandedStyles)
		}

		return result
	}, nil
}

func (b *StyleBuilder) getAllRulesBySelector(widgetType string, classes []string) []*styleRule {
	var classSelectorBitset bitset.BitSet
 	for _, class := range classes {
        classSelectorBitset.Set(b.classBitMapping[class])
    }

	result := make([]*styleRule, 0)
	for _, rule := range b.rules {
		if rule.widgetSelector == widgetType {
			intersectionCount := classSelectorBitset.IntersectionCardinality(&rule.classBitset)
			if intersectionCount == rule.classBitset.Count() {
				result = append(result, rule)
				rule.tmpScore = intersectionCount
			}
		}
	}

	slices.SortFunc(result, func(a *styleRule, b *styleRule) int {
		if a.tmpScore != b.tmpScore {
			return int(a.tmpScore) - int(b.tmpScore)
		}
		return 0
	})

	return result
}

func LoadStyleRules(config string) (StyleFunc, error) {
	builder := NewStyleBuilder()
	if err := builder.LoadJSON(config); err != nil {
		return nil, err
	}
	return builder.Build()
}
