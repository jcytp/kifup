package auxi

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func NewULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func PInt(i int) *int {
	return &i
}

func PInt64(i int64) *int64 {
	return &i
}
