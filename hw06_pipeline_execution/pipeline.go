package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 || stages == nil {
		return in
	}

	out := in
	for _, s := range stages {
		if s != nil {
			out = s(doneWrapper(out, done))
		}
	}
	return out
}

func doneWrapper(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()

	return out
}
