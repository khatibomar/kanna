package browser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

func (a Anime) Title() string       { return a.AnimeName }
func (a Anime) Description() string { return a.AnimeId }
func (a Anime) FilterValue() string { return a.Title() }

type Animes struct {
	Animes []Anime `json:"data"`
}

type Response struct {
	Res Animes `json:"response"`
}

type model struct {
	list list.Model
}

func InitialModel(cfg *config.Config) (model, error) {
	var animes []Anime
	var err error

	animes, err = getAnimeData(cfg, 0, 10)
	if err != nil {
		return model{}, err
	}
	items := []list.Item{}
	for _, anime := range animes {
		items = append(items, anime)
	}
	return model{
		list: list.NewModel(items, list.NewDefaultDelegate(), 0, 0),
	}, err
}

func (m model) Init() tea.Cmd {
	return nil
}

func getAnimeData(cfg *config.Config, offset, limit int) ([]Anime, error) {
	var animes []Anime
	url := fmt.Sprintf(LATEST_EPISODES_URL, offset, limit)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return animes, err
	}

	req.Header.Add("Client-Id", cfg.ClientID)
	req.Header.Add("Client-Secret", cfg.ClientSecret)
	req.Header.Add("Accept", "*/*")

	res, err := client.Do(req)
	if err != nil {
		return animes, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var resp Response
	err = json.Unmarshal(body, &resp)
	return resp.Res.Animes, err
}
