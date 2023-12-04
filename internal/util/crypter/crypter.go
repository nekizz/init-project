package crypter

import (
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/segmentio/ksuid"

	"github.com/nekizz/init-project/pkg/util/crypter"
)

// New initalizes crypter service
func New() *Service {
	return &Service{Service: crypter.New()}
}

// Service holds crypter methods
type Service struct {
	*crypter.Service
}

// XIDWithTime returns unique string ID with 4 bytes of time (seconds) + 3 byte machine id + 2 byte process id + 3 bytes random
// Sample output: b50vl5e54p1000fo3gh0
func (*Service) XIDWithTime(t time.Time) string {
	return xid.NewWithTime(t).String()
}

// UIDWithTime returns unique string ID with 4 bytes of time (seconds) + 3 byte machine id + 2 byte process id + 3 bytes random
// Sample output: b50vl5e54p1000fo3gh0
func (*Service) UIDWithTime(t time.Time) string {
	uid, err := ksuid.NewRandomWithTime(t)
	if err != nil {
		panic(fmt.Sprintf("Couldn't generate KSUID, inconceivable! error: %v", err))
	}
	return uid.String()
}
