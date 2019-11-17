package statistics

// Statistic is a structure containing informations on the last response times and status codes retrieved
type Statistic struct {
	recentStats       EvictingQueue
	totalResponseTime int
	statusCodeCount   map[int]int
}

// EvictingQueue is a queue with a fixed size. When full, enqueueing an element will dequeue the oldest element
type EvictingQueue struct {
	items    []Item
	size     int
	selector int
	filled   bool
}

// Item is an element of the EvictingQueue. It contains a response time and a status code.
type Item struct {
	responseTime int
	statuscode   int
}

// AddRecord adds a record of response time and status code to the Statistic Structure
func (s *Statistic) AddRecord(responseTime int, statuscode int) {
	// Enqueue the new Item
	oldestItem := s.recentStats.Enqueue(Item{responseTime: responseTime, statuscode: statuscode})
	// Update totalResponseTime and statusCodeCount
	s.totalResponseTime += responseTime - oldestItem.responseTime
	if oldestItem.statuscode > 0 {
		s.statusCodeCount[oldestItem.statuscode]--
	}
	s.statusCodeCount[statuscode]++
}

// NewStatistic returns a new Statistic
func NewStatistic(size int) *Statistic {
	return &Statistic{*NewEvictingQueue(size), 0, make(map[int]int)}
}

// NewEvictingQueue returns a initialized EvictingQueue
func NewEvictingQueue(size int) *EvictingQueue {
	return &EvictingQueue{make([]Item, size), size, 0, false}
}

// Enqueue enqueu an element in the queue, and return the oldest element evicted
func (q *EvictingQueue) Enqueue(i Item) Item {
	// Update filled status
	if !q.filled && q.selector == q.size-1 {
		q.filled = true
	}
	previousItem := q.items[q.selector]
	q.items[q.selector] = i
	// Increment selector
	q.selector = (q.selector + 1) % q.size
	return previousItem
}

// Lenght returns the lenght of an EvictingQueue (size if the queue is filled)
func (q *EvictingQueue) Lenght() int {
	if q.filled {
		return q.size
	}
	return q.selector
}

// MaxResponseTime returns the biggest response time of the queue's elements
func (q *EvictingQueue) MaxResponseTime() int {
	max := 0
	for _, item := range q.items {
		if item.responseTime > max {
			max = item.responseTime
		}
	}
	return max
}

// MaxResponseTime returns the biggest response Time of a Statistic
func (s *Statistic) MaxResponseTime() int {
	return s.recentStats.MaxResponseTime()
}

// Average returns the responseTime average of a Statistic
func (s *Statistic) Average() float64 {
	return float64(s.totalResponseTime) / float64(s.recentStats.Lenght())
}

// StatusCodeCount returns the status code count map of a Statistic
func (s *Statistic) StatusCodeCount() map[int]int {
	return s.statusCodeCount
}

// Availability returns the availability of a Statistic. A 1.0 availability means there are only 200 Status code responses.
func (s *Statistic) Availability() float64 {
	return float64(s.statusCodeCount[200]) / float64(s.recentStats.Lenght())
}
