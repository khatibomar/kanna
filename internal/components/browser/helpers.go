package browser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/khatibomar/tkanna/config"
)

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
