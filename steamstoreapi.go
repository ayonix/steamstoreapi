// Steamstore - an API written in go to make requests to the Steam Api.
// The Result is a struct with (hopefully) most of the data coming from the request.
package steamstoreapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// The URL that is used to query the API
const (
	StoreApiUrl = "http://store.steampowered.com/api/appdetails/?l=%v&cc=%v&v=%v&appids=%v"
)

// Stores the locale, currency and version so they don't have to be given everytime
type StoreApi struct {
	locale, currency string
	version          uint64
}

// Returns a pointer to a new StoreApi struct with the given values:
//		locale - For example 'english'
//		currency - For example 'cc'
//		version - For example 1
func NewStoreApi(locale, currency string, version uint64) *StoreApi {
	return &StoreApi{locale, currency, version}
}

// Returns the StoreApi url as a string (without ids)
func (s *StoreApi) toUrl() string {
	return fmt.Sprintf(StoreApiUrl, s.locale, s.currency, s.version)
}

// Makes a request to the API and stores the result in a StoreResponse
//		ids - Slice of the ids the API should be queried for
//		v - Pointer to a StoreResponse struct the result will be stored in
func (s *StoreApi) Request(ids []uint64, v *StoreResponse) error {
	id_strings := make([]string, len(ids))
	for i, id := range ids {
		id_strings[i] = strconv.FormatUint(id, 10)
	}

	url := fmt.Sprintf(s.toUrl(), strings.Join(id_strings, ","))
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	d := json.NewDecoder(resp.Body)
	d.Decode(&v)
	return nil
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
			Highlighted struct {
				Name string `json:",omitempty"`
				Path string `json:",omitempty"`
			} `json:",omitempty"`
			Total json.Number
		}
		Categories []struct {
			Description string
			Id          json.Number
		}
		DetailedDescription string `json:"detailed_description"`
		Developers          []string
		Dlc                 []json.Number
		Genres              []struct {
			Description string
			Id          json.Number
		}
		HeaderImage       string `json:"header_image"`
		LegalNotice       string `json:"legal_notice"`
		LinuxRequirements struct {
			Minimum string
		} `json:"linux_requirements"`
		MacRequirements struct {
			Minimum string
		} `json:"mac_requirements"`
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
		}
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
				Packageid          json.Number
				PercentSavingsText string      `json:"percent_savings_text"`
				PercentSavings     json.Number `json:"percent_savings"`
				OptionText         string      `json:"option_text"`
				OptionDescription  string      `json:"option_description"`
			} `json:"package_groups"`
		}
		Packages       []json.Number
		PcRequirements struct {
			Minimum string
		} `json:"pc_requirements"`
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
		RequiredAge json.Number
		Reviews     string
		Screenshots []struct {
			Id            json.Number
			PathFull      string `json:"path_full"`
			PathThumbnail string `json:"path_thumbnail"`
		} `json:"required_age"`
		SteamAppid         json.Number `json:"steam_appid"`
		SupportedLanguages string      `json:supported_languages`
		Type               string
		Website            string
	}
}
