package steamstoreapi

import (
	"testing"
)

var (
	ids     = []uint64{10, 20, 30}
	id      = []uint64{10}
	id_long = []uint64{224760, 400, 41050, 209540, 2420, 4920, 99700, 241600, 3590, 102600, 205790, 304930, 223220, 70600, 10, 65800, 48240, 219680, 41070, 203210, 242110, 32460, 24800, 205910, 12900, 17480, 107100, 55110, 6120, 42910, 6850, 203140, 98200, 222730, 32360, 8850, 201790, 41060, 24740, 32800, 24980, 41000, 232910, 201420, 206440, 24780, 47790, 203810, 214970, 207570, 2400, 93200, 204360, 63000, 8930, 219190, 65530, 225260, 17460, 730, 50650, 207170, 22600, 63710, 108800, 550, 204300, 214870, 219890, 204880, 227080, 219200, 221260, 620, 41500, 49520, 6900, 644, 39690, 35720, 218620, 94200, 239070, 227780, 220780, 200390, 211820, 2430, 4000, 200710, 500, 219150, 22000, 40800, 17410, 210770, 218060, 233720, 8190, 9350, 91600, 39650, 38440, 238010, 9420, 113020, 17470, 214770, 237530, 7670, 70, 238430, 41800, 201480, 95300, 47830, 204340, 48000, 221380, 244870, 240, 223530}
)

func TestToUrl(t *testing.T) {
	storeapi := newStoreApi("en", "cc", 1)
	url := storeapi.toUrl(ids)
	if url != "http://store.steampowered.com/api/appdetails/?l=en&cc=cc&v=1&appids=10,20,30" {
		t.Errorf("Wrong url: %s", url)
	}
}

func TestRequest(t *testing.T) {
	storeapi := newStoreApi("en", "cc", 1)
	var resp StoreResponse
	err := storeapi.request(id, &resp)
	if err != nil {
		t.Error(err)
	}
}

func TestGetStoreResponse(t *testing.T) {
	_, err := GetStoreResponse(ids, "en", "cc")
	if err != nil {
		t.Error(err)
	}
}

func TestLength(t *testing.T) {
	_, err := GetStoreResponse(id_long[:60], "en", "cc")
	if err != nil {
		t.Error(err)
	}
}
