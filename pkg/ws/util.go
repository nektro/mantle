package ws

import (
	"container/list"
)

func listHas(l *list.List, n interface{}) bool {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == n {
			return true
		}
	}
	return false
}

func listRemove(l *list.List, n interface{}) {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == n {
			l.Remove(e)
			break
		}
	}
}

func listToArray(l *list.List) []string {
	res := []string{}
	for e := l.Front(); e != nil; e = e.Next() {
		res = append(res, e.Value.(string))
	}
	return res
}
