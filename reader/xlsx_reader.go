package reader

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"path/filepath"
	"sort"
	"xCelFlow/entities"
	"xCelFlow/util"
)

type XLSXReader struct {
	*Reader
}

type headPair struct {
	headMap1       map[entities.TupleT]int
	headList1      []keyIndex
	fieldNameMap1  map[entities.TupleT]int
	headMap2       map[entities.TupleT]int
	headList2      []keyIndex
	fieldNameList2 []keyIndex
}

type keyIndex struct {
	key   entities.TupleT
	index int
}

func init() {
	Register("xlsx", newXLSXReader)
}

func newXLSXReader(reader *Reader) IReader {
	return &XLSXReader{reader}
}

func (r *XLSXReader) Read() ([][]string, error) {
	file, err := excelize.OpenFile(r.Path)
	if err != nil {
		return nil, errorTableReadFailed(r.Path, err)
	}

	defer func() { _ = file.Close() }()

	var records [][]string
	filename := filepath.Base(r.Path)
	tableName, err := util.SubTableName(filename)
	if err != nil {
		return nil, errorTableReadFailed(r.Path, err)
	}

	for _, sheetName := range file.GetSheetList() {
		if name, _ := util.SubTableName(sheetName); tableName == name {
			if tempRecords, err := file.GetRows(sheetName); err != nil {
				return nil, errorTableReadFailed(r.Path, err)
			} else {
				r.alignHead(&tempRecords)
				if len(records) == 0 {
					records = tempRecords
				} else if p, err := r.createHeadMap(&records, &tempRecords, sheetName); err != nil {
					return nil, errorTableReadFailed(r.Path, err)
				} else {
					r.mergeRecords(p, &records, &tempRecords)
				}
			}
		}
	}

	if records == nil || len(records) == 0 {
		return nil, errorTableNotSheet(r.Path)
	}

	return records, nil
}

func (r *XLSXReader) createHeadMap(records *[][]string, tempRecords *[][]string, sheetName string) (headPair, error) {
	if headMap1, fieldNameMap1, err := r.initHeadMap(records, sheetName); err != nil {
		return headPair{}, err
	} else if headMap2, fieldNameMap2, err := r.initHeadMap(tempRecords, sheetName); err != nil {
		return headPair{}, err
	} else {
		p := headPair{
			headMap1, initHeadList(headMap1), fieldNameMap1,
			headMap2, initHeadList(headMap2), initHeadList(fieldNameMap2),
		}
		return p, nil
	}
}

// 对齐表头
func (r *XLSXReader) alignHead(records *[][]string) {
	headRows := (*records)[:r.BodyStartIndex] // 前五行
	// 获取最大列数
	maxColNum := 0
	for _, row := range headRows {
		maxColNum = max(maxColNum, len(row))
	}

	// 补齐每行数据
	for i := 0; i < len(headRows); i++ {
		row := &headRows[i]
		if len(*row) < maxColNum {
			*row = append(*row, make([]string, maxColNum-len(*row))...)
		}
	}
}

// initHeadMap 获取表头列信息
func (r *XLSXReader) initHeadMap(records *[][]string, sheetName string) (map[entities.TupleT]int, map[entities.TupleT]int, error) {
	headMap := make(map[entities.TupleT]int)
	fieldNameMap := make(map[entities.TupleT]int)
	headRows := (*records)[:r.BodyStartIndex] // 前五行
	fieldNameRorNum := len(r.FieldNameIndexList)

	// 获取表头列信息
	for colIndex := 0; colIndex < len(headRows[0]); colIndex++ {
		fullKey := entities.TupleT{}
		isEmptyCol := true
		for rowIndex := 0; rowIndex < r.BodyStartIndex; rowIndex++ {
			cell := headRows[rowIndex][colIndex]
			if cell != "" {
				isEmptyCol = false
			}
			fullKey[rowIndex] = cell
		}

		// 字段名+类型
		fnKey := entities.TupleT{}
		for index, rowIndex := range r.FieldNameIndexList {
			if cell := headRows[rowIndex][colIndex]; cell != "" {
				fnKey[index] = cell
			}
		}
		fnKey[fieldNameRorNum] = headRows[r.FieldTypeIndex][colIndex]

		if isEmptyCol {
			fmt.Printf("页签：%s 单元格：%s 存在空表头，多页签数据合并可能无法正确进行，建议至少添加一个表头数据或删除此列\n", sheetName, util.ToCell(0, colIndex))
		} else if _, ok := headMap[fullKey]; ok {
			return nil, nil, fmt.Errorf("页签：%s 单元格：%s 存在重复表头\n", sheetName, util.ToCell(0, colIndex))
		} else if _, ok := fieldNameMap[fnKey]; ok {
			return nil, nil, fmt.Errorf("页签：%s 单元格：%s 存在重复表头\n", sheetName, util.ToCell(0, colIndex))
		} else {
			headMap[fullKey] = colIndex
			fieldNameMap[fnKey] = colIndex
		}
	}

	return headMap, fieldNameMap, nil
}

// initHeadList 获取表头列信息
func initHeadList(headMap map[entities.TupleT]int) []keyIndex {
	headList := make([]keyIndex, 0, len(headMap))
	for k, v := range headMap {
		headList = append(headList, keyIndex{k, v})
	}
	sort.Slice(headList, func(i, j int) bool {
		return headList[i].index < headList[j].index
	})
	return headList
}

// mergeRecords 合并记录
func (r *XLSXReader) mergeRecords(p headPair, records *[][]string, newRecords *[][]string) {
	// 扩充records
	maxColNum := len(p.headList1)
	maxRowNum := len(*records)
	extendNum := len(*newRecords) - r.BodyStartIndex
	for i := 0; i < extendNum; i++ {
		*records = append(*records, make([]string, maxColNum))
	}

	headRows := (*records)[:r.BodyStartIndex]
	bodyRows := (*records)[maxRowNum:]
	for index, e := range p.headList2 {
		if colIndex, ok := p.headMap1[e.key]; ok {
			// 总sheet表中有新sheet表中的字段
			for rowIndex, row := range (*newRecords)[r.BodyStartIndex:] {
				if len(row) > index {
					bodyRows[rowIndex][colIndex] = row[index]
				}
			}
		} else if colIndex = fuzzyIndex(p.fieldNameList2[index], p.fieldNameMap1); colIndex != -1 {
			// 总sheet表中有新sheet表中的字段
			for rowIndex, row := range (*newRecords)[r.BodyStartIndex:] {
				if len(row) > index {
					bodyRows[rowIndex][colIndex] = row[index]
				}
			}
		} else {
			// 总sheet表中没有新sheet表中的字段
			// 追加表头
			for rowIndex := 0; rowIndex < len(headRows); rowIndex++ {
				headRows[rowIndex] = append(headRows[rowIndex], (*newRecords)[rowIndex][e.index])
			}
			// 追加数据
			for rowIndex, row := range (*newRecords)[r.BodyStartIndex:] {
				if len(row) > index {
					bodyRows[rowIndex] = append(bodyRows[rowIndex], row[index])
				}
			}
		}
	}
}

func fuzzyIndex(e keyIndex, fieldNameMap map[entities.TupleT]int) int {
	if index, ok := fieldNameMap[e.key]; ok {
		return index
	}
	return -1
}
