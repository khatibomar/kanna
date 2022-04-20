package ui

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"codeberg.org/omarkhatib/kanna/app/core"
	"codeberg.org/omarkhatib/kanna/app/ui/utils"
	"codeberg.org/omarkhatib/tohru"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	EpisodesOffsetRange   = 500
	contextCancelledError = "CANCELLED"
	readStatus            = "Y"
)

type EpisodeSelection struct {
	animeID string
	episode *tohru.Episode
}

// AnimePage : This struct contains the required primitives for the anime page.
type AnimePage struct {
	Anime *tohru.AnimeDetails
	Grid  *tview.Grid
	Info  *tview.TextView
	Table *tview.Table
	Core  *core.Kanna

	sWrap *utils.SelectorWrapper
	cWrap *utils.ContextWrapper // For context cancellation.
}

// ShowAnimePage : Make the app show the anime page.
func ShowAnimePage(core *core.Kanna, anime *tohru.Anime) {
	id, err := strconv.Atoi(anime.AnimeID)
	if err != nil {
		log.Println(err)
		return
	}
	animeDetails, err := core.Client.AnimeService.GetAnimeDetails(id)
	if err != nil {
		log.Println(err)
		return
	}
	animePage := newAnimePage(core, &animeDetails)

	core.TView.SetFocus(animePage.Grid)
	core.PageHolder.AddAndSwitchToPage(utils.AnimePageID, animePage.Grid, true)
}

// newAnimePage : Creates a new anime page.
func newAnimePage(core *core.Kanna, anime *tohru.AnimeDetails) *AnimePage {
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

	// Use a table to show the episodes for the anime.
	table := tview.NewTable()

	// Set episode headers
	titleHead := tview.NewTableCell("Episode Number").
		SetTextColor(utils.AnimePageChapNumColor).
		SetSelectable(false)
	table.SetCell(0, 0, titleHead).
		SetFixed(1, 0)
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(utils.AnimePageTableBorderColor).
		SetTitle("Episodes").
		SetTitleColor(utils.AnimePageTableTitleColor).
		SetBorder(true)

	// Add info and table to the grid. Set the focus to the episode table.
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
		Core: core,
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

	p.Core.TView.QueueUpdateDraw(func() {
		p.Info.SetText(infoText)
	})
}

// setChapterTable : Fill up the episodes table.
func (p *AnimePage) setEpisodesTable() {
	log.Println("Setting up anime page episodes table...")
	ctx, cancel := p.cWrap.ResetContext()
	p.setHandlers(cancel)

	time.Sleep(loadDelay)
	defer cancel()

	p.Core.TView.QueueUpdateDraw(func() {
		loadingCell := tview.NewTableCell("Loading...").SetSelectable(false)
		p.Table.SetCell(1, 1, loadingCell)
	})

	if p.cWrap.ToCancel(ctx) {
		return
	}
	id, _ := strconv.Atoi(p.Anime.AnimeID)
	episodes, err := p.getAllEpisodes(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), contextCancelledError) {
			return
		}
		log.Println(fmt.Sprintf("Error getting anime episodes: %s", err.Error()))
		p.Core.TView.QueueUpdateDraw(func() {
			modal := okModal(p.Core, utils.GenericAPIErrorModalID, "Error getting anime episodes.\nCheck log for details.")
			ShowModal(p.Core, utils.GenericAPIErrorModalID, modal)
		})
		return
	} else if len(episodes) == 0 {
		p.Core.TView.QueueUpdateDraw(func() {
			noResultsCell := tview.NewTableCell("No episodes!").SetSelectable(false)
			p.Table.SetCell(1, 1, noResultsCell)
		})
		return
	}
	markers := map[string]struct{}{}
	// Fill in the episodes
	for index := 0; index < len(episodes); index++ {
		if p.cWrap.ToCancel(ctx) {
			return
		}
		selection := EpisodeSelection{
			animeID: p.Anime.AnimeID,
			episode: &episodes[index],
		}

		// Episode Number
		episodeText := fmt.Sprintf("الحلقة: %d", index+1)
		episodeNumCell := tview.NewTableCell(
			fmt.Sprintf("%-6s", episodeText)).
			SetMaxWidth(10).SetTextColor(utils.AnimePageChapNumColor).SetReference(selection)
		// Read marker
		var read string
		if _, ok := markers[selection.episode.EpisodeID]; ok {
			read = readStatus
		}
		readCell := tview.NewTableCell(read).SetTextColor(utils.AnimePageWatchStatColor)
		p.Table.SetCell(index+1, 0, episodeNumCell).
			SetCell(index+1, 1, readCell)
	}

	p.Core.TView.QueueUpdateDraw(func() {
		p.Table.Select(1, 0)
		p.Table.ScrollToBeginning()
	})
}

// getAllChapters : Get All episodes for the anime.
func (p *AnimePage) getAllEpisodes(ctx context.Context, animeID int) ([]tohru.Episode, error) {
	var (
		episodes   []tohru.Episode
		currOffset = 0
	)
	for {
		if p.cWrap.ToCancel(ctx) {
			return []tohru.Episode{}, fmt.Errorf(contextCancelledError)
		}
		list, err := p.Core.Client.EpisodeService.GetEpisodesList(animeID)
		if err != nil {
			return []tohru.Episode{}, err
		}
		log.Printf("Got %d of %d episodes\n", currOffset, len(list))
		episodes = list
		currOffset += EpisodesOffsetRange
		if currOffset >= len(list) {
			break
		}
	}
	return episodes, nil
}

// markSelected : Mark an episode as being selected by the user on the main page table.
func (p *AnimePage) markSelected(row int) {
	episodeCell := p.Table.GetCell(row, 0)
	episodeCell.SetTextColor(tcell.ColorBlack).SetBackgroundColor(utils.AnimePageHighlightColor)

	// Add to the Selection wrapper
	p.sWrap.AddSelection(row)
}

// markUnselected : Mark an episode as being unselected by the user on the main page table.
func (p *AnimePage) markUnselected(row int) {
	episodeCell := p.Table.GetCell(row, 0)
	episodeCell.SetTextColor(utils.AnimePageChapNumColor).SetBackgroundColor(tcell.ColorBlack)

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
