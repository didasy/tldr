package tldr

func maxHeapifyEdge(tosort []*Edge, position int) {
	size := len(tosort)
	maximum := position
	leftChild := 2*position + 1
	rightChild := leftChild + 1
	if leftChild < size && tosort[leftChild].weight > tosort[position].weight {
		maximum = leftChild
	}
	if rightChild < size && tosort[rightChild].weight > tosort[maximum].weight {
		maximum = rightChild
	}

	if position != maximum {
		tosort[position], tosort[maximum] = tosort[maximum], tosort[position]
		maxHeapifyEdge(tosort, maximum) //recursive
	}
}

func buildMaxHeapEdge(tosort []*Edge) {

	// from http://en.wikipedia.org/wiki/Heapsort
	// iParent = floor((i-1) / 2)

	for i := (len(tosort) - 1) / 2; i >= 0; i-- {
		maxHeapifyEdge(tosort, i)
	}
}

func HeapSortEdge(tosort []*Edge) {
	buildMaxHeapEdge(tosort)
	for i := len(tosort) - 1; i >= 1; i-- {
		tosort[i], tosort[0] = tosort[0], tosort[i]
		maxHeapifyEdge(tosort[:i-1], 0)
	}
}

func ReverseEdge(num []*Edge) {
	for i, j := 0, len(num)-1; i < j; i, j = i+1, j-1 {
		num[i], num[j] = num[j], num[i]
	}
}

func maxHeapifyRank(tosort []*Rank, position int) {
	size := len(tosort)
	maximum := position
	leftChild := 2*position + 1
	rightChild := leftChild + 1
	if leftChild < size && tosort[leftChild].score > tosort[position].score {
		maximum = leftChild
	}
	if rightChild < size && tosort[rightChild].score > tosort[maximum].score {
		maximum = rightChild
	}

	if position != maximum {
		tosort[position], tosort[maximum] = tosort[maximum], tosort[position]
		maxHeapifyRank(tosort, maximum) //recursive
	}
}

func buildMaxHeapRank(tosort []*Rank) {

	// from http://en.wikipedia.org/wiki/Heapsort
	// iParent = floor((i-1) / 2)

	for i := (len(tosort) - 1) / 2; i >= 0; i-- {
		maxHeapifyRank(tosort, i)
	}
}

func HeapSortRank(tosort []*Rank) {
	buildMaxHeapRank(tosort)
	for i := len(tosort) - 1; i >= 1; i-- {
		tosort[i], tosort[0] = tosort[0], tosort[i]
		maxHeapifyRank(tosort[:i-1], 0)
	}
}

func ReverseRank(num []*Rank) {
	for i, j := 0, len(num)-1; i < j; i, j = i+1, j-1 {
		num[i], num[j] = num[j], num[i]
	}
}
