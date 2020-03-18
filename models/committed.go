package models

import "time"

type Committed struct {
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
	UpdatedBy string    `json:"updatedBy"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"updated"`
}

func (Committed) newCommitted(userName string) Committed {
	return Committed{
		CreatedBy: userName,
		UpdatedBy: userName,
	}
}
