package data

import "database/sql"

type Models struct {
	Tasks                TaskModel
	TaskSessions         TaskSessionModel
	TaskSessionIntervals TaskSessionIntervalModel
	Commits              CommitModel
}

func NewModels(DB *sql.DB) Models {
	return Models{
		Tasks:                TaskModel{DB},
		TaskSessions:         TaskSessionModel{DB},
		TaskSessionIntervals: TaskSessionIntervalModel{DB},
		Commits:              CommitModel{DB},
	}
}
