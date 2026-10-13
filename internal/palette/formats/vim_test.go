package formats

import (
	"strings"
	"testing"
)

func assertVimOutput(t *testing.T, theme Theme, wants []string) {
	t.Helper()
	out, err := FormatTheme(theme, "vim")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range wants {
		if !strings.Contains(out, want) {
			t.Errorf("vim output missing %q:\n%s", want, out)
		}
	}
}

func TestFormatVimHeader(t *testing.T) {
	dark := coherenceTheme()
	light := coherenceTheme()
	light.Background = "#f2ede8"

	tests := []struct {
		name  string
		theme Theme
		want  string
	}{
		{"colors name", dark, `let g:colors_name = "palette"`},
		{"dark background", dark, "set background=dark"},
		{"light background", light, "set background=light"},
		{"highlight clear", dark, "hi clear"},
		{"syntax guard", dark, `if exists("syntax_on")`},
		{"syntax reset", dark, "syntax reset"},
		{"normal", dark, "hi Normal guifg=#e8e1dd guibg=#19120c ctermfg=15 ctermbg=NONE"},
		{"normal light", light, "hi Normal guifg=#e8e1dd guibg=#f2ede8 ctermfg=15 ctermbg=NONE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertVimOutput(t, tt.theme, []string{tt.want})
		})
	}
}

func TestFormatVimHighlights(t *testing.T) {
	theme := coherenceTheme()

	tests := []struct {
		name string
		want string
	}{
		{"comment", "hi Comment guifg=#776b5f guibg=NONE ctermfg=8 ctermbg=NONE gui=italic cterm=italic"},
		{"string", "hi String guifg=#a2b275 guibg=NONE ctermfg=2 ctermbg=NONE"},
		{"line number", "hi LineNr guifg=#776b5f guibg=NONE ctermfg=8 ctermbg=NONE"},
		{"cursor line", "hi CursorLine guifg=NONE guibg=#261e15 ctermfg=NONE ctermbg=0"},
		{"visual", "hi Visual guifg=#19120c guibg=#ffc079 ctermfg=0 ctermbg=11"},
		{"search", "hi Search guifg=#19120c guibg=#d6a163 ctermfg=0 ctermbg=3"},
		{"pmenu", "hi Pmenu guifg=#e8e1dd guibg=#261e15 ctermfg=15 ctermbg=0"},
		{"pmenu thumb", "hi PmenuThumb guifg=NONE guibg=#776b5f ctermfg=NONE ctermbg=8"},
		{"diff add", "hi DiffAdd guifg=#19120c guibg=#a2b275 ctermfg=0 ctermbg=2"},
		{"match paren", "hi MatchParen guifg=#ffc079 guibg=NONE ctermfg=11 ctermbg=NONE"},
		{"match paren attrs", "gui=bold,underline cterm=bold,underline"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertVimOutput(t, theme, []string{tt.want})
		})
	}
}

func TestFormatVimLinks(t *testing.T) {
	theme := coherenceTheme()

	tests := []struct {
		name string
		want string
	}{
		{"character", "hi! link Character String"},
		{"float", "hi! link Float Number"},
		{"conditional", "hi! link Conditional Statement"},
		{"cursor column", "hi! link CursorColumn CursorLine"},
		{"wildmenu", "hi! link WildMenu PmenuSel"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertVimOutput(t, theme, []string{tt.want})
		})
	}
}

func TestFormatVimRequiredGroups(t *testing.T) {
	theme := coherenceTheme()
	out, err := FormatTheme(theme, "vim")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	groups := []string{
		"Normal", "Comment", "Constant", "String", "Character", "Number", "Boolean", "Float",
		"Identifier", "Function", "Statement", "Conditional", "Repeat", "Label", "Operator",
		"Keyword", "Exception", "PreProc", "Include", "Define", "Macro", "PreCondit", "Type",
		"StorageClass", "Structure", "Typedef", "Special", "SpecialChar", "Tag", "Delimiter",
		"SpecialComment", "Underlined", "Todo", "Error", "WarningMsg", "LineNr", "CursorLineNr",
		"CursorLine", "CursorColumn", "ColorColumn", "Visual", "Search", "IncSearch", "Pmenu",
		"PmenuSel", "PmenuSbar", "PmenuThumb", "StatusLine", "StatusLineNC", "VertSplit",
		"TabLine", "TabLineSel", "TabLineFill", "Folded", "FoldColumn", "SignColumn",
		"WildMenu", "Question", "MoreMsg", "ModeMsg", "NonText", "SpecialKey", "Whitespace",
		"MatchParen", "DiffAdd", "DiffChange", "DiffDelete", "DiffText",
	}
	for _, group := range groups {
		defined := strings.Contains(out, "hi "+group+" ") || strings.Contains(out, "hi! link "+group+" ")
		if !defined {
			t.Errorf("vim output does not define group %s:\n%s", group, out)
		}
	}
}
