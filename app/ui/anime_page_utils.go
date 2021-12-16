package ui

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/khatibomar/kanna/app/core"
	"github.com/khatibomar/tohru"
)

const (
	maxRetries = 5
)

// save: Save a Episode.
func (p *AnimePage) saveEpisode(episode *tohru.Episode, errChan chan error) {
	url, err := getDwnLink(episode)
	if err != nil {
		errChan <- err
		return
	}
	filePath := p.getDownloadPath(episode)
	filename := fmt.Sprintf("%s%s.%s", filePath, episode.EpisodeName, "mp4")

	err = os.MkdirAll(filePath, 0777)
	if err != nil {
		errChan <- err
		return
	}
	cmd := exec.Command("wget", "-c", url, "-O", filename)
	if err := cmd.Run(); err != nil {
		o, _ := cmd.CombinedOutput()
		log.Println(string(o))
		errChan <- err
		return
	}
}

// save: Save a Episode.
func (p *AnimePage) streamEpisode(episode *tohru.Episode, errChan chan error) {
	url, err := getDwnLink(episode)
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

func getDwnLink(episode *tohru.Episode) (string, error) {
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
	url, err := core.App.Client.EpisodeService.GetFirstDirectDownloadLink(animeSpecialName, nb)
	if err != nil {
		return "", err
	}
	return url, nil
}

// getDownloadFolder : Get the download folder for an episode.
func (p *AnimePage) getDownloadPath(episode *tohru.Episode) string {
	animeName := p.Anime.AnimeName
	episodeNumber := episode.EpisodeNumber

	// Remove invalid characters from the folder name
	restricted := []string{"<", ">", ":", "/", "|", "?", "*", "\"", "\\", ".", ",", " "}
	for _, c := range restricted {
		animeName = strings.ReplaceAll(animeName, c, "_")
	}
	fullPath := path.Join(core.App.Config.DownloadDir, animeName, episodeNumber) + "/"

	return fullPath
}
