package formats

import (
	"fmt"
	"strconv"
	"strings"
)

const vimHighlightSpec = `Normal fg bg
Comment 8 - italic
Constant 3 -
String 2 -
Number 3 -
Boolean 3 -
Identifier fg -
Function 4 bold
Statement 5 bold
Operator 6 -
PreProc 6 -
Type 3 -
Special 4 -
Underlined - - underline
Todo 3 bg bold
Error 1 bg
WarningMsg 3 -
LineNr 8 -
CursorLineNr cur - bold
CursorLine - 0
ColorColumn - 8
Visual bg cur
Search bg 3
IncSearch bg cur
Pmenu fg 0
PmenuSel bg cur
PmenuSbar - 0
PmenuThumb - 8
StatusLine fg 0
StatusLineNC 8 bg
VertSplit 8 -
TabLine 8 bg underline
TabLineSel fg 0 bold
TabLineFill - bg
Folded 8 0
SignColumn 8 bg
Question 2 - bold
MoreMsg 2 -
ModeMsg fg - bold
NonText 8 -
SpecialKey 8 -
Whitespace 8 -
MatchParen cur - bold,underline
DiffAdd bg 2
DiffChange fg 0
DiffDelete bg 1
DiffText bg 3 bold
`

const vimLinkSpec = `Character String
Float Number
Conditional Statement
Repeat Statement
Label Statement
Keyword Statement
Exception Statement
Include PreProc
Define PreProc
Macro PreProc
PreCondit PreProc
StorageClass Type
Structure Type
Typedef Type
SpecialChar Special
Tag Special
Delimiter Special
SpecialComment Comment
CursorColumn CursorLine
FoldColumn SignColumn
WildMenu PmenuSel
`

func formatVim(t Theme) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("set background=%s\n", vimBackgroundMode(t)))
	b.WriteString("hi clear\n")
	b.WriteString("if exists(\"syntax_on\")\n  syntax reset\nendif\n")
	b.WriteString("let g:colors_name = \"palette\"\n\n")
	writeVimHighlights(&b, t)
	b.WriteString("\n")
	writeVimLinks(&b)
	return b.String()
}

func writeVimHighlights(b *strings.Builder, t Theme) {
	for _, row := range strings.Split(vimHighlightSpec, "\n") {
		fields := strings.Fields(row)
		if len(fields) < 3 {
			continue
		}
		fg, ctermFg := vimColor(t, fields[1], false)
		bg, ctermBg := vimColor(t, fields[2], true)
		line := fmt.Sprintf("hi %s guifg=%s guibg=%s ctermfg=%s ctermbg=%s",
			fields[0], fg, bg, ctermFg, ctermBg)
		if len(fields) > 3 {
			line += fmt.Sprintf(" gui=%s cterm=%s", fields[3], fields[3])
		}
		b.WriteString(line + "\n")
	}
}

func writeVimLinks(b *strings.Builder) {
	for _, row := range strings.Split(vimLinkSpec, "\n") {
		fields := strings.Fields(row)
		if len(fields) < 2 {
			continue
		}
		b.WriteString(fmt.Sprintf("hi! link %s %s\n", fields[0], fields[1]))
	}
}

func vimColor(t Theme, token string, isBackground bool) (string, string) {
	switch token {
	case "-", "":
		return "NONE", "NONE"
	case "fg":
		return t.Foreground, vimCterm(t, t.Foreground)
	case "bg":
		if isBackground {
			return t.Background, "NONE"
		}
		return t.Background, vimCterm(t, t.Background)
	case "cur":
		return t.Cursor, vimCterm(t, t.Cursor)
	default:
		index, err := strconv.Atoi(token)
		if err != nil || index < 0 || index > 15 {
			return "NONE", "NONE"
		}
		return t.Palette[index], strconv.Itoa(index)
	}
}

func vimCterm(t Theme, hex string) string {
	r, g, b := vimRGB(hex)
	best, bestDist := 0, 1<<62
	for i, c := range t.Palette {
		pr, pg, pb := vimRGB(c)
		dist := (r-pr)*(r-pr) + (g-pg)*(g-pg) + (b-pb)*(b-pb)
		if dist < bestDist {
			best, bestDist = i, dist
		}
	}
	return strconv.Itoa(best)
}

func vimRGB(hex string) (int, int, int) {
	digits := strings.TrimPrefix(hex, "#")
	if len(digits) != 6 {
		return 0, 0, 0
	}
	r, _ := strconv.ParseUint(digits[0:2], 16, 32)
	g, _ := strconv.ParseUint(digits[2:4], 16, 32)
	b, _ := strconv.ParseUint(digits[4:6], 16, 32)
	return int(r), int(g), int(b)
}

func vimBackgroundMode(t Theme) string {
	r, g, b := vimRGB(t.Background)
	if (299*r+587*g+114*b)/1000 >= 128 {
		return "light"
	}
	return "dark"
}
