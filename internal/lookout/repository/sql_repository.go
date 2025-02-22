package repository

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/lib/pq"

	"github.com/G-Research/armada/pkg/api/lookout"
)

// Emulates JobStates enum
// can't use protobuf enums because gogoproto + grpc-gateway is hard with K8s specific messages
type JobState string

const (
	JobQueued    JobState = "QUEUED"
	JobPending   JobState = "PENDING"
	JobRunning   JobState = "RUNNING"
	JobSucceeded JobState = "SUCCEEDED"
	JobFailed    JobState = "FAILED"
	JobCancelled JobState = "CANCELLED"
	JobDuplicate JobState = "DUPLICATE"
)

type JobRepository interface {
	GetQueueInfos(ctx context.Context) ([]*lookout.QueueInfo, error)
	GetJobSetInfos(ctx context.Context, opts *lookout.GetJobSetsRequest) ([]*lookout.JobSetInfo, error)
	GetJobs(ctx context.Context, opts *lookout.GetJobsRequest) ([]*lookout.JobInfo, error)
}

type SQLJobRepository struct {
	goquDb *goqu.Database
	clock  Clock
}

var (
	// Tables
	jobTable                  = goqu.T("job")
	jobRunTable               = goqu.T("job_run")
	jobRunContainerTable      = goqu.T("job_run_container")
	userAnnotationLookupTable = goqu.T("user_annotation_lookup")

	// Columns: job table
	job_jobId      = goqu.I("job.job_id")
	job_queue      = goqu.I("job.queue")
	job_owner      = goqu.I("job.owner")
	job_jobset     = goqu.I("job.jobset")
	job_priority   = goqu.I("job.priority")
	job_submitted  = goqu.I("job.submitted")
	job_cancelled  = goqu.I("job.cancelled")
	job_job        = goqu.I("job.job")
	job_state      = goqu.I("job.state")
	job_duplicate  = goqu.I("job.duplicate")
	job_jobUpdated = goqu.I("job.job_updated")

	// Columns: job_run table
	jobRun_runId     = goqu.I("job_run.run_id")
	jobRun_jobId     = goqu.I("job_run.job_id")
	jobRun_podNumber = goqu.I("job_run.pod_number")
	jobRun_cluster   = goqu.I("job_run.cluster")
	jobRun_node      = goqu.I("job_run.node")
	jobRun_created   = goqu.I("job_run.created")
	jobRun_started   = goqu.I("job_run.started")
	jobRun_finished  = goqu.I("job_run.finished")
	jobRun_succeeded = goqu.I("job_run.succeeded")
	jobRun_error     = goqu.I("job_run.error")

	// Columns: annotation table
	annotation_jobId = goqu.I("user_annotation_lookup.job_id")
	annotation_key   = goqu.I("user_annotation_lookup.key")
	annotation_value = goqu.I("user_annotation_lookup.value")
)

type JobRow struct {
	JobId     sql.NullString  `db:"job_id"`
	Queue     sql.NullString  `db:"queue"`
	Owner     sql.NullString  `db:"owner"`
	JobSet    sql.NullString  `db:"jobset"`
	Priority  sql.NullFloat64 `db:"priority"`
	Submitted pq.NullTime     `db:"submitted"`
	Cancelled pq.NullTime     `db:"cancelled"`
	JobJson   sql.NullString  `db:"job"`
	State     sql.NullInt64   `db:"state"`
	RunId     sql.NullString  `db:"run_id"`
	PodNumber sql.NullInt64   `db:"pod_number"`
	Cluster   sql.NullString  `db:"cluster"`
	Node      sql.NullString  `db:"node"`
	Created   pq.NullTime     `db:"created"`
	Started   pq.NullTime     `db:"started"`
	Finished  pq.NullTime     `db:"finished"`
	Succeeded sql.NullBool    `db:"succeeded"`
	Error     sql.NullString  `db:"error"`
}

var AllJobStates = []JobState{
	JobQueued,
	JobPending,
	JobRunning,
	JobSucceeded,
	JobFailed,
	JobCancelled,
	JobDuplicate,
}

var JobStateToIntMap = map[JobState]int{
	JobQueued:    1,
	JobPending:   2,
	JobRunning:   3,
	JobSucceeded: 4,
	JobFailed:    5,
	JobCancelled: 6,
	JobDuplicate: 7,
}

var IntToJobStateMap = map[int]JobState{
	1: JobQueued,
	2: JobPending,
	3: JobRunning,
	4: JobSucceeded,
	5: JobFailed,
	6: JobCancelled,
	7: JobDuplicate,
}

var defaultQueryStates = []JobState{
	JobQueued,
	JobPending,
	JobRunning,
	JobSucceeded,
	JobFailed,
	JobCancelled,
}

func NewSQLJobRepository(db *goqu.Database, clock Clock) *SQLJobRepository {
	return &SQLJobRepository{goquDb: db, clock: clock}
}
