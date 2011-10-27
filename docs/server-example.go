package main

type request struct {
	op     func(int, int) int
	a, b   int
	replyc chan int
}

func server(service chan *request, quit chan bool) {
	for {
		select {
		case req := <-service:
			go func() { // don't wait for computation
				req.replyc <- req.op(req.a, req.b)
			}()
		case <-quit:
			return
		}
	}
}

func startServer() (service chan *request, quit chan bool) {
	service = make(chan *request)
	quit = make(chan bool)
	go server(service, quit)
	return service, quit
}

func main() {
	server, quit := startServer()
	const N = 100
	var reqs [N]request
	mul := func(x, y int) int { return x * y }
	for i := 0; i < N; i++ {
		req := &reqs[i]
		req.op = mul
		req.a = i
		req.b = i + N
		req.replyc = make(chan int)
		server <- req
	}
	for i := N - 1; i >= 0; i-- { // doesn't matter what order
		if <-reqs[i].replyc != i*(i+N) {
			println("fail at", i)
		} else {
			print(".")
		}
	}
	quit <- true
}
