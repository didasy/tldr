package tldr

type ByWeight []*Edge

func (b ByWeight) Len() int {
	return len(b)
}

func (b ByWeight) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByWeight) Less(i, j int) bool {
	return b[i].weight < b[j].weight
}

type ByScore []*Rank

func (b ByScore) Len() int {
	return len(b)
}

func (b ByScore) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByScore) Less(i, j int) bool {
	return b[i].score < b[j].score
}

func ReverseEdge(num []*Edge) {
	for i, j := 0, len(num)-1; i < j; i, j = i+1, j-1 {
		num[i], num[j] = num[j], num[i]
	}
}

func ReverseRank(num []*Rank) {
	for i, j := 0, len(num)-1; i < j; i, j = i+1, j-1 {
		num[i], num[j] = num[j], num[i]
	}
}
