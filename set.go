package main

func newSet() set {
	return make(set)
}

type set map[interface{}]struct{}

func (set *set) Add(i interface{}) bool {
	_, found := (*set)[i]
	(*set)[i] = struct{}{}
	return !found
}

func (set *set) Contains(i ...interface{}) bool {
	for _, val := range i {
		if _, ok := (*set)[val]; !ok {
			return false
		}
	}

	return true
}

func (set *set) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for elem := range *set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *set) Size() int {
	return len(*set)
}

func (set *set) ToSliceString() []string {
	keys := make([]string, 0, set.Size())
	for elem := range *set {
		keys = append(keys, elem.(string))
	}

	return keys
}
