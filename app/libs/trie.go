package libs

import (
	"errors"
)

const (
	FIND     = 1
	NOT_FIND = 0
)

var FundTrie *TrieNode

// Trie
type TrieNode struct {
	key      string
	data     *TrieData            // 数据
	dataList map[string]*TrieNode // 子节点数据
}

type TrieData struct {
	FundCode         string `json:"fund_code"`         // 基金 code
	FundAbbreviation string `json:"fund_abbreviation"` // 基金 缩写
}

// 返回本节点的 数据
func (node *TrieNode) Get() *TrieData {
	return node.data
}

func (node *TrieNode) GetKey() string {
	return node.key
}

func (node *TrieNode) InitDataList() {
	node.dataList = make(map[string]*TrieNode)
}

// 返回所有满足 到此节点为止的所有 子节点 共同组成的单词(数据)
// 此方法会默认遍历所有 数据集, 效率 会 很低(类似mysql的 模糊查询)
func (node *TrieNode) GetFullList(getData []string, isFind int) (list []*TrieData) {
	var (
		data *TrieNode
		str  []string
	)

	// 如果是最后一个,根据前置是否 已经发现对应元素进行 数据返回
	if len(node.dataList) == 0 {
		if isFind == FIND {
			list = append(list, node.Get())
		}

		if len(getData) > 0 && node.GetKey() == getData[len(getData)-1] {
			list = append(list, node.Get())
		}
		return
	}

	if isFind == FIND && node.data != nil {
		list = append(list, node.Get())
	}

	if len(getData) > 0 {
		data = node.dataList[getData[0]]
	}

	if data != nil && isFind == NOT_FIND && len(getData) == 1 {
		isFind = FIND
	}

	if isFind == NOT_FIND && data != nil {
		str = getData[1:]
	} else if isFind == NOT_FIND && data == nil {
		str = getData
	} else if isFind == FIND && len(getData) > 0 {
		str = getData[1:]
	}

	// 未发现, 也就是 isFind 为0的时候
	for _, v := range node.dataList {
		if len(getData) == 1 && v.GetKey() != getData[0] && isFind == FIND {
			for _, v2 := range v.GetFullList(getData, NOT_FIND) {
				list = append(list, v2)
			}
			continue
		}

		for _, v2 := range v.GetFullList(str, isFind) {
			list = append(list, v2)
		}
	}

	return
}

func (node *TrieNode) GetList(getData []string, isFind int) (list []*TrieData) {
	var (
		data *TrieNode
		str  []string
	)

	// 如果是最后一个,根据前置是否 已经发现对应元素进行 数据返回
	if len(node.dataList) == 0 {
		if isFind == FIND {
			list = append(list, node.Get())
		}

		if len(getData) > 0 && node.GetKey() == getData[len(getData)-1] {
			list = append(list, node.Get())
		}
		return
	}

	if isFind == FIND && node.data != nil {
		list = append(list, node.Get())
	}

	if len(getData) > 0 {
		data = node.dataList[getData[0]]
	}

	if data != nil && isFind == NOT_FIND && len(getData) == 1 {
		isFind = FIND
	}

	if len(getData) > 0 {
		str = getData[1:]
	}

	// 未发现, 也就是 isFind 为0的时候
	for _, v := range node.dataList {
		if len(getData) >= 1 && v.GetKey() != getData[0] {
			continue
		}

		for _, v2 := range v.GetList(str, isFind) {
			list = append(list, v2)
		}
	}

	return
}

// 插入数据
func (node *TrieNode) Insert(key []string, insertData *TrieData, setupNum int) (err error) {
	var (
		data     *TrieNode
		trieData *TrieData
	)

	if len(key) == 0 {
		if setupNum == 1 {
			err = errors.New("请传入非空数据")
			return
		}
		//fmt.Println("数据插入成功")
		return
	}

	data = node.dataList[key[0]]

	if data != nil {
		// 数据已存在, 进行递归插入
		if len(key) == 1 {
			// 插入数据
			if data.Get() == nil {
				data.data = insertData
			}
		}

		err = data.Insert(key[1:], insertData, setupNum+1)

	} else {
		if len(key) == 1 {
			trieData = insertData
		} else {
			trieData = nil
		}
		// 插入数据
		node.dataList[key[0]] = &TrieNode{
			key:      key[0],
			data:     trieData,
			dataList: make(map[string]*TrieNode),
		}

		err = node.dataList[key[0]].Insert(key[1:], insertData, setupNum+1)
	}

	return
}
