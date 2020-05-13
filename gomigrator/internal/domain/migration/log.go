package migration

import (
	"time"

	"github.com/Kalinin-Andrey/otus-go/gomigrator/pkg/sqlmigrator/api"
)

// MigrationLog struct
type MigrationLog struct {
	ID				uint
	Status			uint
	Name			string
	Time			time.Time
}


type MigrationsLogsList map[uint]MigrationLog
type MigrationsLogsSlice []MigrationLog


const (
	// TableName const
	TableName        = "migration"
	StatusNotApplied = 0
	StatusApplied    = 1
	StatusError      = 2
)

var SQLCreateTable string = `CREATE TABLE IF NOT EXISTS public."migration" (
	id int4 NOT NULL,
	status int4 NOT NULL DEFAULT 0,
	name varchar(100) NOT NULL,
	"time" timestamptz NOT NULL DEFAULT Now(),
	CONSTRAINT migration_pkey PRIMARY KEY (id)
);`

// QueryCondition struct for defining a query condition
type QueryCondition struct {
	Where	*WhereCondition
}
// WhereCondition struct
type WhereCondition struct {
	Status	uint
}


func (l MigrationsLogsList) GetSlice() (mls []MigrationLog) {
	mls = make([]MigrationLog, 0, len(l))

	for _, ml := range l {
		mls = append(mls, ml)
	}

	return mls
}

func (l MigrationsLogsList) GetIDs() (ids []int) {
	ids = make([]int, 0, len(l))

	for id := range l {
		ids = append(ids, int(id))
	}
	return ids
}

func (l MigrationsLogsList) Copy() (mls MigrationsLogsList) {
	mls = make(MigrationsLogsList, len(l))

	for id, ml := range l {
		mls[id] = ml
	}

	return mls
}

func GroupLogsByStatus(list []MigrationLog) (l map[uint]MigrationsLogsList) {
	l = make(map[uint]MigrationsLogsList, 2)
	l[StatusNotApplied] = make(MigrationsLogsList, 0)
	l[StatusApplied] = make(MigrationsLogsList, 0)

	for _, i := range list {
		if i.Status == StatusApplied {
			l[StatusApplied][i.ID] = i
		} else {
			l[StatusNotApplied][i.ID] = i
		}
	}

	return l
}


func MigrationsListFilterExceptByKeys(sourceList api.MigrationsList, exceptList MigrationsLogsList) (l api.MigrationsList) {
	l = make(api.MigrationsList, 0)

	for id, m := range sourceList {
		if _, ok := exceptList[id]; !ok {
			l[id] = m
		}
	}

	return l
}


func MigrationsListFilterExistsByKeys(sourceList api.MigrationsList, exceptList MigrationsLogsList) (l api.MigrationsList) {
	l = make(api.MigrationsList, 0)

	for id, m := range sourceList {
		if _, ok := exceptList[id]; ok {
			l[id] = m
		}
	}

	return l
}


func MigrationsLogsFilterExceptByKeys(sourceList MigrationsLogsList, exceptList MigrationsLogsList) (l MigrationsLogsList) {
	l = make(MigrationsLogsList, 0)

	for id, m := range sourceList {
		if _, ok := exceptList[id]; !ok {
			l[id] = m
		}
	}

	return l
}


func MigrationsLogsFilterExistsByKeys(sourceList MigrationsLogsList, exceptList MigrationsLogsList) (l MigrationsLogsList) {
	l = make(MigrationsLogsList, 0)

	for id, m := range sourceList {
		if _, ok := exceptList[id]; ok {
			l[id] = m
		}
	}

	return l
}


func NewMigrationLog (m api.Migration, status uint) *MigrationLog {
	return &MigrationLog{
		ID:		m.ID,
		Status:	status,
		Name:	m.Name,
	}
}


func (s MigrationsLogsSlice) Len() int {
	return len(s)
}

func (s MigrationsLogsSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s MigrationsLogsSlice) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

