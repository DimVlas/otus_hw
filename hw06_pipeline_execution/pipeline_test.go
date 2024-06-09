package hw06pipelineexecution

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})

	t.Run("empty case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		stages := []Stage{}

		result := make([]int, 0, 5)
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(int))
		}

		require.Len(t, result, 5)
		require.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("empty data", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 5)
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}

		require.Len(t, result, 0)
		require.Equal(t, []string{}, result)
	})
}

func TestChanWrap(t *testing.T) {
	t.Run("To end", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
				time.Sleep(time.Millisecond * 10)
			}
			close(in)
		}()

		res := make([]int, 0, 5)
		for i := range chanWrap(in, done) {
			res = append(res, i.(int))
		}

		require.Len(t, res, 5)
		require.Equal(t, []int{1, 2, 3, 4, 5}, res)
	})

	t.Run("Close Done", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for i, v := range data {
				if i == 2 {
					close(done)
					return
				}
				in <- v
				time.Sleep(time.Millisecond * 10)
			}
			close(in)
		}()

		res := make([]int, 0, 5)
		for i := range chanWrap(in, done) {
			res = append(res, i.(int))
		}

		require.Len(t, res, 2)
		require.Equal(t, []int{1, 2}, res)
	})

	t.Run("Close In", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for i, v := range data {
				if i == 2 {
					close(in)
					return
				}
				in <- v
				time.Sleep(time.Millisecond * 10)
			}
			close(in)
		}()

		res := make([]int, 0, 5)
		for i := range chanWrap(in, done) {
			res = append(res, i.(int))
		}

		require.Len(t, res, 2)
		require.Equal(t, []int{1, 2}, res)
	})
}
