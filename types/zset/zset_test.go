package zset

import "testing"

func TestSkipList(t *testing.T) {
	skipList := makeSkiplist()
	skipList.insert("a1", 200)
	skipList.insert("a2", 300)
	skipList.insert("a3", 400)
	skipList.insert("a4", 120)

	t.Logf("%v", skipList.getRank("a1", 200))
}
