package ui

import (
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/khatibomar/chunky/dwn"
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
	filename := fmt.Sprintf("%s%s", episode.EpisodeName, filepath.Ext("mp4"))
	filePath := p.getDownloadFolder(episode)

	d := dwn.NewFileDownloader(url, filename, filePath)
	errChan <- d.Download()
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
	urls, err := core.App.Client.EpisodeService.GetDirectDownloadLinks(animeSpecialName, nb)
	if err != nil {
		return "", err
	}
	if len(urls) == 0 {
		return "", fmt.Errorf("No direct links available")
	}
	return urls[0], nil
}

// getDownloadFolder : Get the download folder for a manga's chapter.
func (p *AnimePage) getDownloadFolder(episode *tohru.Episode) string {
	animeName := p.Anime.AnimeName
	episodeNumber := episode.EpisodeNumber
	// Remove invalid characters from the folder name
	restricted := []string{"<", ">", ":", "/", "|", "?", "*", "\"", "\\", "."}
	for _, c := range restricted {
		animeName = strings.ReplaceAll(animeName, c, "-")
		episodeNumber = strings.ReplaceAll(episodeNumber, c, "-")
	}

	folder := filepath.Join(core.App.Config.DownloadDir, animeName, episodeNumber)
	return folder
}
