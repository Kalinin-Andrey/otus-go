package calendar_error

import(
	"github.com/pkg/errors"
)

var DateBusy error = errors.New("Date is busy")
var NotFound error = errors.New("Not found")

