package models

import "github.com/microcosm-cc/bluemonday"

func (req *RequestUpdateBalance) Sanitize(sanitizer bluemonday.Policy) {
}
