package collection

import (
	"easyrpc/timex"
	"sync"
	"time"
)

//滑动时间窗口算法

type Bucket struct {
	Sum   float64 //样本窗口的值
	Count int64   //样本窗口增加的次数
}

func (b *Bucket) add(v float64) {
	b.Sum += v
	b.Count++
}

//重置样本窗口,样本窗口过期时
func (b *Bucket) reset() {
	b.Sum = 0
	b.Count = 0
}

//window 滑动窗口
type window struct {
	buckets []*Bucket //样本窗口
	size    int       //样本窗口个数
}

func newWindow(size int) *window {
	buckets := make([]*Bucket, size)
	for i := 0; i < size; i++ {
		buckets[i] = new(Bucket)
	}
	return &window{
		buckets: buckets,
		size:    size,
	}
}

func (w *window) add(offset int, v float64) {
	w.buckets[offset%w.size].add(v)
}

func (w *window) reduce(start, count int, fn func(b *Bucket)) {
	for i := 0; i < count; i++ {
		fn(w.buckets[(start+i)%w.size])
	}
}

func (w *window) resetBucket(offset int) {
	w.buckets[offset%w.size].reset()
}

type (
	RollingWindowOption func(rollingwindow *RollingWindow)

	//RollingWindow 滑动时间窗口
	//
	//
	//
	//
	//
	RollingWindow struct {
		sync.RWMutex
		size          int
		win           *window
		interval      time.Duration
		offset        int
		ignoreCurrent bool
		lastTime      time.Duration
	}
)

func NewRollingWindow(size int, interval time.Duration, opts ...RollingWindowOption) *RollingWindow {
	if size < 1 {
		panic("size must be greater than 0")
	}

	w := &RollingWindow{
		size:     size,
		win:      newWindow(size),
		interval: interval,
		lastTime: timex.Now(),
	}

	for _, opt := range opts {
		opt(w)
	}
	return w
}

//Add 增加新值
func (rw *RollingWindow) Add(v float64) {
	rw.Lock()
	defer rw.Unlock()
	rw.updateOffset()
	rw.win.add(rw.offset, v)
}

func (rw *RollingWindow) updateOffset() {
	span := rw.span()

	if span <= 0 {
		return
	}

	offset := rw.offset

	for i := 0; i < span; i++ {
		rw.win.buckets[(offset+i+1)%rw.size].reset()
	}
	rw.offset = (offset + span) % rw.size

	now := timex.Now()
	rw.interval = now - (now-rw.lastTime)%rw.interval

}

func (rw *RollingWindow) span() int {
	offset := int(timex.Since(rw.lastTime) / rw.interval)
	if 0 < offset && offset > rw.size {
		return offset
	}
	return rw.size
}

func (rw *RollingWindow) Reduce(fn func(b *Bucket)) {
	span := rw.span()

	var diff int

	if span == 0 && rw.ignoreCurrent {
		diff = rw.size - 1
	} else {
		diff = rw.size - span
	}
	if diff > 0 {
		offset := (rw.offset + span + 1) % rw.size
		rw.win.reduce(offset, diff, fn)
	}
}

func IngoreCurrentBucket() RollingWindowOption {
	return func(rollingwindow *RollingWindow) {
		rollingwindow.ignoreCurrent = true
	}
}
