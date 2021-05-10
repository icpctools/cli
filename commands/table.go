package commands

import (
	"fmt"
)

type rowStr []string

const ALIGN_LEFT = 0
const ALIGN_RIGHT = 1

type Table struct {
	Header rowStr
	Rows   []rowStr
	Align  []int
}

func (table *Table) appendRow(row []string) {
	table.Rows = append(table.Rows, row)
}

func (table Table) print() error {
	// determine the amount of padding needed
	var numCol = len(table.Header)
	var maxLength []int = make([]int, numCol)
	var format []string = make([]string, numCol)

	// find max header width
	for i, s := range table.Header {
		if maxLength[i] < len(s) {
			maxLength[i] = len(s)
		}
	}

	// find max cell width
	for _, r := range table.Rows {
		for i, s := range r {
			if maxLength[i] < len(s) {
				maxLength[i] = len(s)
			}
		}
	}

	// create format for each column, respecting width and alignment
	for i := range table.Header {
		if table.Align[i] == ALIGN_LEFT {
			format[i] = fmt.Sprintf(" %%-%vv ", maxLength[i])
		} else {
			format[i] = fmt.Sprintf(" %%%vv ", maxLength[i])
		}
	}

	// output header bold and underlined
	fmt.Printf("  \033[1;4m")
	for i, k := range table.Header {
		fmt.Printf(format[i], k)
	}
	fmt.Printf("\033[0m\n")

	// output each cell
	for _, r := range table.Rows {
		fmt.Printf("  ")
		for i, s := range r {
			fmt.Printf(format[i], s)
		}
		fmt.Printf("\n")
	}

	return nil
}
