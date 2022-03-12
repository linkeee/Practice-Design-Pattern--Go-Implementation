package output

import (
	"demo/db"
	"demo/monitor/plugin"
	"demo/monitor/record"
	"fmt"
	"reflect"
)

type MemoryDbOutput struct {
	db        db.Db
	tableName string
}

func (m *MemoryDbOutput) Install() {
	m.db = db.MemoryDbInstance()
	table := db.NewTable(m.tableName).WithType(reflect.TypeOf(new(record.MonitorRecord)))
	m.db.CreateTableIfNotExist(table)
}

func (m *MemoryDbOutput) Uninstall() {
}

func (m *MemoryDbOutput) SetContext(ctx plugin.Context) {
	if name, ok := ctx.GetString("tableName"); ok {
		m.tableName = name
	}
}

func (m *MemoryDbOutput) Output(event *plugin.Event) error {
	r, ok := event.Payload().(*record.MonitorRecord)
	if !ok {
		return fmt.Errorf("memory db output unknown event type %T", event.Payload())
	}
	return m.db.Insert(m.tableName, r.Id, r)
}
