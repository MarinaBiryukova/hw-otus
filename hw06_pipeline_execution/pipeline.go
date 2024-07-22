package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		res := make(Bi)
		close(res)
		return res
	}

	prevOut := in

	for _, stage := range stages {
		curIn := make(Bi)

		go func(next Bi, prev Out) {
			defer func() {
				close(next)
				for range prev {
					continue
				}
			}()

			for {
				select {
				case <-done:
					return
				case elem, ok := <-prev:
					if !ok {
						return
					}
					next <- elem
				}
			}
		}(curIn, prevOut)

		prevOut = stage(curIn)
	}

	return prevOut
}
