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

	"codeberg.org/omarkhatib/tohru"
	"github.com/cavaliergopher/grab/v3"
)

const (
	maxRetries = 5
)

func (p *AnimePage) saveEpisode(episode *tohru.Episode, errChan chan error, infoChan chan string) {
	url, err := getDwnLink(episode, p.Core.Client.EpisodeService.GetFirstDirectDownloadInfo)
	if err != nil {
		errChan <- err
		return
	}
	filePath := p.getDownloadPath(episode, p.Core.Config.DownloadDir)
	fullPath := fmt.Sprintf("%s%s.%s", filePath, removeRestrictedChars(episode.EpisodeName), "mp4")

	log.Printf("downloading episode with id : %s , name : %s\n from server %s \nto %s", episode.EpisodeID, episode.EpisodeName, url, fullPath)

	err = os.MkdirAll(filePath, 0777)
	if err != nil {
		errChan <- err
		return
	}
	resp, err := grab.Get(fullPath, url.EpisodeDirectDownloadLink)
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
		log.Printf("download complete and saved to %s\n", fullPath)
		return
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
		return tohru.DownloadInfo{}, fmt.Errorf("No Download links available")
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
func (p *AnimePage) getDownloadPath(episode *tohru.Episode, dwnDir string) string {
	animeName := removeRestrictedChars(p.Anime.AnimeName)
	animeYear := p.Anime.AnimeReleaseYear

	// Remove invalid characters from the folder name
	fullPath := path.Join(dwnDir, animeName+"_"+animeYear) + "/"

	return fullPath
}

func removeRestrictedChars(s string) string {
	restricted := []string{"<", ">", ":", "/", "|", "?", "*", "\"", "\\", ".", ",", " "}
	for _, c := range restricted {
		s = strings.ReplaceAll(s, c, "_")
	}
	return s
}
