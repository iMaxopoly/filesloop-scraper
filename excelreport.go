package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

func (fl *filesloop) generateExcelReport(timestamp string) {
	var (
		file        *xlsx.File
		sheet       *xlsx.Sheet
		row         *xlsx.Row
		cell        *xlsx.Cell
		headerStyle *xlsx.Style
		dataStyle   *xlsx.Style
		err         error
	)

	file = xlsx.NewFile()

	// create sheets
	sheet, err = file.AddSheet("links")
	if err != nil {
		fmt.Printf(err.Error())
	}
	sheet.SetColWidth(0, 0, 15) // index
	sheet.SetColWidth(1, 1, 20) // licensor
	sheet.SetColWidth(2, 2, 30) // siteLink
	sheet.SetColWidth(3, 3, 60) // pageTitle
	sheet.SetColWidth(4, 4, 20) // searchServerClass
	sheet.SetColWidth(5, 5, 30) // searchServerLink
	sheet.SetColWidth(6, 6, 30) // cyberlockerLink
	sheet.SetColWidth(7, 7, 20) // cyberlockerFiletype
	sheet.SetColWidth(8, 8, 20) // cyberlockerFilesize

	// setup style
	headerStyle = xlsx.NewStyle()
	headerStyle.Font.Bold = true
	headerStyle.Font.Name = "Calibri"
	headerStyle.Font.Size = 11
	headerStyle.Alignment.Horizontal = "center"
	headerStyle.Border.Top = "thin"
	headerStyle.Border.Bottom = "thin"
	headerStyle.Border.Right = "thin"
	headerStyle.Border.Left = "thin"

	dataStyle = xlsx.NewStyle()
	dataStyle.Font.Name = "Calibri"
	dataStyle.Font.Size = 11

	// header internal sheet
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "index"
	cell.SetStyle(headerStyle)

	cell = row.AddCell()
	cell.Value = "licensor"
	cell.SetStyle(headerStyle)

	cell = row.AddCell()
	cell.Value = "siteLink"
	cell.SetStyle(headerStyle)

	cell = row.AddCell()
	cell.Value = "pageTitle"
	cell.SetStyle(headerStyle)

	cell = row.AddCell()
	cell.Value = "searchServerClass"
	cell.SetStyle(headerStyle)

	cell = row.AddCell()
	cell.Value = "searchServerLink"
	cell.SetStyle(headerStyle)

	cell = row.AddCell()
	cell.Value = "cyberlockerLink"
	cell.SetStyle(headerStyle)

	cell = row.AddCell()
	cell.Value = "cyberlockerFiletype"
	cell.SetStyle(headerStyle)

	cell = row.AddCell()
	cell.Value = "cyberlockerFilesize"
	cell.SetStyle(headerStyle)

	data := readFileIntoList(DEBUGFILEPATH)

	for i, v := range data {
		// [0] = licensor; [1] = siteLink; [2] = pageTitle; [3] = searchServerClass; [4] = searchServerLink
		// [5] = cyberlockerLink; [6] = cyberlockerFiletype; [7] = cyberlockerFilesize

		values := strings.Split(v, "\t")
		debugLog(values)

		licensor := values[0]
		siteLink := values[1]
		pageTitle := values[2]
		searchServerClass := values[3]
		searchServerLink := values[4]
		cyberlockerLink := values[5]
		cyberlockerFiletype := values[6]
		cyberlockerFilesize := values[7]

		thisRow := sheet.AddRow()

		// index
		thisCell := thisRow.AddCell()
		thisCell.Value = strconv.Itoa(i)
		thisCell.SetStyle(headerStyle)

		// licensor
		thisCell = thisRow.AddCell()
		thisCell.Value = licensor
		thisCell.SetStyle(dataStyle)

		// siteLink
		thisCell = thisRow.AddCell()
		thisCell.Value = siteLink
		thisCell.SetStyle(dataStyle)

		// pageTitle
		thisCell = thisRow.AddCell()
		thisCell.Value = pageTitle
		thisCell.SetStyle(dataStyle)

		// searchServerClass
		thisCell = thisRow.AddCell()
		thisCell.Value = searchServerClass
		thisCell.SetStyle(dataStyle)

		// searchServerLink
		thisCell = thisRow.AddCell()
		thisCell.Value = searchServerLink
		thisCell.SetStyle(dataStyle)

		// cyberlockerLink
		thisCell = thisRow.AddCell()
		thisCell.Value = cyberlockerLink
		thisCell.SetStyle(dataStyle)

		// cyberlockerFiletype
		thisCell = thisRow.AddCell()
		thisCell.Value = cyberlockerFiletype
		thisCell.SetStyle(dataStyle)

		// cyberlockerFilesize
		thisCell = thisRow.AddCell()
		thisCell.Value = cyberlockerFilesize
		thisCell.SetStyle(dataStyle)
	}

	err = file.Save(fmt.Sprintf("%s/report_%s.xlsx", DEBUGFOLDER, timestamp))
	if err != nil {
		fmt.Printf(err.Error())
	}
}
