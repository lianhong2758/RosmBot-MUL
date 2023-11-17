// Package wallet 货币系统
package walle

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/FloatTech/floatbox/file"
	sql "github.com/FloatTech/sqlite"
)

// Storage 货币系统
type Storage struct {
	sync.RWMutex
	db *sql.Sqlite
}

// Wallet 钱包
type Wallet struct {
	UID   string
	Money int
}

var (
	swdb = &Storage{
		db: &sql.Sqlite{
			DBPath: "data/wallet/wallet.db",
		},
	}
)

func init() {
	if file.IsNotExist("data/wallet") {
		err := os.MkdirAll("data/wallet", 0755)
		if err != nil {
			panic(err)
		}
	}
	err := swdb.db.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
	err = swdb.db.Create("storage", &Wallet{})
	if err != nil {
		panic(err)
	}
}

// GetWalletOf 获取钱包数据
func GetWalletOf(uid string) (money int) {
	return swdb.getWalletOf(uid).Money
}

// GetWalletInfoGroup 获取多人钱包数据
//
// if sort == true,由高到低排序; if sort == false,由低到高排序
func GetGroupWalletOf(uids []string, sortable bool) (wallets []Wallet, err error) {
	return swdb.getGroupWalletOf(uids, sortable)
}

// InsertWalletOf 更新钱包(money > 0 增加,money < 0 减少)
func InsertWalletOf(uid string, money int) error {
	lastMoney := swdb.getWalletOf(uid)
	return swdb.updateWalletOf(uid, lastMoney.Money+money)
}

// 获取钱包数据
func (sql *Storage) getWalletOf(uid string) (Wallet Wallet) {
	sql.RLock()
	defer sql.RUnlock()
	uidstr := uid
	_ = sql.db.Find("storage", &Wallet, "where uid is "+uidstr)
	return
}

// 获取钱包数据组
func (sql *Storage) getGroupWalletOf(uids []string, issorted bool) (wallets []Wallet, err error) {
	sql.RLock()
	defer sql.RUnlock()
	wallets = make([]Wallet, 0, len(uids))
	sort := "ASC"
	if issorted {
		sort = "DESC"
	}
	info := Wallet{}
	err = sql.db.FindFor("storage", &info, "where uid IN ("+strings.Join(uids, ", ")+") ORDER BY money "+sort, func() error {
		wallets = append(wallets, info)
		return nil
	})
	return
}

// 更新钱包
func (sql *Storage) updateWalletOf(uid string, money int) (err error) {
	sql.Lock()
	defer sql.Unlock()
	return sql.db.Insert("storage", &Wallet{
		UID:   uid,
		Money: money,
	})
}
