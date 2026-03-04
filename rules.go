package loglint

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

var sensitiveKeywords = []string{
	"password", "passwd",
	"secret",
	"token",
	"api_key", "apikey",
	"credential",
	"private_key", "privatekey",
	"access_key", "accesskey",
	"session_id", "sessionid",
}

func checkAllRules(pass *analysis.Pass, msgExpr ast.Expr) {
	lits := extractStringLiterals(msgExpr)
	if len(lits) == 0 {
		return
	}

	msg := fullMessage(lits)
	first := lits[0]

	// rule 1 — message must start with a lowercase letter
	checkLowercaseStart(pass, first, msg)

	// rule 2 — only Latin (English) text allowed
	checkEnglishOnly(pass, msgExpr, msg)

	// rule 3 — no emoji / special decorative chars
	checkSpecialChars(pass, msgExpr, msg)

	// rule 4 — no sensitive data leaking via concatenation
	checkSensitiveData(pass, msgExpr, lits)
}

// ------------------------------------------------------------------
// Rule 1: first letter must be lowercase
// ------------------------------------------------------------------

func checkLowercaseStart(pass *analysis.Pass, first stringLiteral, msg string) {
	if msg == "" {
		return
	}
	r, _ := utf8.DecodeRuneInString(msg)
	if !unicode.IsUpper(r) {
		return
	}

	diag := analysis.Diagnostic{
		Pos:     first.node.Pos(),
		Message: "log message should start with a lowercase letter",
	}

	// suggested fix: lowercase the first character of the literal
	val := first.value
	if len(val) > 0 {
		fr, sz := utf8.DecodeRuneInString(val)
		fixed := string(unicode.ToLower(fr)) + val[sz:]
		diag.SuggestedFixes = []analysis.SuggestedFix{{
			Message: "lowercase the first letter",
			TextEdits: []analysis.TextEdit{{
				Pos:     first.node.Pos(),
				End:     first.node.End(),
				NewText: []byte(strconv.Quote(fixed)),
			}},
		}}
	}

	pass.Report(diag)
}

// ------------------------------------------------------------------
// Rule 2: message must be in English (Latin script only)
// ------------------------------------------------------------------

func checkEnglishOnly(pass *analysis.Pass, expr ast.Expr, msg string) {
	for _, r := range msg {
		if unicode.IsLetter(r) && !unicode.Is(unicode.Latin, r) {
			pass.Reportf(expr.Pos(), "log message must be in English, found non-Latin characters")
			return
		}
	}
}

// ------------------------------------------------------------------
// Rule 3: no special characters or emoji
// ------------------------------------------------------------------

func checkSpecialChars(pass *analysis.Pass, expr ast.Expr, msg string) {
	hasEmoji := false
	for _, r := range msg {
		if isEmoji(r) {
			hasEmoji = true
			break
		}
	}
	hasSpecial := strings.Contains(msg, "!") ||
		strings.Contains(msg, "...") ||
		strings.ContainsRune(msg, '\u2026') // …

	if !hasEmoji && !hasSpecial {
		return
	}

	diag := analysis.Diagnostic{
		Pos:     expr.Pos(),
		Message: "log message should not contain special characters or emoji",
	}

	// suggest fix when the whole expression is a single string literal
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		cleaned := cleanMessage(msg)
		if cleaned != msg {
			diag.SuggestedFixes = []analysis.SuggestedFix{{
				Message: "remove special characters and emoji",
				TextEdits: []analysis.TextEdit{{
					Pos:     lit.Pos(),
					End:     lit.End(),
					NewText: []byte(strconv.Quote(cleaned)),
				}},
			}}
		}
	}

	pass.Report(diag)
}

func cleanMessage(msg string) string {
	var b strings.Builder
	for _, r := range msg {
		if isEmoji(r) || r == '!' || r == '\u2026' {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	for strings.Contains(out, "...") {
		out = strings.ReplaceAll(out, "...", "")
	}
	return strings.TrimSpace(out)
}

func isEmoji(r rune) bool {
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols & Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport & Map
		(r >= 0x1F1E0 && r <= 0x1F1FF) || // Flags
		(r >= 0x2600 && r <= 0x26FF) || // Misc Symbols
		(r >= 0x2700 && r <= 0x27BF) || // Dingbats
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols
		(r >= 0x1FA00 && r <= 0x1FAFF) // Symbols Extended-A
}

// ------------------------------------------------------------------
// Rule 4: no sensitive data in concatenated log messages
// ------------------------------------------------------------------

func checkSensitiveData(pass *analysis.Pass, expr ast.Expr, lits []stringLiteral) {
	// only flag when the message mixes literals with runtime values
	if !hasNonLiteralParts(expr) {
		return
	}

	combined := strings.ToLower(fullMessage(lits))
	for _, kw := range sensitiveKeywords {
		if strings.Contains(combined, kw) {
			pass.Reportf(expr.Pos(), "log message may contain sensitive data (keyword: %q)", kw)
			return
		}
	}
}
