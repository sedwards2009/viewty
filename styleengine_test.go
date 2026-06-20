package viewty

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestLoadStyleRulesJSON(t *testing.T) {
	config := `{
		"*": {
			"foregroundColor": "#ffffff",
			"backgroundColor": "#000000"
		},
		"Button": {
			"foregroundColor": "#ff0000"
		}
	}`

	styleFunc, err := LoadStyleRules(config)
	assert.NoError(t, err)
	assert.NotNil(t, styleFunc)

	result := styleFunc(nil, "Button", nil)
	assert.Equal(t, tcell.GetColor("#000000"), result["backgroundColor"])
	assert.Equal(t, tcell.GetColor("#ff0000"), result["foregroundColor"])
}

func TestLoadStyleRulesJSONNoDefault(t *testing.T) {
	config := `{
		"Button": {
			"foregroundColor": "#ff0000"
		}
	}`

	styleFunc, err := LoadStyleRules(config)
	assert.NoError(t, err)
	assert.NotNil(t, styleFunc)

	result := styleFunc(nil, "Button", nil)
	assert.Equal(t, tcell.GetColor("#ff0000"), result["foregroundColor"])
}

func TestLoadStyleRulesJSONInvalid(t *testing.T) {
	_, err := LoadStyleRules(`invalid json`)
	assert.Error(t, err)
}

func TestLoadStyleRulesJSONInheritance(t *testing.T) {
	config := `{
		"Base": {
			"foregroundColor": "#ffffff"
		},
		"Button": {
			"from": "Base",
			"backgroundColor": "#000000"
		}
	}`

	styleFunc, err := LoadStyleRules(config)
	assert.NoError(t, err)
	assert.NotNil(t, styleFunc)

	result := styleFunc(nil, "Button", nil)
	assert.Equal(t, tcell.GetColor("#000000"), result["backgroundColor"])
	assert.Equal(t, tcell.GetColor("#ffffff"), result["foregroundColor"])
}

func TestNewStyleBuilder(t *testing.T) {
	builder := NewStyleBuilder()
	assert.NotNil(t, builder)
}

func TestStyleBuilderSetDefaultStyles(t *testing.T) {
	defaultMap := StyleMap{"foregroundColor": tcell.GetColor("#ffffff")}
	builder := NewStyleBuilder().SetDefaultStyles(defaultMap)
	assert.NotNil(t, builder)
}

func TestStyleBuilderAddWidgetRule(t *testing.T) {
	builder := NewStyleBuilder().AddWidgetRule("Button", StyleMap{"foregroundColor": tcell.GetColor("#ff0000")})
	assert.NotNil(t, builder)
}

func TestStyleBuilderLoadJSON(t *testing.T) {
	config := `{
		"*": {
			"foregroundColor": "#ffffff"
		},
		"Button": {
			"foregroundColor": "#ff0000"
		}
	}`

	err := NewStyleBuilder().LoadJSON(config)
	assert.NoError(t, err)
}

func TestStyleBuilderBuild(t *testing.T) {
	config := `{
		"*": {
			"foregroundColor": "#ffffff",
			"backgroundColor": "#000000"
		},
		"Button": {
			"foregroundColor": "#ff0000"
		}
	}`

	builder := NewStyleBuilder()
	err := builder.LoadJSON(config)
	assert.NoError(t, err)

	styleFunc, err := builder.Build()
	assert.NoError(t, err)
	assert.NotNil(t, styleFunc)

	result := styleFunc(nil, "Button", nil)
	assert.Equal(t, tcell.GetColor("#000000"), result["backgroundColor"])
	assert.Equal(t, tcell.GetColor("#ff0000"), result["foregroundColor"])
}

func TestStyleBuilderBuildWithBaseStyles(t *testing.T) {
	config := `{
		"Button": {
			"foregroundColor": "#ff0000"
		}
	}`

	builder := NewStyleBuilder()
	err := builder.LoadJSON(config)
	assert.NoError(t, err)

	baseStyles := StyleMap{"backgroundColor": tcell.GetColor("#00ff00")}
	styleFunc, err := builder.Build()
	assert.NoError(t, err)

	result := styleFunc(baseStyles, "Button", nil)
	assert.Equal(t, tcell.GetColor("#ff0000"), result["foregroundColor"])
	assert.Equal(t, tcell.GetColor("#00ff00"), result["backgroundColor"])
}

func TestStyleBuilderLoadJSONInvalid(t *testing.T) {
	err := NewStyleBuilder().LoadJSON(`invalid json`)
	assert.Error(t, err)
}

func TestParseRuleStylesColors(t *testing.T) {
	config := `{
		"Button": {
			"foregroundColor": "#ff0000",
			"backgroundColor": "#00ff00",
			"size": "large"
		}
	}`

	err := NewStyleBuilder().LoadJSON(config)
	assert.NoError(t, err)

	builder := NewStyleBuilder()
	err = builder.LoadJSON(config)
	assert.NoError(t, err)

	styleFunc, err := builder.Build()
	assert.NoError(t, err)

	result := styleFunc(nil, "Button", nil)
	assert.Equal(t, tcell.GetColor("#ff0000"), result["foregroundColor"])
	assert.Equal(t, tcell.GetColor("#00ff00"), result["backgroundColor"])
	assert.Equal(t, "large", result["size"])
}

func TestParseRuleStylesWithInheritance(t *testing.T) {
	config := `{
		"Base": {
			"foregroundColor": "#ffffff",
			"size": "small"
		},
		"Button": {
			"from": "Base",
			"backgroundColor": "#0000ff",
			"foregroundColor": "#ff0000"
		}
	}`

	err := NewStyleBuilder().LoadJSON(config)
	assert.NoError(t, err)

	styleFunc, err := LoadStyleRules(config)
	assert.NoError(t, err)

	result := styleFunc(nil, "Button", nil)
	assert.Equal(t, tcell.GetColor("#0000ff"), result["backgroundColor"])
	assert.Equal(t, tcell.GetColor("#ff0000"), result["foregroundColor"])
	assert.Equal(t, "small", result["size"])
}
