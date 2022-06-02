package model

import (
	"gorm.io/gorm"
	"time"
)

type Buyer struct {
	Id       int    `json:"id"`       // ID
	Username string `json:"username"` // 用户名
	QQGNum   int    `json:"qqgNum"`   // 已购买数量
}

func (*Buyer) TableName() string {
	return "buyer"
}

func NewBuyQQG() *BuyQQG {
	return &BuyQQG{}
}

type BuyQQG struct {
	BuyerId  int        // 购买者
	Code     string     // QQ群
	DateTime *time.Time // 购买时间
}

func (*BuyQQG) TableName() string {
	return "buy_qqg"
}

type BuyService struct {
	db *gorm.DB
}

func (s *BuyService) Add(buyer *Buyer, qqg *QQG) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 中间表
		now := time.Now()
		buyQQG := NewBuyQQG()
		buyQQG.BuyerId = buyer.Id
		buyQQG.Code = qqg.Code
		buyQQG.DateTime = &now
		s.db.Create(buyQQG)
		// 购买累计
		s.db.Model(&Buyer{}).Update("QQGNum", gorm.Expr("QQGNum + ?", 1))
		// 出售次数
		s.db.Model(&QQG{}).Update("SellTimes", gorm.Expr("SellTimes + ?", 1))
		return nil
	})
	if err != nil {
		return
	}
}
