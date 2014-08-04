// Steamstore - an API written in go to make requests to the Steam Api.
// The Result is a struct with (hopefully) most of the data coming from the request.
package steamstoreapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// The URL that is used to query the API
const (
	storeApiUrl = "http://store.steampowered.com/api/appdetails/?l=%v&cc=%v&v=%v&appids=%v"
)

// Stores the locale, currency and version so they don't have to be given everytime
type storeApi struct {
	locale, currency string
	version          int64
}

// Returns a pointer to a new StoreApi struct with the given values:
//		locale - For example 'english'
//		currency - For example 'cc'
//		version - For example 1
func newStoreApi(locale, currency string, version int64) *storeApi {
	return &storeApi{locale, currency, version}
}

// Returns the StoreApi url as a string (without ids)
func (s *storeApi) toUrl(ids []uint64) string {
	id_strings := make([]string, len(ids))
	for i, id := range ids {
		id_strings[i] = strconv.FormatUint(id, 10)
	}
	return fmt.Sprintf(storeApiUrl, s.locale, s.currency, s.version, strings.Join(id_strings, ","))
}

// Makes a request to the API and stores the result in a StoreResponse
//		ids - Slice of the ids the API should be queried for
//		v - Pointer to a StoreResponse struct the result will be stored in
func (s *storeApi) request(ids []uint64, ch chan StoreResponse, errch chan error) {
	var v StoreResponse
	url := s.toUrl(ids)
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		errch <- err
		return
	}

	if resp.StatusCode != 200 {
		errch <- errors.New(fmt.Sprintf("Server responded with status %d", resp.StatusCode))
		return
	}

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&v)
	if err != nil {
		errch <- err
	} else {
		ch <- v
	}
}

func GetStoreResponse(ids []uint64, locale, currency string) (resp StoreResponse, err error) {
	rcount := 0
	store := newStoreApi(locale, currency, 1)
	resp = make(StoreResponse)

	step, lo := 50, 0
	hi := step

	ch := make(chan StoreResponse)
	errors := make(chan error)

	for hi < len(ids) {
		if lo == hi {
			break
		}

		rcount++
		go store.request(ids[lo:hi], ch, errors)

		// Update bounds
		lo = hi
		if hi+step > len(ids) {
			hi = len(ids) - 1
		} else {
			hi += step
		}
	}

	// wait for all goroutines to finish
	for i := 0; i < rcount; i++ {
		select {
		case err := <-errors:
			return nil, err
		case r := <-ch:
			// add the results to the returned map
			for k, v := range r {
				resp[k] = v
			}
		}
	}
	return
}

// The response from the API is a map[string]AppResponse
// with the appid as key
type StoreResponse map[string]AppResponse

// A struct that covers at least most of the data coming
// from the API
type AppResponse struct {
	Success bool
	Data    struct {
		AboutTheGame string `json:"about_the_game"`
		Achievements struct {
			Highlighted []struct {
				Name string `json:"name"`
				Path string `json:"path"`
			} `json:"highlighted,omitempty"`
			Total json.Number
		} `json:",omitempty"`
		Categories []struct {
			Description string
			Id          json.Number
		}
		DetailedDescription string `json:"detailed_description"`
		Developers          []string
		Dlc                 []json.Number `json:",omitempty"`
		Genres              []struct {
			Description string
			Id          json.Number
		}
		HeaderImage string `json:"header_image"`
		LegalNotice string `json:"legal_notice"`
		// LinuxRequirements struct {
		// Minimum     string
		// Recommended string
		// } `json:"linux_requirements,omitempty"`
		// MacRequirements struct {
		// Minimum     string
		// Recommended string
		// } `json:"mac_requirements,omitempty"`
		Metacritic struct {
			Score json.Number
			Url   string
		}
		Movies []struct {
			Highlight bool
			Id        json.Number
			Name      string
			Thumbnail string
			Webm      map[string]string
		} `json:",omitempty"`
		Name          string
		PackageGroups []struct {
			Name                    string
			Title                   string
			Description             string
			SelectionText           string      `json:"selection_text"`
			SaveText                string      `json:"save_text"`
			DisplayType             json.Number `json:"display_type"`
			IsRecurringSubscription string      `json:"is_recurring_subscription"`
			Subs                    []struct {
				CanGetFreeLicense  json.Number `json:"can_get_free_license,omitempty"`
				PackageId          json.Number `json:"packageid"`
				PercentSavingsText string      `json:"percent_savings_text,omitempty"`
				PercentSavings     json.Number `json:"percent_savings,omitempty"`
				OptionText         string      `json:"option_text,omitempty"`
				OptionDescription  string      `json:"option_description,omitempty"`
			} `json:"package_groups,omitempty"`
		}
		Packages []json.Number
		// PcRequirements struct {
		// Minimum     string
		// Recommended string
		// } `json:"pc_requirements,omitempty"`
		Platforms     map[string]bool
		PriceOverview struct {
			Currency        string
			Initial         json.Number
			Final           json.Number
			DiscountPercent json.Number `json:"discount_percent"`
		} `json:"price_overview"`
		Publishers      []string
		Recommendations struct {
			Total json.Number
		}
		ReleaseDate struct {
			ComingSoon bool `json:"coming_soon"`
			Date       string
		} `json:"release_date"`
		RequiredAge json.Number `json:"required_age"`
		Reviews     string
		Screenshots []struct {
			Id            json.Number
			PathFull      string `json:"path_full"`
			PathThumbnail string `json:"path_thumbnail"`
		}
		SteamAppid  json.Number `json:"steam_appid"`
		SupportInfo struct {
			Email string
			Url   string
		} `json:"support_info,omitempty"`
		SupportedLanguages string `json:"supported_languages"`
		Type               string
		Website            string
	}
}
