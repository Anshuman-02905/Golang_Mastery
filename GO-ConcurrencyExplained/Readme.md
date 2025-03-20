Mulptiple Go without waitgroup
Your issue occurs because the main function exits before the goroutines get a chance to run. Since goroutines run asynchronously, the program does not wait for them to complete. Once main exits, all spawned goroutines are killed.

With waitgroup

in main or pass i function
var wg sync.WaitGroup
	wg.Add(3)

    	wg.Wait()


Mutex avoid Race Conditio