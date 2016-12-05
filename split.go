package splitjob

func newSplit(spawn SpawnFn, chanSize int) *Split {
	Split := &Split{
		think: spawn(),
		in:    make(chan interface{}, chanSize),
		out:   make(chan error),
	}
	return Split
}

func (split *Split) queueMessage(obj interface{}) {
	split.in <- obj
}

func (split *Split) signalStart() {
	go func() {
		var err error
		var in <-chan interface{} = split.in
		var out chan<- error = split.out
	readLoop:
		for {
			select {
			case obj := <-in:
				if obj == nil {
					break readLoop
				}
				if err == nil {
					if err = split.think(obj); err != nil {
						out <- err
					}
				}
			}
		}
		close(out)
	}()
}

func (split *Split) signalStop() {
	close(split.in)
}

func (split *Split) wait() error {
	err := <-split.out
	return err
}
