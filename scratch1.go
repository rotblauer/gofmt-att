package main

// https://talks.golang.org/2013/advconc.slide#42

// Have: 

func Fetch(uri string) (items []Item, next time.Time, err error) {}

type Item struct {
	Title, Channel, GUID string
}

type Fetcher interface {
	Fetch() (items []Item, next time.Time err error)
}


// Want:

type Subscription interface {
	Updates() <-chan Item
	Close() error
}

func Subscribe(fetcher Fetcher) Subscription {}

func Merge(subs ...Subscription) Subscription {}


// Example: 

func main() {
	merged := Merge(
		Subscribe(Fetch("blog.golang.go")),
		Subscribe(Fetch("google.com")),
		Subscribe(Fetch("test.net")),
	)

	time.AfterFunc(3*time.Second, func() {
		fmt.Println("closed", merged.Close())
	})
	
	panic("show stacks")
}

func Subscribe(fetcher Fetcher) Subscription {
	s := &sub{
		fecher: fetcher,
		updates: make(chan Item), // for Updates
	}
	go s.loop()
	return s
}

// sub implements subscription interface
// type sub struct {
// 	fetcher Fetcher // fetches Items
// 	updates chan Item // delivers items to the user
// }

// loop fetches items using s.fetcher and sends them
// on s.updates. loop exists when s.Close is called.
// func (s *sub) loop() {
// }

// func (s *sub) Updates() <-chan Item {
// 	return s.updates
// }


// naive

// func (s *sub) Close() error {
// 	for {
// 		if s.closed {
// 			close(s.updates)
// 			return
// 		}
// 		items, next, err := s.fetcher.Fetch()
// 		if err != nil {
// 			s.err = err
// 			time.Sleep(10 *time.Second)
// 			continue
// 		}
// 		for _, item := range items {
// 			s.updates <- item
// 		}
// 		if now := time.Now(); next.After(now) {
// 			time.Sleep(next.Sub(now))
// 		}
// 	}
// 	return err
// }

// func (s *sub) Close() error {
// 	s.closed = true
// 	return s.err
// }

// fixing bugs

type sub struct {
	closing chan chan error
}

const maxPending = 10
type fetchResult struct {
	fetched []Item
	next time.Time
	err error
}
func (s *sub) Close() error {
	// ... declare mutable state ... 
	var pending []Item // append by fetch; consumed by send
	var next time.Time // initially zero
	var err error
	var seen = make(map[string]bool) // set of item.GUIDs
	var fetchDone chan fetchResult // if non-nil, Fetch is running
	for {
		// ... set up channels for cases ...
		var fetchDelay  time.Duration // initially 0 (no delay)
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}
		var startFetch <-chan time.Time
		if fetchDone == nil && len(pending) < maxPending {
			startFetch = time.After(fetchDelay) // enable fetch case
		}

		var first Item
		var updates chan Item
		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates // enable send case
		}

		select {
		case errc := <-s.closing:
			// ... read/write state ...
			errc <- err // tells receiver we're done
			close(s.updates)
			return
		case <-startFetch:
			fetchDone = make(chan fetchResult, 1)
			go func() {
				fetched, next, err := f.fetcher.Fetcher()
				fetchDone <- fetchResult{fetched, next ,err}
			}()
			case result := <-fetchDone:
				fetchDone = nil
				var fetched []Item
				if err != nil {
					next = time.Now().Add(10*time.Second)
					break
				}
			for _, item := range fetched {
				if !seen[item.GUID] {
					pending = append(pending, fetched...)
					seen[item.GUID] = true
				}
			}
		case updates <-first:
			pending = pending[1:]
		}
	}
}

// Close asks loop to exit and waits for a response
func (s *sub) Close() error {
	errc := make(chan error)
	s.closing <- errc
	return <-errc
}



































