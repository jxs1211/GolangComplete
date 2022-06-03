package unittest

import (
	"strings"
)

type S struct {
	name string `json:"name,omitempty" form:"name,omitempty" orm:"name,omitempty"`
	age  int    `json:"age,omitempty" form:"age,omitempty" orm:"age,omitempty"`
}

func Split(s, sep string) (result []string) {
	i := strings.Index(s, sep)

	for i > -1 {
		result = append(result, s[:i])
		s = s[i+len(sep):] // 这里使用len(sep)获取sep的长度
		i = strings.Index(s, sep)
	}
	result = append(result, s)
	return
}

func Add(a, b int) int {
	return a + b
}

type Learner interface {
	Learn()
	Speaker(s string) string
}

type Student struct{}

func (s *Student) Learn() {
	panic("not implemented") // TODO: Implement
}

func (s *Student) Speaker(str string) string {
	panic("not implemented") // TODO: Implement
}

func ExtractFuncTest(a, b, c int) {
	if a > 0 {
	}
	if b > 0 {
	}
	if c > 0 {
	}
}

func multiply(a, b int) int {
	return a * b
}
