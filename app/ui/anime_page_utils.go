package ui

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/cavaliercoder/grab"
	"github.com/khatibomar/tohru"
)

const (
	maxRetries = 5
)

// save: Save a Episode.
func (p *AnimePage) saveEpisode(episode *tohru.Episode, errChan chan error, infoChan chan string) {
	url, err := getDwnLink(episode, p.Core.Client.EpisodeService.GetFirstDirectDownloadLink)
	if err != nil {
		errChan <- err
		return
	}
	filePath := p.getDownloadPath(episode, p.Core.Config.DownloadDir)
	fullPath := fmt.Sprintf("%s%s.%s", filePath, removeRestrictedChars(episode.EpisodeName), "mp4")

	err = os.MkdirAll(filePath, 0777)
	if err != nil {
		errChan <- err
		return
	}
	resp, err := grab.Get(fullPath, url)
	if err != nil {
		errChan <- err
		return
	}
	if resp.Err() != nil {
		errChan <- resp.Err()
		return
	}
	if resp.IsComplete() {
		infoChan <- fmt.Sprintf("Download is complete and file can be found at: %s", fullPath)
		return
	}

}

// save: Save a Episode.
func (p *AnimePage) streamEpisode(episode *tohru.Episode, errChan chan error) {
	url, err := getDwnLink(episode, p.Core.Client.EpisodeService.GetFirstDirectDownloadLink)
	if err != nil {
		errChan <- err
		return
	}
	mpv := exec.Command("mpv", url)
	if err := mpv.Start(); err != nil {
		errChan <- fmt.Errorf("%q: failed to start mpv", err)
		return
	}
}

func getDwnLink(episode *tohru.Episode, getFirstDwnLinkF func(string, int) (string, error)) (string, error) {
	if len(episode.EpisodeUrls) == 0 {
		return "", fmt.Errorf("No Download links available")
	}
	input_url := episode.EpisodeUrls[0].EpisodeURL
	u, err := url.Parse(input_url)
	if err != nil {
		return "", err
	}
	params := u.Query()
	var animeSpecialName string
	if len(params["f"]) > 0 {
		animeSpecialName = params["f"][0]
	} else {
		animeSpecialName = params["n"][0]
		animeSpecialName = strings.Split(animeSpecialName, "\\")[0]
	}
	nb, err := strconv.Atoi(episode.EpisodeNumber)
	if err != nil {
		return "", err
	}
	url, err := getFirstDwnLinkF(animeSpecialName, nb)
	if err != nil {
		return "", err
	}
	return url, nil
}

// getDownloadFolder : Get the download folder for an episode.
func (p *AnimePage) getDownloadPath(episode *tohru.Episode, dwnDir string) string {
	animeName := removeRestrictedChars(p.Anime.AnimeName)
	episodeNumber := episode.EpisodeNumber

	// Remove invalid characters from the folder name
	fullPath := path.Join(dwnDir, animeName, episodeNumber) + "/"

	return fullPath
}

func removeRestrictedChars(s string) string {
	restricted := []string{"<", ">", ":", "/", "|", "?", "*", "\"", "\\", ".", ",", " "}
	for _, c := range restricted {
		s = strings.ReplaceAll(s, c, "_")
	}
	return s
}
