package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	output := in

	for _, stage := range stages {
		input := chanWrap(output, done)
		output = stage(input)
	}

	return output
}

func chanWrap(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for val := range in {
			select {
			case out <- val:
			case <-done:
				return
			}
		}
	}()

	return out
}
