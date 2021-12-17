package list

import "testing"

func TestLinkedList(t *testing.T) {
	linkedList := Make()
	linkedList.Add(20)
	linkedList.Add(40)
	linkedList.Add(60)
	linkedList.Add(80)
	linkedList.Add(100)

	l := linkedList.Len()
	t.Logf("Linked List length: %d", l)
	t.Logf("%v", linkedList.Range(0, 4))
}
