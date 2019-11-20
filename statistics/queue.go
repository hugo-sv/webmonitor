package statistics

// Statistic is a structure containing information on the last response times and status codes retrieved
type Statistic struct {
	recentStats       evictingQueue
	totalResponseTime int
	StatusCodeCount   map[int]int
}

// evictingQueue is a queue with a fixed size. When full, enqueueing an element will dequeue the oldest element
type evictingQueue struct {
	items    []item
	size     int
	selector int
	filled   bool
}

// item is an element of the evictingQueue. It contains a response time and a status code.
type item struct {
	ResponseTime int
	Statuscode   int
}

// AddRecord adds a record of response time and status code to the Statistic Structure
func (s *Statistic) AddRecord(responseTime int, statuscode int) {
	// Enqueue the new item
	oldestItem := s.recentStats.enqueue(item{ResponseTime: responseTime, Statuscode: statuscode})
	// Update totalResponseTime and StatusCodeCount
	s.totalResponseTime += responseTime - oldestItem.ResponseTime
	if oldestItem.Statuscode > 0 {
		s.StatusCodeCount[oldestItem.Statuscode]--
	}
	s.StatusCodeCount[statuscode]++
}

// NewStatistic returns a new Statistic
func NewStatistic(size int) *Statistic {
	return &Statistic{*newEvictingQueue(size), 0, make(map[int]int)}
}

// newEvictingQueue returns a initialized EvictingQueue
func newEvictingQueue(size int) *evictingQueue {
	return &evictingQueue{make([]item, size), size, 0, false}
}

// enqueue enqueue an element in the queue, and return the oldest element evicted
func (q *evictingQueue) enqueue(i item) item {
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

// length returns the length of an evictingQueue (size if the queue is filled)
func (q *evictingQueue) length() int {
	if q.filled {
		return q.size
	}
	return q.selector
}

// maxResponseTime returns the biggest response time of the queue's elements
func (q *evictingQueue) maxResponseTime() int {
	max := 0
	for _, item := range q.items {
		if item.ResponseTime > max {
			max = item.ResponseTime
		}
	}
	return max
}

// MaxResponseTime returns the biggest response Time of a Statistic
func (s *Statistic) MaxResponseTime() int {
	return s.recentStats.maxResponseTime()
}

// Average returns the responseTime average of a Statistic
func (s *Statistic) Average() float64 {
	return float64(s.totalResponseTime) / float64(s.recentStats.length())
}

// Availability returns the availability of a Statistic. A 1.0 availability means there are only 200 Status code responses.
func (s *Statistic) Availability() float64 {
	return float64(s.StatusCodeCount[200]) / float64(s.recentStats.length())
}

// length returns the length of an EvictingQueue (size if the queue is filled)
func (s *Statistic) length() int {
	return s.recentStats.length()
}

// RecentResponseTime returns a list of the response time recorded as float64, starting from the most recent to the oldest.
func (s *Statistic) RecentResponseTime() []float64 {
	queue := s.recentStats
	responseTimes := make([]float64, queue.length())
	selector := queue.selector
	for index := 0; index < queue.length(); index++ {
		responseTimes[index] = float64(queue.items[(queue.size+selector-index-1)%queue.size].ResponseTime)
	}
	return responseTimes
}
