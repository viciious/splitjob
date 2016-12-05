package splitjob

import (
	"sync"
)

func New(opts *Options) *Job {
	t := &Job{
		make([]*Split, opts.NumSplits),
		opts,
	}
	for i := range t.splits {
		t.splits[i] = newSplit(opts.Spawn, opts.ChanSize)
	}
	return t
}

func (t *Job) Do() error {
	var err error

	splits := t.splits
	numSplits := t.opts.NumSplits
	pull := t.opts.Pull

	done := make(chan error, numSplits)

	for _, s := range t.splits {
		s.signalStart()
	}

	var wg sync.WaitGroup
	for _, s := range splits {
		wg.Add(1)
		go func(s *Split) {
			err := s.wait()
			if err != nil {
				done <- err
			}
			wg.Done()
		}(s)
	}
	go func() {
		wg.Wait()
		close(done)
	}()

pullLoop:
	for {
		select {
		case err = <-done:
			break pullLoop
		default:
			msg, k, stop := pull()
			if stop {
				break pullLoop
			}
			s := splits[k%numSplits]
			s.queueMessage(msg)
		}
	}
	for _, s := range splits {
		s.signalStop()
	}

	if err == nil {
		err = <-done
	}
	return err
}
