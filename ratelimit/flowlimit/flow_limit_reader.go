package flowlimit

import (
	"context"
	"io"
	"time"

	"golang.org/x/time/rate"
)

type LimitReader struct {
	r       io.Reader
	limiter *rate.Limiter
	ctx     context.Context
}

const burstLimit = 1000 * 1000 * 1000

func NewLimitReader(r io.Reader) *LimitReader {
	return &LimitReader{
		r:   r,
		ctx: context.Background(),
	}
}

func (s *LimitReader) SetRateLimit(bytesPerSec float64) {
	s.limiter = rate.NewLimiter(rate.Limit(bytesPerSec), burstLimit)
	s.limiter.AllowN(time.Now(), burstLimit)
}

func (s *LimitReader) Read(p []byte) (int, error) {
	if s.limiter == nil {
		return s.r.Read(p)
	}
	n, err := s.r.Read(p)
	if err != nil {
		return n, err
	}
	if err := s.limiter.WaitN(s.ctx, n); err != nil {
		return n, err
	}
	return n, nil
}
