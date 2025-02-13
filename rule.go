package httprouter

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const (
	MatchErrorCodeWrongMethod int = iota + 1
	MatchErrorCodeNoMatch
	MatchErrorCodeWrongQueryParam
	ValidationErrorCodeWrongPattern int = iota + 1
)

// RuleMatchError returns an error during matching.
type RuleMatchError struct {
	Code        int
	QueryParams map[string]string
}

// RuleValidationError returns an error during rule validation.
type RuleValidationError struct {
	Code int
}

var (
	MatchErrorWrongMethod       = &RuleMatchError{MatchErrorCodeWrongMethod, nil}
	MatchErrorNoMatch           = &RuleMatchError{MatchErrorCodeNoMatch, nil}
	ValidationErrorWrongPattern = &RuleValidationError{ValidationErrorCodeWrongPattern}
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
	case ValidationErrorCodeWrongPattern:
		return "wrong formed rule pattern"
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
	if len(strings.Trim(r.pattern, " \t\b\r\n")) == 0 {
		return errors.New("attempted to add a rule that was not passed the validation")
	}
	_, err := regexp.Compile(r.pattern)
	if err != nil {
		return fmt.Errorf("%w: %s", ValidationErrorWrongPattern, err.Error())
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
	re, _ := regexp.Compile(`<(\S+)>`)
	pattern = re.ReplaceAllString(pattern, `(?P<$1>\S+)`)
	r := &Rule{pattern, method, h}
	err := r.Validate()
	if err != nil {
		return nil, err
	}
	return r, nil
}
