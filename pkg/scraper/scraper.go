// Package implements utility routines to get files from the web.
package scraper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"github.com/eric-gagnon/mgl7460-tp2-linktograph-go/pkg/link"
)

func ScrapFilesToCache(sourceLinks []string, cachefolderpath string) {
	// Concurrency: https://nathanleclaire.com/blog/2014/02/15/how-to-wait-for-all-goroutines-to-finish-executing-before-continuing/
	messages := make(chan string)
	var wg sync.WaitGroup

	wg.Add(len(sourceLinks))

	for i, l := range sourceLinks {

		filename := link.GetSha1FileNameForLink(l)
		cacheFilePath := filepath.Join(cachefolderpath, filename)

		go func(link string, cacheFilePath string, index int) {
			defer wg.Done()

			// https://golangcode.com/check-if-a-file-exists/
			if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
				// todo : Ajouter retour erreur.
				downloadFileForLink(link, cacheFilePath)
				messages <- fmt.Sprintf("getFileFromLink finished for %s, starting order: %d", link, index)
			} else {
				messages <- fmt.Sprintf("Skip downloadFileForLink, file already in cache : %s", link)
			}

		}(l, cacheFilePath, i)
	}

	go func() {
		for i := range messages {
			fmt.Println(i)
		}
	}()

	wg.Wait()
}

func downloadFileForLink(link string, cacheFilePath string) {

	client := &http.Client{}

	// todo : handle the error?
	req, _ := http.NewRequest("GET", link, nil)

	resp, err := client.Do(req)

	if err != nil {
		// todo : add why.
		fmt.Printf("Skip download to cache, failed request. err: %v, link : %s,  \n", err, link)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Note : will rewrite file if the file exist (refresh).

		out, err := os.Create(cacheFilePath)
		if err != nil {
			panic(err)
		}

		defer out.Close()
		io.Copy(out, resp.Body)
		fmt.Printf("Downloaded link %s to : %s\n", link, cacheFilePath)
	}
}
