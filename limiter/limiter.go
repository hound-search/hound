package limiter

type Limiter chan bool

func New(n int) Limiter {
	return limiter(make(chan bool, n))
}

func (l Limiter) Acquire() {
	l <- true
}

func (l Limiter) Release() {
	<-l
}
