package model

import (
	"time"
)

type UserModel struct {
	ID        uint       `gorm:"primaryKey"` // Standard field for the primary key
	Name      string     // A regular string field
	Email     string     // A pointer to a string, allowing for null values
	Age       uint8      // An unsigned 8-bit integer
	Birthday  *time.Time // A pointer to time.Time, can be null
	CreatedAt time.Time  // Automatically managed by GORM for creation time
	UpdatedAt time.Time  // Automatically managed by GORM for update time
}

func (u *UserModel) TableName() string {
	return "trade_user"
}

type TradeOrderModel struct {
	ID          uint `gorm:"primaryKey"` // Standard field for the primary key
	UserId      uint
	OrderNo     uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	TotalAmount uint
	PaiedAmount uint
}

func (t *TradeOrderModel) TableName() string {
	return "trade_order"
}

type OrderProductModel struct {
	ID          uint `gorm:"primaryKey"` // Standard field for the primary key
	OrderId     uint
	ProductId   uint
	Quantity    uint
	TotalAmount uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (o *OrderProductModel) TableName() string {
	return "trade_order_product"
}

type ProductModel struct {
	ID        uint `gorm:"primaryKey"` // Standard field for the primary key
	Name      string
	Price     uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (o *ProductModel) TableName() string {
	return "trade_product"
}
