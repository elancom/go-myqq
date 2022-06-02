package model

import (
	"gorm.io/gorm"
	"time"
)

type QQ struct {
	Id             int        `json:"id"`             // ID
	MachineId      string     `json:"machineId"`      // 所在机器ID
	Code           string     `json:"code"`           // QQ号
	SearchLastTime *time.Time `json:"searchLastTime"` // 最后搜索时间
}

func (*QQ) TableName() string {
	return "qq"
}

// QQSortWithSearch 按搜索时间顺序排序
type QQSortWithSearch []*QQ

func (sq QQSortWithSearch) Len() int {
	return len(sq)
}

func (sq QQSortWithSearch) Less(i, j int) bool {
	if sq[i].SearchLastTime == nil {
		return true
	}
	if sq[j].SearchLastTime == nil {
		return false
	}
	return sq[i].SearchLastTime.UnixMilli() < sq[j].SearchLastTime.UnixMilli()
}

func (sq QQSortWithSearch) Swap(i, j int) {
	sq[i], sq[j] = sq[j], sq[i]
}

func NewQQG() *QQG {
	return &QQG{}
}

type QQG struct {
	Id int `json:"id"` // ID
	// 基础信息
	Code         string `json:"code"`         // 群号
	Gid          string `json:"gid"`          // 群号
	Class        string `json:"class"`        // 类别
	ClassId      int    `json:"classId"`      // 类别ID
	Tags         string `json:"tags"`         // 分类标签
	Features     string `json:"features"`     // 状态特征
	Labels       string `json:"labels"`       // 群标签
	MaxMemberNum int    `json:"maxMemberNum"` // 最大会员数
	MemberNum    int    `json:"memberNum"`    // 会员数
	Memo         string `json:"memo"`         // 群描述
	Name         string `json:"name"`         // 群名称
	Image        string `json:"image"`        // 群头像
	CityId       int    `json:"cityId"`       // 城市ID
	Province     string `json:"province"`     // 省份
	City         string `json:"city"`         // 城市
	Level        int    `json:"level"`        // 等级
	Activity     int    `json:"activity"`     // 活动情况?
	// 采集信息
	UpdateTimes int        `json:"updateTimes"` // 更新次数
	UpdateTime  *time.Time `json:"updateTime"`  // 更新时间
	DateTime    *time.Time `json:"dateTime"`    // 创建时间
	SellTimes   int        `json:"sellTimes"`   // 出售次数
}

func (*QQG) TableName() string {
	return "qqg"
}

func NewQQDao(gdb *gorm.DB) *QQDao {
	return &QQDao{db: gdb}
}

type QQDao struct {
	db *gorm.DB
}

func (dao *QQDao) Create(qq *QQ) {
	dao.db.Create(qq)
}

func (dao *QQDao) UpdateLastSearchTime(id int, time *time.Time) {
	dao.db.Model(&QQ{}).Where("id", id).Update("SearchLastTime", time)
}

func (dao *QQDao) List() []*QQ {
	qqs := make([]*QQ, 0)
	dao.db.Find(&qqs)
	return qqs
}

func NewQQService(dao *QQDao) *QQService {
	return &QQService{dao: dao}
}

type QQService struct {
	dao *QQDao
}

func (s *QQService) Add(qq *QQ) {
	s.dao.Create(qq)
}

func (s *QQService) UpdateLastSearchTime(id int, time *time.Time) {
	s.dao.UpdateLastSearchTime(id, time)
}

func (s *QQService) List() []*QQ {
	return s.dao.List()
}

func NewQQGDao(db *gorm.DB) *QQGDao {
	return &QQGDao{db: db}
}

type QQGDao struct {
	db *gorm.DB
}

func (dao *QQGDao) Create(qqg *QQG) {
	dao.db.Create(qqg)
}

func (dao *QQGDao) FindByCode(code string) (*QQG, error) {
	one, err := One[QQG](dao.db.Where("Code = ?", code).First(NewQQG()))
	if err != nil {
		return nil, err
	}
	return one, nil
}

func (dao *QQGDao) Save(qqg *QQG) {
	dao.db.Save(qqg)
}

func NewQQGService(dao *QQGDao) *QQGService {
	return &QQGService{dao: dao}
}

type QQGService struct {
	dao *QQGDao
}

func (s *QQGService) Add(qqg *QQG) {
	s.dao.Create(qqg)
}

func (s *QQGService) FindByCode(code string) (*QQG, error) {
	return s.dao.FindByCode(code)
}

func (s *QQGService) Save(qqg *QQG) {
	s.dao.Save(qqg)
}
