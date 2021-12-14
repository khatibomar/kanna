package ui

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/khatibomar/angoslayer"
	"github.com/khatibomar/tkanna/app/core"
	"github.com/khatibomar/tkanna/app/ui/utils"
	"github.com/rivo/tview"
)

const (
	EpisodesOffsetRange   = 500
	contextCancelledError = "CANCELLED"
	readStatus            = "Y"
)

// AnimePage : This struct contains the required primitives for the anime page.
type AnimePage struct {
	Anime *angoslayer.AnimeDetails
	Grid  *tview.Grid
	Info  *tview.TextView
	Table *tview.Table

	sWrap *utils.SelectorWrapper
	cWrap *utils.ContextWrapper // For context cancellation.
}

// ShowAnimePage : Make the app show the anime page.
func ShowAnimePage(anime *angoslayer.Anime) {
	id, err := strconv.Atoi(anime.AnimeID)
	if err != nil {
		log.Fatalln(err)
	}
	animeDetails, err := core.App.Client.AnimeService.GetAnimeDetails(id)
	if err != nil {
		log.Fatalln(err)
	}
	animePage := newAnimePage(&animeDetails)

	core.App.TView.SetFocus(animePage.Grid)
	core.App.PageHolder.AddAndSwitchToPage(utils.AnimePageID, animePage.Grid, true)
}

// newAnimePage : Creates a new anime page.
func newAnimePage(anime *angoslayer.AnimeDetails) *AnimePage {
	var dimensions []int
	for i := 0; i < 15; i++ {
		dimensions = append(dimensions, -1)
	}
	grid := utils.NewGrid(dimensions, dimensions)
	// Set grid attributes
	grid.SetTitleColor(utils.AnimePageGridTitleColor).
		SetBorderColor(utils.AnimePageGridBorderColor).
		SetTitle("Anime Information").
		SetBorder(true)

	// Use a TextView for basic information of the anime.
	info := tview.NewTextView()
	// Set textview attributes
	info.SetWrap(true).SetWordWrap(true).
		SetBorderColor(utils.AnimePageInfoViewBorderColor).
		SetTitleColor(utils.AnimePageInfoViewTitleColor).
		SetTitle("About").
		SetBorder(true)

	// Use a table to show the chapters for the anime.
	table := tview.NewTable()
	// Set chapter headers
	numHeader := tview.NewTableCell("Chap").
		SetTextColor(utils.AnimePageChapNumColor).
		SetSelectable(false)
	titleHeader := tview.NewTableCell("Name").
		SetTextColor(utils.AnimePageTitleColor).
		SetSelectable(false)
	downloadHeader := tview.NewTableCell("Download Status").
		SetTextColor(utils.AnimePageDownloadStatColor).
		SetSelectable(false)
	scanGroupHeader := tview.NewTableCell("ScanGroup").
		SetTextColor(utils.AnimePageScanGroupColor).
		SetSelectable(false)
	readMarkerHeader := tview.NewTableCell("Read Status").
		SetTextColor(utils.AnimePageReadStatColor).
		SetSelectable(false)
	table.SetCell(0, 0, numHeader).
		SetCell(0, 1, titleHeader).
		SetCell(0, 2, downloadHeader).
		SetCell(0, 3, scanGroupHeader).
		SetCell(0, 4, readMarkerHeader).
		SetFixed(1, 0)
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(utils.AnimePageTableBorderColor).
		SetTitle("Chapters").
		SetTitleColor(utils.AnimePageTableTitleColor).
		SetBorder(true)

	// Add info and table to the grid. Set the focus to the chapter table.
	grid.AddItem(info, 0, 0, 5, 15, 0, 0, false).
		AddItem(table, 5, 0, 10, 15, 0, 0, true).
		AddItem(info, 0, 0, 15, 5, 0, 80, false).
		AddItem(table, 0, 5, 15, 10, 0, 80, true)

	ctx, cancel := context.WithCancel(context.Background())
	animePage := &AnimePage{
		Anime: anime,
		Grid:  grid,
		Info:  info,
		Table: table,
		sWrap: &utils.SelectorWrapper{
			Selection: map[int]struct{}{},
		},
		cWrap: &utils.ContextWrapper{
			Ctx:    ctx,
			Cancel: cancel,
		},
	}

	// Set up values
	go animePage.setAnimeInfo()
	go animePage.setEpisodesTable()

	return animePage
}

// setAnimeInfo: Set up anime information.
func (p *AnimePage) setAnimeInfo() {
	// Title
	var title string
	if p.Anime.AnimeEnglishTitle != "" {
		title = p.Anime.AnimeEnglishTitle
	} else {
		title = p.Anime.AnimeName
	}

	// Status
	status := p.Anime.AnimeStatus

	// Description
	desc := tview.Escape(p.Anime.AnimeDescription)

	// Set up information text.
	infoText := fmt.Sprintf("Title: %s\n\nStatus: %s\n\nDescription:\n%s",
		title, status, desc)

	core.App.TView.QueueUpdateDraw(func() {
		p.Info.SetText(infoText)
	})
}

// setChapterTable : Fill up the episodes table.
func (p *AnimePage) setEpisodesTable() {
	log.Println("Setting up anime page episodes table...")
	ctx, cancel := p.cWrap.ResetContext()
	// Set handlers.
	p.setHandlers(cancel)

	time.Sleep(loadDelay)
	defer cancel()

	// Show loading status so user knows it's loading.
	core.App.TView.QueueUpdateDraw(func() {
		loadingCell := tview.NewTableCell("Loading...").SetSelectable(false)
		p.Table.SetCell(1, 1, loadingCell)
	})

	// Get All chapters
	if p.cWrap.ToCancel(ctx) {
		return
	}
	id, _ := strconv.Atoi(p.Anime.AnimeID)
	episodes, err := p.getAllEpisodes(ctx, id)
	if err != nil { // If error getting chapters.
		if strings.Contains(err.Error(), contextCancelledError) {
			return
		}
		log.Println(fmt.Sprintf("Error getting anime episodes: %s", err.Error()))
		core.App.TView.QueueUpdateDraw(func() {
			modal := okModal(utils.GenericAPIErrorModalID, "Error getting anime episodes.\nCheck log for details.")
			ShowModal(utils.GenericAPIErrorModalID, modal)
		})
		return
	} else if len(episodes) == 0 { // If there are no chapters.
		core.App.TView.QueueUpdateDraw(func() {
			noResultsCell := tview.NewTableCell("No episodes!").SetSelectable(false)
			p.Table.SetCell(1, 1, noResultsCell)
		})
		return
	}

	markers := map[string]struct{}{}
	// Fill in the chapters
	for index := 0; index < len(episodes); index++ {
		if p.cWrap.ToCancel(ctx) {
			return
		}
		episode := episodes[index]
		// Chapter Number
		chapterNumCell := tview.NewTableCell(
			fmt.Sprintf("%-6s", episode.EpisodeNumber)).
			SetMaxWidth(10).SetTextColor(utils.AnimePageChapNumColor).SetReference(&episode)

		// Chapter title
		titleCell := tview.NewTableCell(fmt.Sprintf("%-30s", episode.EpisodeName)).SetMaxWidth(30).
			SetTextColor(utils.AnimePageTitleColor)

		// Chapter download status
		var downloadStatus string
		// Check for the presence of the download folder.
		if _, err = os.Stat(p.getDownloadFolder(&episode)); err == nil {
			downloadStatus = "Y"
		}
		downloadCell := tview.NewTableCell(downloadStatus).SetTextColor(utils.AnimePageDownloadStatColor)

		// Read marker
		var read string
		if _, ok := markers[episode.EpisodeID]; ok {
			read = readStatus
		}
		readCell := tview.NewTableCell(read).SetTextColor(utils.AnimePageReadStatColor)

		p.Table.SetCell(index+1, 0, chapterNumCell).
			SetCell(index+1, 1, titleCell).
			SetCell(index+1, 2, downloadCell)

		p.Table.SetCell(index+1, 4, readCell)
	}
	core.App.TView.QueueUpdateDraw(func() {
		p.Table.Select(1, 0)
		p.Table.ScrollToBeginning()
	})
}

// getAllChapters : Get All episodes for the anime.
func (p *AnimePage) getAllEpisodes(ctx context.Context, animeID int) ([]angoslayer.Episode, error) {
	var (
		episodes   []angoslayer.Episode
		currOffset = 0
	)
	for {
		if p.cWrap.ToCancel(ctx) {
			return []angoslayer.Episode{}, fmt.Errorf(contextCancelledError)
		}
		list, err := core.App.Client.EpisodeService.GetEpisodesList(animeID)
		if err != nil {
			return []angoslayer.Episode{}, err
		}
		log.Printf("Got %d of %d chapters\n", currOffset, len(list))
		episodes = list
		currOffset += EpisodesOffsetRange
		if currOffset >= len(list) {
			break
		}
	}
	return episodes, nil
}

// markSelected : Mark a chapter as being selected by the user on the main page table.
func (p *AnimePage) markSelected(row int) {
	chapterCell := p.Table.GetCell(row, 0)
	chapterCell.SetTextColor(tcell.ColorBlack).SetBackgroundColor(utils.AnimePageHighlightColor)

	// Add to the Selection wrapper
	p.sWrap.AddSelection(row)
}

// markUnselected : Mark a chapter as being unselected by the user on the main page table.
func (p *AnimePage) markUnselected(row int) {
	chapterCell := p.Table.GetCell(row, 0)
	chapterCell.SetTextColor(utils.AnimePageChapNumColor).SetBackgroundColor(tcell.ColorBlack)

	// Remove from the Selection wrapper
	p.sWrap.RemoveSelection(row)
}

// markAll : Marks All rows as selected or unselected.
func (p *AnimePage) markAll() {
	if p.sWrap.All {
		for row := 1; row < p.Table.GetRowCount(); row++ {
			p.markUnselected(row)
		}
	} else {
		for row := 1; row < p.Table.GetRowCount(); row++ {
			p.markSelected(row)
		}
	}
	p.sWrap.All = !p.sWrap.All
}
