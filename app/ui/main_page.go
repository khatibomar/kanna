package ui

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/khatibomar/tkanna/app/core"
	"github.com/khatibomar/tkanna/app/ui/utils"
	"github.com/rivo/tview"
)

const (
	limit     = 100
	loadDelay = time.Millisecond * 50
	maxOffset = 10000
)

// MainPage : This struct contains the grid and the entry table.
type MainPage struct {
	Grid          *tview.Grid
	Table         *tview.Table
	CurrentOffset int

	cWrap *utils.ContextWrapper // For context cancellation.
}

// ShowMainPage : Make the app show the main page.
func ShowMainPage() {
	// Create the new main page
	log.Println("Creating new main page...")
	mainPage := newMainPage()

	core.App.TView.SetFocus(mainPage.Grid)
	core.App.PageHolder.AddAndSwitchToPage(utils.MainPageID, mainPage.Grid, true)
}

// newMainPage : Creates a new main page.
func newMainPage() *MainPage {
	var dimensions []int
	for i := 0; i < 15; i++ {
		dimensions = append(dimensions, -1)
	}
	grid := utils.NewGrid(dimensions, dimensions)
	// Set grid attributes.
	grid.SetTitleColor(utils.MainPageGridTitleColor).
		SetBorderColor(utils.MainPageGridBorderColor).
		SetBorder(true)

	// Create the base main table.
	table := tview.NewTable()
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(utils.MainPageTableBorderColor).
		SetTitleColor(utils.MainPageTableTitleColor).
		SetBorder(true)

	// Add the table to the grid. Table spans the whole page.
	grid.AddItem(table, 0, 0, 15, 15, 0, 0, true)

	ctx, cancel := context.WithCancel(context.Background())
	mainPage := &MainPage{
		Grid:  grid,
		Table: table,
		cWrap: &utils.ContextWrapper{
			Ctx:    ctx,
			Cancel: cancel,
		},
	}

	go mainPage.setTableGrid()
	go mainPage.setLatestUpdatedAnimeTable()
	return mainPage
}

func (p *MainPage) setTableGrid() {
	log.Println("Setting anime grid...")
	core.App.TView.QueueUpdateDraw(func() {
		name, _ := os.Hostname()
		p.Grid.SetTitle(fmt.Sprintf("Welcome to Tkanna, [yellow]%s!", name))
	})
	log.Println("Finished setting table grid.")
}

func (p *MainPage) setLatestUpdatedAnimeTable() {
	log.Println("Setting latest updated anime table...")
	ctx, cancel := p.cWrap.ResetContext()
	p.setHandlers(cancel)

	time.Sleep(loadDelay)
	defer cancel()

	tableTitle := "Latest Updated Anime"

	core.App.TView.QueueUpdateDraw(func() {
		// Clear current entries
		p.Table.Clear()

		// Set headers.
		titleHeader := tview.NewTableCell("Anime").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageTitleColor).
			SetSelectable(false)
		descHeader := tview.NewTableCell("Rating").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageDescColor).
			SetSelectable(false)
		tagHeader := tview.NewTableCell("Episode").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageTagColor).
			SetSelectable(false)
		p.Table.SetCell(0, 0, titleHeader).
			SetCell(0, 1, descHeader).
			SetCell(0, 2, tagHeader).
			SetFixed(1, 0)

		// Set table title.
		page, first, last := p.calculatePaginationData()
		p.Table.SetTitle(fmt.Sprintf("%s. Page %d (%d-%d). [::bu]Loading...", tableTitle, page, first, last))
	})

	// Get list of Animes.
	if p.cWrap.ToCancel(ctx) {
		return
	}
	list, err := core.App.Client.AnimeService.GetLatestAnimes(p.CurrentOffset, limit)
	if err != nil {
		log.Println(err.Error())
		core.App.TView.QueueUpdateDraw(func() {
			modal := okModal(utils.GenericAPIErrorModalID, "Error getting anime list.\nCheck logs for details.")
			ShowModal(utils.GenericAPIErrorModalID, modal)
		})
		return
	}

	// Update table title.
	page, first, last := p.calculatePaginationData()
	core.App.TView.QueueUpdateDraw(func() {
		p.Table.SetTitle(fmt.Sprintf("%s. Page %d (%d-%d).", tableTitle, page, first, last))
	})

	// Fill in the details
	for index := 0; index < len(list); index++ {
		if p.cWrap.ToCancel(ctx) {
			return
		}
		anime := list[index]
		// Anime title cell.
		mtCell := tview.NewTableCell(fmt.Sprintf("%-40s", anime.AnimeName)).
			SetMaxWidth(400).SetTextColor(utils.GuestMainPageTitleColor).SetReference(&anime)

		// Rating cell.
		desc := tview.Escape(fmt.Sprintf("%-60s",
			strings.SplitN(tview.Escape(anime.AnimeRating), "\n", 2)[0]))
		descCell := tview.NewTableCell(desc).SetMaxWidth(50).SetTextColor(utils.GuestMainPageDescColor)

		// Episode cell.
		tagCell := tview.NewTableCell(anime.LatestEpisodeName).SetTextColor(utils.GuestMainPageTagColor)

		p.Table.SetCell(index+1, 0, mtCell).
			SetCell(index+1, 1, descCell).
			SetCell(index+1, 2, tagCell)
	}
	core.App.TView.QueueUpdateDraw(func() {
		p.Table.Select(1, 0)
		p.Table.ScrollToBeginning()
	})
	log.Println("Finished setting latest updated anime table.")
}

// calculatePaginationData : Calculates the current page and first/last entry number.
// Returns (pageNo, firstEntry, lastEntry).
func (p *MainPage) calculatePaginationData() (int, int, int) {
	page := p.CurrentOffset/limit + 1
	firstEntry := p.CurrentOffset + 1
	lastEntry := page * limit

	if firstEntry > lastEntry {
		firstEntry = lastEntry
	}

	return page, firstEntry, lastEntry
}
