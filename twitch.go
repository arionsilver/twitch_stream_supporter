package main

func startTwitchHelper(auth string, q chan bool) (c chan bool) {
	c = make(chan bool)
	go executeTwitchHelper(auth, c, q)

	return
}

func executeTwitchHelper(auth string, c chan bool, q chan bool) {
	defer func() { c <- true }()
	<-q // wait on quit
}
