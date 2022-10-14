package shttp

import (
	"net/url"
	"sort"
	"strings"
)

type Values map[string][]string

func (v Values) Get(key string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func (v Values) Set(key, value string) {
	v[key] = []string{value}
}

func (v Values) Add(key, value string) {
	v[key] = append(v[key], value)
}

func (v Values) Del(key string) {
	delete(v, key)
}

func (v Values) Has(key string) bool {
	_, ok := v[key]
	return ok
}

func (v Values) Encode() string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}
