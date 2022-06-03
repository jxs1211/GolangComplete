package goleak1

func leak() {
	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
	}()
}
