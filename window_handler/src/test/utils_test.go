package test

import (
	"testing"
)

func TestCreateFile(t *testing.T) {
	if !createFile("E:/source/tree_ttt", 1024*512, true) {
		t.Error("error")
	}
}

func TestCreateFileTree(t *testing.T) {
	count := 0
	createFileTree("E:/source/tree_ttt", 5, 5, 8, false, true, &count)
}

func TestRandomContent(t *testing.T) {
	t.Logf("Random result : %s", randomPalindrome(1024))
}
