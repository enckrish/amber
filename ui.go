package main

import (
	"amber/pipeline"
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

func runUI(p *pipeline.Pipeline) {
	app.SetInputCapture(handleInput)

	view := buildLayout(p)
	err := app.SetRoot(view, true).EnableMouse(true).Run()
	if err != nil {
		panic(err)
	}
}

func buildLayout(p *pipeline.Pipeline) tview.Primitive {
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

func getOverviewTable(p *pipeline.Pipeline) *tview.Table {
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

func updateOverviewTableData(p *pipeline.Pipeline) {
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

	data := getUIDataFromPipeline(p)
	for r, values := range data {
		severity := values[severityCol]
		for c, v := range values {
			overviewTable.SetCell(r+1, c, severityCell(v, severity))
		}
	}
	overviewTable.ScrollToBeginning()
}

func updateLogsBox(p *pipeline.Pipeline) {
	logsBox.Clear()
	data := getSelectionData(p)
	cat := ""
	for _, line := range data.Logs() {
		cat += line + "\n"
	}

	logsBox.SetText(cat)
}

func updateResultsBox(p *pipeline.Pipeline) {
	resultsBox.Clear()
	data := getSelectionData(p)
	if data == nil {
		return
	}
	b, err := json.MarshalIndent(data.Result(), "", "\t")
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

func getSelectionData(p *pipeline.Pipeline) *pipeline.StoreItem {
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

func onUpdateSample(p *pipeline.Pipeline) {
	//app.Draw()
	if overviewTable != nil {
		updateOverviewTableData(p)
		app.Draw()
	}
}

func getUIDataFromPipeline(p *pipeline.Pipeline) [][4]string {
	res := make([][4]string, 0)
	for pair := p.Store().Newest(); pair != nil; pair = pair.Prev() {
		messageId := pair.Key.(string)
		data := pair.Value.(pipeline.StoreItem)
		timestamp := data.RequestTime().Round(0).String()
		resultsFetched := data.Result() != nil
		var status, severity string
		if resultsFetched {
			status = "DONE"
			severity = data.Result().StrRating()
		} else {
			status = "WAITING"
			severity = ""
		}
		values := [4]string{timestamp, status, messageId, severity}
		res = append(res, values)
	}
	return res
}
