package splitjob

type ThinkFn func(obj interface{}) error
type SpawnFn func() ThinkFn
type PullFn func() (interface{}, uint32, bool)

type Options struct {
	Spawn     SpawnFn
	Pull      PullFn
	ChanSize  int
	NumSplits uint32
}

type Split struct {
	think ThinkFn
	in    chan interface{}
	out   chan error
}

type Job struct {
	splits []*Split
	opts   *Options
}
