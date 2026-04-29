package traceunic

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateTraceID(prefix string) string {
	date := time.Now().Format("20060102")
	uid := uuid.New().String()[:8]
	return fmt.Sprintf("%s_%s_%s", prefix, date, uid)
}
