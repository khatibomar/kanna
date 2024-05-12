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

	"github.com/khatibomar/tohru"
)

func (p *AnimePage) saveEpisodes(episodes []*tohru.Episode, errChan chan error) {
	if len(episodes) == 0 {
		return
	}
	errCount := 0
	filePath := p.getDownloadPath(p.Core.Config.DownloadDir)
	queueName := p.getAnimeNameWithYear()
	for _, episode := range episodes {
		url, err := getDwnLink(episode, p.Core.Client.EpisodeService.GetFirstDirectDownloadInfo)
		if err != nil {
			errChan <- fmt.Errorf("%s: %s", episode.EpisodeName, err)
			errCount++
			continue
		}
		fullPath := fmt.Sprintf("%s%s.%s", filePath, removeRestrictedChars(episode.EpisodeName), "mp4")

		log.Printf("downloading episode with id : %s , name : %s\n from server %s \nto %s", episode.EpisodeID, episode.EpisodeName, url, fullPath)

		err = os.MkdirAll(filePath, 0777)
		if err != nil {
			errChan <- fmt.Errorf("%s: %s", episode.EpisodeName, err)
			continue
		}
		err = p.Core.Fafnir.Add(queueName, url.EpisodeDirectDownloadLink, filePath, fmt.Sprintf("%s.%s", removeRestrictedChars(episode.EpisodeName), "mp4"), url.EpisodeHostLink)
		if err != nil {
			errChan <- err
		}
	}
	if errCount > 0 {
		errChan <- fmt.Errorf("%d errors appeared please check logs", errCount)
	}
	err := p.Core.Fafnir.StartQueueDownload(queueName)
	if err != nil {
		errChan <- err
	}
}

func (p *AnimePage) streamEpisode(episode *tohru.Episode, errChan chan error) {
	url, err := getDwnLink(episode, p.Core.Client.EpisodeService.GetFirstDirectDownloadInfo)
	if err != nil {
		errChan <- err
		return
	}
	log.Printf("streaming episode with id : %s , name : %s\n from server %s", episode.EpisodeID, episode.EpisodeName, url)
	mpv := exec.Command("mpv", url.EpisodeDirectDownloadLink)
	if err := mpv.Start(); err != nil {
		errChan <- fmt.Errorf("%q: failed to start mpv", err)
		return
	}
}

func getDwnLink(episode *tohru.Episode, getFirstDwnLinkF func(string, int) (tohru.DownloadInfo, error)) (tohru.DownloadInfo, error) {
	if len(episode.EpisodeUrls) == 0 {
		return tohru.DownloadInfo{}, fmt.Errorf("no Download links available")
	}
	input_url := episode.EpisodeUrls[0].EpisodeURL
	u, err := url.Parse(input_url)
	if err != nil {
		return tohru.DownloadInfo{}, err
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
		return tohru.DownloadInfo{}, err
	}
	url, err := getFirstDwnLinkF(animeSpecialName, nb)
	if err != nil {
		return tohru.DownloadInfo{}, err
	}
	return url, nil
}

// getDownloadFolder : Get the download folder for an episode.
func (p *AnimePage) getDownloadPath(dwnDir string) string {
	animeName := removeRestrictedChars(p.Anime.AnimeName)
	animeYear := p.Anime.AnimeReleaseYear

	// Remove invalid characters from the folder name
	fullPath := path.Join(dwnDir, animeName+"_"+animeYear) + "/"

	return fullPath
}

func (p *AnimePage) getAnimeNameWithYear() string {
	animeName := removeRestrictedChars(p.Anime.AnimeName)
	animeYear := p.Anime.AnimeReleaseYear
	return animeName + "_" + animeYear
}

func removeRestrictedChars(s string) string {
	restricted := []string{"<", ">", ":", "/", "|", "?", "*", "\"", "\\", ".", ",", " "}
	for _, c := range restricted {
		s = strings.ReplaceAll(s, c, "_")
	}
	return s
}
