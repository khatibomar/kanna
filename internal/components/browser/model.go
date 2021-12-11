package browser

var (
	offset = 0
	limit  = 30
)

const (
	LATEST_EPISODES_URL = `https://anslayer.com/anime/public/animes/get-published-animes?json={"_offset":%d,"_limit":%d,"_order_by":"latest_first","list_type":"latest_updated_episode_new","just_info":"Yes"}`
)

type Anime struct {
	AnimeId            int     `json:"anime_id,omitempty"`
	AnimeName          string  `json:"anime_name,omitempty"`
	AnimeType          string  `json:"anime_type"`
	AnimeStatus        string  `json:"anime_status"`
	JustInfo           string  `json:"just_info"`
	AnimeSeason        string  `json:"anime_season"`
	AnimeReleaseYear   int     `json:"anime_release_year"`
	AnimeRating        float32 `json:"anime_rating,omitempty"`
	LatestEpisodeID    int     `json:"latest_episode_id,omitempty"`
	LatestEpisodeName  string  `json:"latest_episode_name,omitempty"`
	AnimeCoverImageURL string  `json:"anime_cover_image_url"`
	AnimeTrailerURL    string  `json:"anime_trailer_url"`
	AnimeReleaseDay    string  `json:"anime_release_day"`
}

type Animes struct {
	Animes []Anime `json:"data"`
}

type Response struct {
	Res Animes `json:"response"`
}

type model struct {
	animes   []Anime
	cursor   int
	selected map[int]int
}
