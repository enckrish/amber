package main

import (
	"encoding/json"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
)

const columnsUTable = 5
const messageIdCol = 2
const severityCol = 3

var app = tview.NewApplication()
var (
	overviewTable *tview.Table
	logsBox       *tview.TextView
	resultsBox    *tview.TextView
)

func runUI(p *Pipeline) {
	app.SetInputCapture(handleInput)

	view := buildLayout(p)
	err := app.SetRoot(view, true).EnableMouse(true).Run()
	if err != nil {
		panic(err)
	}
}

func buildLayout(p *Pipeline) tview.Primitive {
	overviewTable = getOverviewTable(p)
	overviewTable.SetBackgroundColor(tcell.ColorDefault)
	detailView := getDetailedFlex()

	baseFlex := tview.NewFlex()
	baseFlex.SetDirection(tview.FlexRow).
		AddItem(overviewTable, 0, 1, false).
		AddItem(detailView, 0, 1, false)

	return baseFlex
}

func getDetailedFlex() tview.Primitive {
	flex := tview.NewFlex()

	logsBox = tview.NewTextView()
	logsBox.SetBackgroundColor(tcell.ColorDefault)
	logsBox.SetBorder(true).SetTitle("Logs")
	logsBox.SetScrollable(true)

	resultsBox = tview.NewTextView()
	resultsBox.SetBackgroundColor(tcell.ColorDefault)
	resultsBox.SetBorder(true).SetTitle("Results")
	resultsBox.SetScrollable(true)

	flex.SetDirection(tview.FlexColumn).
		AddItem(logsBox, 0, 1, false).
		AddItem(resultsBox, 0, 1, false)

	return flex
}

func getOverviewTable(p *Pipeline) *tview.Table {
	table := tview.NewTable().SetFixed(1, columnsUTable).SetSelectable(true, false)
	table.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			table.SetSelectable(false, false)
		}
	})
	table.SetFocusFunc(func() {
		updateLogsBox(p)
		updateResultsBox(p)
	})

	table.SetBlurFunc(
		func() {
			//logsBox.Clear()
			//resultsBox.Clear()
		})

	table.SetBorder(true)
	overviewTable = table
	updateOverviewTableData(p)
	return table
}

func updateOverviewTableData(p *Pipeline) {
	headerCell := func(s string) *tview.TableCell {
		return tview.NewTableCell(s).
			SetBackgroundColor(tcell.ColorBlueViolet).
			SetTextColor(tcell.ColorBlack).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
	}

	severityCell := func(s string, severity string) *tview.TableCell {
		var color tcell.Color
		switch {
		case severity == "high" || severity == "critical":
			color = tcell.ColorRed
		case severity == "medium":
			color = tcell.ColorOrange
		case severity == "err":
			color = tcell.ColorViolet
		}
		return tview.NewTableCell(s).
			SetTextColor(color)
	}

	headers := []string{"timestamp", "status", "message id", "severity"}
	for i, key := range headers {
		overviewTable.SetCell(0, i, headerCell(key))
	}

	data := p.GetUIDataOverview()
	for r, values := range data {
		severity := values[severityCol]
		for c, v := range values {
			overviewTable.SetCell(r+1, c, severityCell(v, severity))
		}
	}
	overviewTable.ScrollToBeginning()
}

func updateLogsBox(p *Pipeline) {
	logsBox.Clear()
	data := getSelectionData(p)
	cat := ""
	for _, line := range data.logs {
		cat += line + "\n"
	}

	logsBox.SetText(cat)
}

func updateResultsBox(p *Pipeline) {
	resultsBox.Clear()
	data := getSelectionData(p)
	if data == nil {
		return
	}
	b, err := json.MarshalIndent(data.result, "", "\t")
	if err != nil {
		return
	}
	cat := string(b[:])
	resultsBox.SetText(cat)
}

func getSelectionMessageId() string {
	row, _ := overviewTable.GetSelection()
	v := overviewTable.GetCell(row, messageIdCol).Text
	return v
}

func getSelectionData(p *Pipeline) *StoreItem {
	data, ok := p.GetResultsOfMessage(getSelectionMessageId())
	if !ok {
		return nil
	}
	return &data
}

func handleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q':
		app.Stop()
		os.Exit(0)
	}
	return event
}
