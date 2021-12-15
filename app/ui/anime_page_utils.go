package ui

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/khatibomar/kanna/app/core"
	"github.com/khatibomar/tohru"
)

const (
	maxRetries = 5
)

// save: Save a Episode.
func (p *AnimePage) saveEpisode(episode *tohru.Episode) error {
	input_url := episode.EpisodeUrls[0].EpisodeURL
	u, err := url.Parse(input_url)
	if err != nil {
		return err
	}
	params := u.Query()
	var animeSpecialName string
	if len(params["f"]) > 0 {
		animeSpecialName = params["f"][0]
	} else {
		animeSpecialName = params["n"][0]
		animeSpecialName = strings.Split(animeSpecialName, "\\")[0]
	}
	// nb, err := strconv.Atoi(episode.EpisodeNumber)
	// if err != nil {
	// 	return nil
	// }
	// urls, err := core.App.Client.EpisodeService.GetDownloadLinks(animeSpecialName, nb)
	// if err != nil {
	// 	return err
	// }

	// // Get the pages to download
	// pages := chapter.Attributes.Data

	// link, err := downloader.GetChapterPage(page)
	// if err != nil {
	// 	return err
	// }

	// filename := fmt.Sprintf("%s%s", p.getDownloadFolder(episode) , filepath.Ext(page))
	// filePath := filepath.Join(core.App.Config.DownloadDir, filename)
	// // Save image
	// if err = ioutil.WriteFile(filePath, link, os.ModePerm); err != nil {
	// 	return err
	// }

	return nil
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
