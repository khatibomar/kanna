package browser

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/khatibomar/tkanna/config"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

const (
	LATEST_EPISODES_URL = `https://anslayer.com/anime/public/animes/get-published-animes?json={"_offset":%d,"_limit":%d,"_order_by":"latest_first","list_type":"latest_updated_episode_new","just_info":"Yes"}`
)

type Anime struct {
	AnimeId            string `json:"anime_id,omitempty"`
	AnimeName          string `json:"anime_name,omitempty"`
	AnimeType          string `json:"anime_type"`
	AnimeStatus        string `json:"anime_status"`
	JustInfo           string `json:"just_info"`
	AnimeSeason        string `json:"anime_season"`
	AnimeReleaseYear   string `json:"anime_release_year"`
	AnimeRating        string `json:"anime_rating,omitempty"`
	LatestEpisodeID    string `json:"latest_episode_id,omitempty"`
	LatestEpisodeName  string `json:"latest_episode_name,omitempty"`
	AnimeCoverImageURL string `json:"anime_cover_image_url"`
	AnimeTrailerURL    string `json:"anime_trailer_url"`
	AnimeReleaseDay    string `json:"anime_release_day"`
}

func (a Anime) Title() string { return a.AnimeName }

func (a Anime) Description() string {
	return fmt.Sprintf("Episode:%s , Season:%s , Rating:%s", a.EpisodeName(), a.AnimeSeason, a.AnimeRating)
}

func (a Anime) FilterValue() string { return a.Title() }

func (a Anime) EpisodeName() string {
	ep := a.LatestEpisodeName
	return strings.TrimPrefix(ep, "الحلقة : ")
}

type Animes struct {
	Animes []Anime `json:"data"`
}

type Response struct {
	Res Animes `json:"response"`
}

type errMsg error

type model struct {
	list list.Model
	err  error
}

func InitialModel(cfg *config.Config) model {
	var animes []Anime
	var err error

	animes, err = getAnimeData(cfg, 0, 15)
	if err != nil {
		return model{err: err}
	}
	items := []list.Item{}
	for _, anime := range animes {
		items = append(items, anime)
	}
	m := model{
		list: list.NewModel(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "Latest added animes"
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}
