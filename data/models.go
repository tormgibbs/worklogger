package data

import "database/sql"

type Models struct {
	Tasks                TaskModel
	SessionTags          SessionTagModel
	SessionKPI           SessionKPIModel
	TaskSessions         TaskSessionModel
	TaskSessionIntervals TaskSessionIntervalModel
	Commits              CommitModel
	Logs                 LogModel
}

func NewModels(DB *sql.DB) Models {
	return Models{
		Tasks:                TaskModel{DB},
		TaskSessions:         TaskSessionModel{DB},
		TaskSessionIntervals: TaskSessionIntervalModel{DB},
		Commits:              CommitModel{DB},
		SessionTags:          SessionTagModel{DB},
		SessionKPI:           SessionKPIModel{DB},
		Logs:                 LogModel{DB},
	}
}
