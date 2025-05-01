package xkcd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SearchComic ...
func SearchComic(n int) (*ComicInfo, error) {
	url := fmt.Sprintf(xkcdURL, n)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// We must close resp.Body on all execution paths.
	// defer should make this simpler
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result ComicInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}
