package model

import (
	"github.com/elancom/go-util/str"
	"gorm.io/gorm"
	"log"
)

// 机器ID
// 服务器地址

type Machine struct {
	Id     string `json:"id"`
	Url    string `json:"url"`
	Online bool   `json:"online"`
}

func (*Machine) TableName() string {
	return "machine"
}

func NewMachineStore(machine *Machine, dao *MachineDao) *MachineStore {
	return &MachineStore{Machine: machine, dao: dao}
}

type MachineStore struct {
	*Machine
	dao *MachineDao
}

func (store *MachineStore) UpdateOnline(b bool) {
	if store.Online == b {
		return
	}
	store.Online = b
	go store.dao.UpdateOnline(store.Id, store.Online)
}

func NewMachineDao(gdb *gorm.DB) *MachineDao {
	return &MachineDao{db: gdb}
}

type MachineDao struct {
	db *gorm.DB
}

func (dao *MachineDao) Save(value *Machine) {
	dao.db.Save(value)
}
func (dao *MachineDao) UpdateOnline(id string, b bool) {
	dao.db.Model(&Machine{}).Where("id", id).Update("online", b)
}

func (dao *MachineDao) List() []Machine {
	list := make([]Machine, 0)
	dao.db.Find(&list)
	return list
}

func NewMachineStoreManager(dao *MachineDao) *MachineStoreManager {
	return &MachineStoreManager{dao: dao}
}

type MachineStoreManager struct {
	dao    *MachineDao
	stores []*MachineStore
}

func (mgr *MachineStoreManager) Init() {
	mgr.Reload()
}

func (mgr *MachineStoreManager) Reload() {
	machines := mgr.dao.List()
	stores := make([]*MachineStore, 0, len(machines))
	for _, machine := range machines {
		stores = append(stores, NewMachineStore(&machine, mgr.dao))
	}
	mgr.stores = stores
	log.Println("load stores size:" + str.String(len(mgr.stores)))
}

func (mgr *MachineStoreManager) Get(id string) *MachineStore {
	for _, s := range mgr.stores {
		if s.Id == id {
			return s
		}
	}
	return nil
}
