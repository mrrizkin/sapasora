package logger

import (
	"regexp"
	"strings"
)

const redactedValue = "[REDACTED]"

var (
	sensitiveKeyPattern     = regexp.MustCompile(`(?i)(password|passwd|secret|token|api[_-]?key|authorization|credential|private[ _-]?key|dsn|webhook|url|message[_-]?(body|text)?|body|content|phone|email|cookie|session)`)
	credentialURLPattern    = regexp.MustCompile(`(?i)([a-z][a-z0-9+.-]*://)[^/\s:@]+:[^/\s@]+@`)
	secretQueryPattern      = regexp.MustCompile(`(?i)([?&](?:password|passwd|secret|token|api[_-]?key|authorization|credential|sig(?:nature)?)[^=]*=)[^&#\s]+`)
	structuredSecretPattern = regexp.MustCompile(`(?i)\b(password|passwd|secret|token|api[_-]?key|authorization|credential|private[ _-]?key|dsn)\b(\s*[:=]\s*)("[^"]*"|'[^']*'|[^\s,;]+)`)
	sqlStringPattern        = regexp.MustCompile(`'(?:''|[^'])*'`)
	sqlNumericValuePattern  = regexp.MustCompile(`(?i)(\b(?:limit|offset)\s+|[=<>!]+\s*|,\s*|\(\s*)(-?\d+(?:\.\d+)?)(\s*(?:[,)]|\s|$))`)
)

// redactField removes values that must not appear in application logs. It is
// intentionally applied at the logger boundary so callers do not have to
// remember to sanitize every provider, HTTP, or database log field.
func redactField(key string, value any) any {
	if sensitiveKeyPattern.MatchString(key) {
		return redactedValue
	}
	if err, ok := value.(error); ok {
		return redactText(err.Error())
	}
	if key == "sql" || strings.EqualFold(key, "query") {
		if text, ok := value.(string); ok {
			return redactSQL(text)
		}
	}
	if text, ok := value.(string); ok {
		return redactText(text)
	}
	return value
}

func redactText(text string) string {
	text = redactSQL(text)
	text = credentialURLPattern.ReplaceAllString(text, "$1"+redactedValue+":"+redactedValue+"@")
	text = secretQueryPattern.ReplaceAllString(text, "${1}"+redactedValue)
	return structuredSecretPattern.ReplaceAllString(text, "${1}${2}"+redactedValue)
}

func redactSQL(sql string) string {
	sql = sqlStringPattern.ReplaceAllString(sql, "'"+redactedValue+"'")
	return sqlNumericValuePattern.ReplaceAllString(sql, "${1}"+redactedValue+"${3}")
}
