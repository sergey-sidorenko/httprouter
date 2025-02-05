package httprouter

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
)

// RuleContext - context of rule
type RuleContext map[string]string

// Rule - rule for url-to-action routing
type Rule struct {
	pattern string
	method  HttpMethod
	h       http.HandlerFunc
	ctx     RuleContext
}

// Match - проверка, обрабатывает ли правило ресурс, на который указывает URL
func (rule Rule) Match(req *http.Request) bool {
	if req.Method != string(rule.method) {
		return false
	}
	re, _ := regexp.Compile(rule.pattern)
	matches := re.FindStringSubmatch(req.URL.String())
	if matches != nil {
		paramNames := re.SubexpNames()[1:]
		for index, paramValue := range matches[1:] {
			rule.ctx[paramNames[index]] = paramValue
		}
		return true
	}
	return false
}

// Validate - валидация правила
func (r *Rule) Validate() error {
	if len(strings.Trim(r.pattern, " \t\b\r\n")) == 0 {
		return errors.New("attempted to add a rule that was not passed the validation")
	}
	_, err := regexp.Compile(r.pattern)
	if err != nil {
		return errors.New("wrong formed rule pattern")
	}
	if r.method != HttpGetMethod &&
		r.method != HttpPostMethod &&
		r.method != HttpPutMethod &&
		r.method != HttpPatchMethod &&
		r.method != HttpDeleteMethod &&
		r.method != HttpHeadMethod &&
		r.method != HttpOptionsMethod {
		r.method = HttpGetMethod
	}
	return nil
}
func (r Rule) Handle(w http.ResponseWriter, rq *http.Request) {
	r.h(w, rq)
}

// Создание нового правила - валидация правила
func NewRule(pattern string, method HttpMethod, h http.HandlerFunc) *Rule {
	re, _ := regexp.Compile(`{(\S+)}`)
	pattern = re.ReplaceAllString(pattern, `(?P<$1>\S+)`)
	r := &Rule{pattern, method, h, make(RuleContext)}
	err := r.Validate()
	if err != nil {
		return nil
	}
	return r
}
