package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Exchange struct {
	Id            uint        `json:"id" gorm:"primaryKey"`
	Confirmations Uint64Array `json:"confirmations"`
	IsCompleted   bool        `json:"is_completed"`
	ChatRoomID    uint        `json:"chat_room_id"`
	ChatRoom      ChatRoom    `json:"chat_room" gorm:"foreignKey:ChatRoomID"`
}

type Uint64Array []uint64

func (a *Uint64Array) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Uint64Array: %v", value)
	}
	return json.Unmarshal(bytes, a)
}

func (a Uint64Array) Value() (driver.Value, error) {
	return json.Marshal(a)
}
