package dao

import (
	"context"
	"database/sql"
	"gindemo/webook/payment/domain"
	"time"
)

type PaymentDAO interface {
	Insert(ctx context.Context, pmt Payment) error
	UpdateTxnIDAndStatus(ctx context.Context, bizTradeNO, txnID string, status domain.PaymentStatus) error
	FindExpiredPayment(ctx context.Context, offset, limit int, t time.Time) ([]Payment, error)
	GetPayment(ctx context.Context, bizTradeNo string) (Payment, error)
}

type Payment struct {
	Id  int64 `gorm:"primary_key;autoIncrement" bson:"id,omitempty"`
	Amt int64
	// 你存储枚举也可以，比如说 0-CNY
	// 目前磁盘内存便宜，直接放string也可以
	// "CNY"
	Currency string
	// 可以抽象地认为，这是一个简短的描述
	// 也就是说即便是别的支付方式，这边也可以提供一个简单的描述
	// 你可以认为这算是冗余的数据，因为从原则上来说，我们可以完全不保存的。
	// 而是要求调用者直接 BizID 和 Biz 去找业务方要
	// 管的越少，系统越稳
	Description string `gorm:"column:biz_trade_no;type:varchar(256);unique"`

	// 业务方传过来的
	BizTradeNO string `gorm:"column:biz_trade_no;type:varchar(256);unique"`

	// 第三方支付平台的事务 ID，唯一的
	TxnID sql.NullString `gorm:"column:txn_id;type:varchar(128);unique"`

	Status uint8
	Utime  int64
	Ctime  int64
}
