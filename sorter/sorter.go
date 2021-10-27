package sorter

import (
	"io/fs"
	"sort"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func NatSort(files []fs.FileInfo) []*fileMetadata {
	_files := make([]*fileMetadata, len(files))
	for i, file := range files {
		_files[i] = toFileMetadata(file)
	}

	sort.SliceStable(_files, func(l, r int) bool {
		lf := _files[l]
		rf := _files[r]

		lNameLen := len(lf.StrTokens)
		rNameLen := len(rf.StrTokens)
		nameCmpLen := min(lNameLen, rNameLen)
		for i := 0; i < nameCmpLen; i++ {
			if lf.TokenIsNums[i] && rf.TokenIsNums[i] {
				if lf.NumTokens[i] < rf.NumTokens[i] {
					return true
				} else if lf.NumTokens[i] > rf.NumTokens[i] {
					return false
				}
			} else {
				if lf.StrTokens[i] < rf.StrTokens[i] {
					return true
				} else if lf.StrTokens[i] > rf.StrTokens[i] {
					return false
				}
			}
		}

		if lNameLen < rNameLen {
			return true
		}
		if rNameLen < lNameLen {
			return false
		}

		return lf.Ext < rf.Ext
	})

	return _files
}
