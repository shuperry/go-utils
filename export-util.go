package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type ExcelInfo struct {
	BaseXlsxName *string
	Sheets       []*SheetInfo
}

type SheetInfo struct {
	Name         *string
	Title        *string
	Data         []map[string]interface{}
	DataFields   []*DataFieldInfo
	NeedNoColumn *bool
	NeedTitle    *bool
}

type DataFieldInfo struct {
	FieldName *string
	Title     *string
	Type      *string
	Width     *float64
	Render    *func(data interface{}, fieldData interface{}) interface{}
}

func GenerateExcel(ei *ExcelInfo) string {
	xlsx := excelize.NewFile()

	// 表格名称单元格格式.
	style, _ := xlsx.NewStyle(`{
		"font": {"bold": true, "size": 32},
		"alignment": {"horizontal": "center"},
		"border": [
			{"type": "left", "color": "000000", "style": 1},
			{"type": "top", "color": "000000", "style": 1},
			{"type": "bottom", "color": "000000", "style": 1},
			{"type": "right", "color": "000000", "style": 1}
		]
	}`)

	//生成标题行单元格颜色格式.
	titleRowStyle, _ := xlsx.NewStyle(`{
		"font": {"bold": true, "size": 14},
		"alignment": {"horizontal": "center"},
		"fill": {"type": "pattern", "color": ["#EEC900"], "pattern": 9},
		"border": [
			{"type": "left", "color": "000000", "style": 1},
			{"type": "top", "color": "000000", "style": 1},
			{"type": "bottom", "color": "000000", "style": 1},
			{"type": "right", "color": "000000", "style": 1}
		]
	}`)

	// 生成数据单元格边框格式.
	dataCellStyle, _ := xlsx.NewStyle(`{
		"border": [
			{"type": "left", "color": "000000", "style": 1},
			{"type": "top", "color": "000000", "style": 1},
			{"type": "bottom", "color": "000000", "style": 1},
			{"type": "right", "color": "000000", "style": 1}
		]
	}`)

	characters := []string{
		"A", "B", "C", "D", "E",
		"F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O",
		"P", "Q", "R", "S", "T",
		"U", "V", "W", "X", "Y",
		"Z",
	}

	var (
		sheetName          string
		fieldName          string
		titleFieldPosition string
		fieldPosition      string
	)

	for i, sheet := range ei.Sheets {
		sheetName = *sheet.Name

		if i == 0 {
			xlsx.SetSheetName("Sheet1", sheetName)
		} else {
			xlsx.NewSheet(sheetName)
		}

		startColIndex := 0
		if sheet.NeedNoColumn != nil && *sheet.NeedNoColumn {
			startColIndex = 1
		}

		startRowIndex := -1
		if sheet.NeedTitle != nil && *sheet.NeedTitle {
			// 设置表格名称.
			xlsx.SetCellValue(sheetName, "A1", *sheet.Title)

			// 合并表格名称单元格.
			sheetTitlePosition := fmt.Sprintf("%s%s", characters[startColIndex+len(sheet.DataFields)-1], "1")

			// 表格名称单元格样式.
			xlsx.MergeCell(sheetName, "A1", sheetTitlePosition)
			xlsx.SetCellStyle(sheetName, "A1", sheetTitlePosition, style)

			startRowIndex = 0
		}

		if sheet.NeedNoColumn != nil && *sheet.NeedNoColumn {
			noFieldPosition := fmt.Sprintf("%s%s", "A", strconv.Itoa(2+startRowIndex))
			// 添加序号列.
			xlsx.SetCellValue(sheetName, noFieldPosition, "序号")
			// 设置序号列宽度.
			xlsx.SetColWidth(sheetName, "A", "A", 8)

			// 标题行中序号单元格样式.
			xlsx.SetCellStyle(sheetName, noFieldPosition, noFieldPosition, titleRowStyle)
		}

		for i1, field := range sheet.DataFields {
			if field.FieldName != nil {
				fieldName = *field.FieldName
			}

			// 标题行数据及样式.
			titleFieldPosition = fmt.Sprintf("%s%s", characters[i1+startColIndex], strconv.Itoa(2+startRowIndex))

			xlsx.SetCellValue(sheetName, titleFieldPosition, *field.Title)
			xlsx.SetCellStyle(sheetName, titleFieldPosition, titleFieldPosition, titleRowStyle)

			// 设置各数据列宽度.
			xlsx.SetColWidth(sheetName, characters[i1+startColIndex], characters[i1+startColIndex], *field.Width)

			for i2, data := range sheet.Data {
				if sheet.NeedNoColumn != nil && *sheet.NeedNoColumn {
					// 序号单元格一列数据填充及样式敲定.
					fieldPosition = fmt.Sprintf("%s%s", "A", strconv.Itoa(i2+startColIndex+2+startRowIndex))
					xlsx.SetCellValue(sheetName, fieldPosition, i2+1)
					xlsx.SetCellStyle(sheetName, fieldPosition, fieldPosition, dataCellStyle)
				}

				if sheet.NeedNoColumn != nil && *sheet.NeedNoColumn {
					fieldPosition = fmt.Sprintf("%s%s", characters[i1+startColIndex], strconv.Itoa(startColIndex+i2+2+startRowIndex))
				} else {
					fieldPosition = fmt.Sprintf("%s%s", characters[i1+startColIndex], strconv.Itoa(startColIndex+i2+2+1+startRowIndex))
				}

				// 如果有配置 Render 函数, 则单元格填入 Render 函数返回的值
				if field.FieldName != nil && field.Render != nil && data[fieldName] != nil {
					render := *field.Render
					render(data, data[fieldName])
				} else {
					switch *field.Type {
					case "string":
						xlsx.SetCellValue(sheetName, fieldPosition, data[fieldName].(string))
					case "float64":
						xlsx.SetCellValue(sheetName, fieldPosition, data[fieldName].(float64))
					default:
						xlsx.SetCellValue(sheetName, fieldPosition, data[fieldName])
					}
				}

				// 设置所有数据单元格边框样式.
				xlsx.SetCellStyle(sheetName, fieldPosition, fieldPosition, dataCellStyle)
			}
		}
	}

	xlsxName := strings.Join([]string{*ei.BaseXlsxName, time.Now().Format("20060102150405")}, "_")
	savedExcelFilePath := strings.Join([]string{"./", xlsxName, ".xlsx"}, "")

	err := xlsx.SaveAs(savedExcelFilePath)
	if err != nil {
		panic(err)
	}

	return savedExcelFilePath
}
