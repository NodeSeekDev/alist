package rateg

import (
	"context"
	"errors"
	"io"
	"time"
)

func Copy(dst io.Writer, src io.Reader, rates int64) (written int64, err error) {
	return CopyCtx(context.Background(), dst, src, rates)
}

// CopyCtx copies from src to dst, max rates bytes per second, support context cancellation
func CopyCtx(ctx context.Context, dst io.Writer, src io.Reader, rates int64) (written int64, err error) {
	buf := make([]byte, rates)
outer:
	for {
		select {
		case <-ctx.Done():
			return written, ctx.Err()
		default:
			start := time.Now()
			nr, er := src.Read(buf)
			if nr > 0 {
				nw, ew := dst.Write(buf[0:nr])
				if nw < 0 || nr < nw {
					nw = 0
					if ew == nil {
						ew = errors.New("invalid write result")
					}
				}
				written += int64(nw)
				if ew != nil {
					err = ew
					break
				}
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break outer
			}
			elapsed := time.Since(start)
			if elapsed < time.Second {
				time.Sleep(time.Second - elapsed)
			}
		}
	}
	return written, err
}
