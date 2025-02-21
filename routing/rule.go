package routing

import (
	"fmt"
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

// PathParams - context of rule
type PathParams map[string]string

// Rule - rule for url-to-action routing
type Rule struct {
	pattern string
	method  string
}

// Match - проверка, обрабатывает ли правило ресурс, на который указывает URL
func (rule *Rule) Match(req *Request) (PathParams, error) {
	var params PathParams
	if req.Method != rule.method {
		return params, fmt.Errorf("%w", MatchErrorWrongMethod)
	}
	re, _ := regexp.Compile(rule.pattern)
	var matches []string = re.FindStringSubmatch(req.URL.String())
	if matches != nil {
		var paramNames []string = re.SubexpNames()[1:]
		if len(paramNames) > 0 {
			params = make(PathParams)
			for index, paramValue := range paramNames {
				params[paramNames[index]] = paramValue

			}
		}
		return params, nil
	}
	return params, fmt.Errorf("%w", MatchErrorNoMatch)
}

// Validate - validates the rule
func (r *Rule) Validate() error {
	var trims = " \t\b\r\n"
	r.pattern = strings.Trim(r.pattern, trims)
	if len(r.pattern) == 0 {
		r.pattern = "/"
	} else if r.pattern[0] != '/' {
		r.pattern = "/" + r.pattern
	}
	// sorting alternative variants in OR clause of regexp
	rePattern := `\(.*?\|.*?\)`
	re, _ := regexp.Compile(rePattern)
	if re != nil {
		r.pattern = re.ReplaceAllStringFunc(r.pattern, func(s string) string {
			s = strings.Trim(s, "()")
			variants := strings.Split(s, "|")
			slices.Sort(variants)
			return fmt.Sprintf("(%s)", strings.Join(variants, "|"))
		})
	}
	re, _ = regexp.Compile(`<(\S+)>:?([^\s\/]*)`)
	if re != nil {
		// clear all named groups with empty regexp expressions
		reEmptyNamedGroups, _ := regexp.Compile(`(<\S+>):?\s*((?:$|\/).*)`)
		if reEmptyNamedGroups != nil {
			r.pattern = reEmptyNamedGroups.ReplaceAllString(r.pattern, `$1:\S+$2`)
		}
		// make all groups in pattern non-capturing
		reParentesses, _ := regexp.Compile(`(<\S+>:)\(([^\s\/]*)\)`)
		if reParentesses != nil {
			r.pattern = reParentesses.ReplaceAllString(r.pattern, `$1(?:$2)`)
		}
		r.pattern = re.ReplaceAllString(r.pattern, `(?P<$1>$2)(?:$|\/)`)
	}
	_, err := regexp.Compile(r.pattern)
	if err != nil {
		return fmt.Errorf("%w: %s", RuleValidationErrorWrongPattern, err.Error())
	}
	r.method = strings.Trim(r.method, trims)
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

// Создание нового правила - валидация правила
func NewRule(pattern string, method string) (*Rule, error) {
	r := &Rule{pattern, method}
	err := r.Validate()
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Pattern - returns pattern
func (r Rule) Pattern() string {
	return r.pattern
}

// Method - returns http request method that corresponds to the rule
func (r Rule) Method() string {
	return r.method
}
