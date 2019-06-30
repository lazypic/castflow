package main

import (
	"fmt"
)

// Character 는 캐릭터 정보를 다루는 자료구조이다.
type Character struct {
	ID              string // 캐릭터이름
	Regnum          string // 저작권 등록번호
	Manager         string // 매니저
	FieldOfActivity string // 활동범위
	Concept         string // 컨셉
	StartDate       string // 개발 시작일
	Email           string // 캐릭터 이메일
}

func (c Character) String() string {
	return fmt.Sprintf(`
ID: %s (%s)
Manager: %s
FieldOfActivity: %s
Concept: %s
StartDate: %s
Email: %s`,
		c.ID,
		c.Regnum,
		c.Manager,
		c.FieldOfActivity,
		c.Concept,
		c.StartDate,
		c.Email,
	)
}
