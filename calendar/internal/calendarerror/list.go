package calendarerror

import(
	"github.com/pkg/errors"
)

// ErrTimeIsBusy is error for case when event time is alredy busy
var ErrTimeIsBusy error = errors.New("Time is busy")
// ErrNotFound is error for case when event not found
var ErrNotFound error = errors.New("Not found")

