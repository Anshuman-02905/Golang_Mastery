package main

func orDone(done <-chan interface{}, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		select {
		case <-done:
			return
		case v, ok := <-c:
			if ok == false {
				return
			}
			select {
			case valStream <- v:
			case <-done:
			}
		}
	}()
	return valStream
}

func main() {

}
