package model

import (
	"github.com/elancom/go-util/str"
	"gorm.io/gorm"
)

const (
	TaskStateTodo    = "待运行"
	TaskStateRunning = "运行中"
	TaskStateEnd     = "结束"
)
const (
	TaskResultTodo    = "待定"
	TaskResultError   = "失败"
	TaskResultSuccess = "成功"
)

func NewTask() *Task {
	return &Task{}
}

type Task struct {
	Id      string `json:"id"`      // ID
	State   string `json:"state"`   // 任务状态
	Result  string `json:"result"`  // 任务结果
	Keyword string `json:"keyword"` // 关键词
	Page    int    `json:"page"`    // 当前采集页
	QQGNum  int    `json:"qqgNum"`  // qq群累计
}

func (*Task) TableName() string {
	return "task"
}

type TaskLog struct{}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{db: db}
}

type TaskService struct {
	db *gorm.DB
}

func (t *TaskService) NewTask(task *Task) {
	t.db.Create(task)
}

//func (t *TaskService) Update(task *Task) {
//	t.db.SaveUpdate(task)
//}

func (t *TaskService) UpdatePage(task *Task) {
	t.db.Model(&Task{}).Where("id", task.Id).Update("page", task.Page)
}

func (t *TaskService) UpdateQQGNum(task *Task) {
	t.db.Model(&Task{}).Where("id", task.Id).Update("QQGNum", task.QQGNum)
}

func (t *TaskService) UpdateState(task *Task) {
	t.db.Model(&Task{}).Where("id", task.Id).Update("state", task.State)
}

func (t *TaskService) UpdateResult(task *Task) {
	t.db.Model(&Task{}).Where("id", task.Id).Update("result", task.Result)
}

type ListParam struct {
	State string
}

func (t *TaskService) List(p *ListParam) []*Task {
	db := t.db
	if str.IsNotBlank(p.State) {
		db = db.Where("state = ?", p.State)
	}
	tasks := make([]*Task, 0)
	db.Find(&tasks)
	return tasks
}
