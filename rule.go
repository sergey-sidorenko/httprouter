package httprouter

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
)

const (
	MatchErrorCodeWrongMethod int = iota + 1
	MatchErrorCodeNoMatch
	MatchErrorCodeWrongQueryParam
	RuleValidationErrorCodeWrongPattern int = iota + 1
	RuleValidationErrorCodeWrongField
)

// RuleMatchError returns an error during matching.
type RuleMatchError struct {
	Code        int
	QueryParams map[string]string
}

// RuleValidationError returns an error during rule validation.
type RuleValidationError struct {
	Code   int
	Fields map[string]string
}

var (
	MatchErrorWrongMethod           = &RuleMatchError{MatchErrorCodeWrongMethod, nil}
	MatchErrorNoMatch               = &RuleMatchError{MatchErrorCodeNoMatch, nil}
	RuleValidationErrorWrongPattern = &RuleValidationError{RuleValidationErrorCodeWrongPattern, nil}
)

func (e RuleMatchError) Error() string {
	switch e.Code {
	case MatchErrorCodeWrongMethod:
		return "http method doesnt math the rule"
	case MatchErrorCodeNoMatch:
		return "no match"
	case MatchErrorCodeWrongQueryParam:
		var qp string
		for pn, et := range e.QueryParams {
			qp += fmt.Sprintf("%s: %s", pn, et)
		}
		return fmt.Sprintf("query param %s does not math the rule method doesnt math the rule", qp)
	default:
		return ""
	}
}
func (e RuleValidationError) Error() string {
	switch e.Code {
	case RuleValidationErrorCodeWrongPattern:
		return "wrong formed rule pattern"
	case RuleValidationErrorCodeWrongField:
		var qp string
		for fn, err := range e.Fields {
			qp += fmt.Sprintf("[%s]: %s", fn, err)
		}
		return fmt.Sprintf("these fields of the rule filled wrong: %s", qp)
	default:
		return "unknown validation error"
	}
}

// RuleContext - context of rule
type RuleContext map[string]string

// Rule - rule for url-to-action routing
type Rule struct {
	pattern string
	method  string
	h       http.HandlerFunc
}

// Match - проверка, обрабатывает ли правило ресурс, на который указывает URL
func (rule Rule) Match(req *http.Request) (RuleContext, error) {
	var ctx RuleContext
	if req.Method != rule.method {
		return ctx, fmt.Errorf("%w", MatchErrorWrongMethod)
	}
	re, _ := regexp.Compile(rule.pattern)
	var matches []string = re.FindStringSubmatch(req.URL.String())
	if matches != nil {
		var paramNames []string = re.SubexpNames()[1:]
		if len(paramNames) > 0 {
			ctx = make(RuleContext)
			for index, paramValue := range paramNames {
				ctx[paramNames[index]] = paramValue

			}
		}
		return ctx, nil
	}
	return ctx, fmt.Errorf("%w", MatchErrorNoMatch)
}

// Validate - валидация правила
func (r *Rule) Validate() error {
	var trims = " \t\b\r\n"
	r.method = strings.Trim(r.method, trims)
	r.pattern = strings.Trim(r.pattern, trims)
	if len(r.pattern) == 0 {
		ruleValidationError := &RuleValidationError{RuleValidationErrorCodeWrongField, map[string]string{"pattern": "empty"}}
		return fmt.Errorf("%w", ruleValidationError)
	}
	if r.pattern[0] != '/' {
		r.pattern = "/" + r.pattern
	}
	_, err := regexp.Compile(r.pattern)
	if err != nil {
		return fmt.Errorf("%w: %s", RuleValidationErrorWrongPattern, err.Error())
	}
	if r.method != "GET" &&
		r.method != "POST" &&
		r.method != "PUT" &&
		r.method != "PATCH" &&
		r.method != "DELETE" &&
		r.method != "HEAD" &&
		r.method != "OPTIONS" {
		r.method = "GET"
	}
	return nil
}

// SetPattern - sets pattern
func (r *Rule) SetPattern(p string) {
	r.pattern = p
}

// Handle - shandles client network request
func (r Rule) Handle(w http.ResponseWriter, rq *http.Request, ctx RuleContext) {
	r.h(w, rq)
}

// Создание нового правила - валидация правила
func NewRule(pattern string, method string, h http.HandlerFunc) (*Rule, error) {

	rePattern := `{.*?\|.*?}`
	re, _ := regexp.Compile(rePattern)
	if re != nil {
		pattern = re.ReplaceAllStringFunc(pattern, func(s string) string {
			s = strings.Trim(s, "{}")
			variants := strings.Split(s, "|")
			slices.Sort(variants)
			return fmt.Sprintf("{%s}", strings.Join(variants, "|"))
		})
	}
	re, _ = regexp.Compile(`<(\S+)>`)
	if re != nil {
		pattern = re.ReplaceAllString(pattern, `(?P<$1>\S+)`)
	}
	r := &Rule{pattern, method, h}
	err := r.Validate()
	if err != nil {
		return nil, err
	}
	return r, nil
}
