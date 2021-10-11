package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrDataNotExists = errors.New("object does not exist")
	ErrDataAlreadyExists = errors.New("object already exist")
)

var (
	normal  = 1
	invalid = 2
)

type Account struct {
	Uid          int32     `gorm:"column:uid;primary_key:true;AUTO_INCREMENT;comment:'用户uid'" json:"uid"`
	Username     string    `gorm:"column:username;type:varchar(140);NOT NULL;unique;comment:'用户名'" json:"username"`
	Password     string    `gorm:"column:password;type:varchar(255);comment:'密码'" json:"password"`
	Salt         string    `gorm:"column:salt;type:varchar(80);comment:'盐'" json:"-"`
	Nickname     string    `gorm:"column:nickname;type:varchar(30);comment:'昵称'" json:"nickname"`
	OrderWord    string    `gorm:"column:order_word;type:varchar(140);default:'#';index:IDX_WORD;comment:'排序字段'" json:"order_word"`
	Avatar       string    `gorm:"column:avatar;type:varchar(255);comment:'头像'" json:"avatar"`
	IsFirstLogin int8      `gorm:"column:is_first_login;type:tinyint(4);default:'0';comment:'首次登录 1-是 2-否'" json:"is_first_login"`
	DataState    uint8     `gorm:"column:data_state;type:tinyint(4);default:'1';comment:'数据状态 1-正常 2-删除'" json:"-"`
	CreateTime   time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"-"`
	UpdateTime   time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"-"`
	GoogleId     string    `gorm:"column:google_id;type:varchar(255);comment:'Google账号ID'" json:"google_id"`
	AppleId      string    `gorm:"column:apple_id;type:varchar(255);comment:'Apple账号ID'" json:"apple_id"`
}

// gorm自动建表
func (a *Account) TableName() string {
	return "account"
}

// 表名
func (a *Account) table() string {
	return "account"
}

// 单个创建
func (a *Account) create(db *gorm.DB) error {
	if err := db.Table(a.table()).Create(a).Error; err != nil {
		return err
	}
	return nil
}

// 单个获取
func (a *Account) get(db *gorm.DB) error {
	if err := db.Table(a.table()).Where("uid = ? and data_state= ? ",a.Uid, normal).Scopes(a.username).
		First(&a).Error; err != nil {
		return err
	}
	return nil
}

// 分页获取数据
func (a *Account) getPage(db *gorm.DB, offset, limit int32, like, order string, desc bool) ([]*Account, int32, error) {
	pageAccount := make([]*Account, limit)
	var count int32
	order = a.order(order, desc)
	if err := db.Table(a.table()).Where("data_state = ?", normal).Count(&count).Order(order, desc).
		Offset(offset).Limit(limit).Find(&pageAccount).Error; err != nil {
			return nil, 0, err
	}
	return pageAccount, count, nil
}

// 单个修改-修改全部信息
func (a *Account) update(db *gorm.DB) error {
	if err := db.Table(a.table()).Where("uid = ? ", a.Uid).Update(a).Error; err != nil {
		return err
	}
	return nil
}
// 修改部分信息
func (a *Account) updateT(db *gorm.DB) error {
	if err := db.Table(a.table()).Where("uid = ?", a.Uid).Update(&Account{
		Username: a.Username,
		Nickname: a.Nickname,
		Avatar: a.Avatar,
	}).Error; err != nil {
		return err
	}
	return nil
}


// 真实删除-删除数据
func (a *Account) delete(db *gorm.DB) error {
	if err := db.Table(a.table()).Where("uid = ?", a.Uid).Delete(a).Error; err != nil {
		return err
	}
	return nil
}

// 物理删除-修改数据状态为不可用
func (a *Account) deleteT(db *gorm.DB) error {
	if err := db.Table(a.table()).Where("uid = ? ", a.Uid).Update(map[string]interface{}{
		"data_state": invalid,
	}).Error; err != nil {
		return err
	}
	return nil
}

// 动态条件查询
func (a *Account) username(db *gorm.DB) *gorm.DB {
	if a.Username != "" {
		return db.Where("username = ?", a.Username)
		//db = db.Where("username = ?", a.Username)
	}
	return db
}

func (a *Account) order(column string, desc bool) string {
	var columns = [2]string{"create_time", "update_time"}
	for _, v := range columns {
		if v == column {
			if desc {
				return column + " desc"
			} else {
				return column
			}
		}
	}
	return "uid"
}