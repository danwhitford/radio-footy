package feeds

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSolentDayToMergedMatches(t *testing.T) {
	table := []struct {
		input  []byte
		output []MergedMatch
	}{
		{
			input: solentTestDay,
			output: []MergedMatch{
				{
					Title:       "Gloucestershire v Hampshire Hawks",
					Competition: "T20 Blast cricket",
					Date:        "Tuesday, Jun 20",
					Stations:    []string{"BBC Radio Solent"},
					Time:        "18:00",
					Datetime:    "2023-06-20T18:00:00+01:00",
				},
			},
		},
	}

	for _, test := range table {
		got, err := getSolentDay(test.input)
		if err != nil {
			t.Errorf("getSolentDay() error: %v", err)
		}
		if diff := cmp.Diff(test.output, got); diff != "" {
			t.Errorf("getSolentDay(%s) mismatch (-want +got)", diff)
		}
	}
}

var solentTestDay = []byte(`{"$schema":"https://rms.api.bbc.co.uk/docs/swagger.json#/definitions/ExperienceResponse","data":[{"type":"inline_display_module","id":"schedule_items","style":null,"title":"2023-06-20","description":null,"state":"ok","uris":null,"controls":null,"total":null,"data":[{"type":"broadcast_summary","id":"p0frdd6d","urn":"urn:bbc:radio:episode:p0fr6wdm","start":"2023-06-19T21:00:00Z","end":"2023-06-20T00:00:00Z","service_id":"bbc_radio_solent","duration":10800,"network":{"id":"bbc_radio_solent","key":"radiosolent","short_title":"Radio Solent","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_radio_solent/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Paul Miller","secondary":"Unexpected events","tertiary":""},"synopses":{"short":"Lively late-night chat and stories, all fuelled by your calls, texts, emails and tweets.","medium":null,"long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0fdrskt.jpg","container":{"type":"brand","id":"p001d7lw","urn":"urn:bbc:radio:brand:p001d7lw","title":"Paul Miller","activities":[]},"playable_item":{"type":"playable_item","id":"p0fr6wdr","urn":"urn:bbc:radio:episode:p0fr6wdm","network":{"id":"bbc_local_radio","key":null,"short_title":"Local Radio","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_local_radio/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Paul Miller","secondary":"Unexpected events","tertiary":null},"synopses":{"short":"Lively late-night chat and stories, all fuelled by your calls, texts, emails and tweets.","medium":null,"long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0fdrskt.jpg","duration":{"value":10800,"label":"180 mins"},"progress":null,"container":{"type":"brand","id":"p001d7lw","urn":"urn:bbc:radio:brand:p001d7lw","title":"Paul Miller","synopses":{"short":"Lively late-night chat and stories, all fuelled by your calls, texts, emails and tweets.","medium":null,"long":null},"activities":[]},"download":{"type":"drm","quality_variants":{"low":{"bitrate":96,"file_url":null,"file_size":144000000,"label":"144 MB"},"medium":{"bitrate":128,"file_url":null,"file_size":192000000,"label":"192 MB"},"high":{"bitrate":320,"file_url":null,"file_size":477000000,"label":"477 MB"}}},"availability":{"from":"2023-06-20T00:00:00Z","to":"2023-07-20T00:00:00Z","label":"Available for 26 days"},"release":{"date":"2023-06-19T00:00:00Z","label":"19 Jun 2023"},"guidance":{"competition_warning":false,"warnings":null},"activities":[],"uris":[{"type":"latest","id":null,"label":"Latest","uri":"/v2/programmes/playable?container=p001d7lw&sort=sequential&type=episode&experience=domestic"}],"play_context":null,"recommendation":null}},{"type":"broadcast_summary","id":"p0frddzg","urn":"urn:bbc:radio:episode:p0fql53f","start":"2023-06-20T00:00:00Z","end":"2023-06-20T04:00:00Z","service_id":"bbc_radio_solent","duration":14400,"network":{"id":"bbc_radio_solent","key":"radiosolent","short_title":"Radio Solent","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_radio_solent/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Weekday overnights with BBC Radio 5 live","secondary":"20/06/2023","tertiary":""},"synopses":{"short":"Local BBC Radio joins BBC Radio 5 live through the night on a weekday.","medium":null,"long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p08hys6s.jpg","container":{"type":"brand","id":"p08692ml","urn":"urn:bbc:radio:brand:p08692ml","title":"Weekday overnights with BBC Radio 5 live","activities":[]},"playable_item":null},{"type":"broadcast_summary","id":"p0frddzx","urn":"urn:bbc:radio:episode:p0fql54z","start":"2023-06-20T04:00:00Z","end":"2023-06-20T05:00:00Z","service_id":"bbc_radio_solent","duration":3600,"network":{"id":"bbc_radio_solent","key":"radiosolent","short_title":"Radio Solent","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_radio_solent/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Jason Rosam","secondary":"20/06/2023","tertiary":""},"synopses":{"short":"The early breakfast show with everything you need to know to start your day.","medium":null,"long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0d9qj47.jpg","container":{"type":"brand","id":"p01nmqg2","urn":"urn:bbc:radio:brand:p01nmqg2","title":"Jason Rosam","activities":[]},"playable_item":{"type":"playable_item","id":"p0fql55k","urn":"urn:bbc:radio:episode:p0fql54z","network":{"id":"bbc_local_radio","key":null,"short_title":"Local Radio","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_local_radio/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Jason Rosam","secondary":"20/06/2023","tertiary":null},"synopses":{"short":"The early breakfast show with everything you need to know to start your day.","medium":null,"long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0d9qj47.jpg","duration":{"value":7200,"label":"120 mins"},"progress":null,"container":{"type":"brand","id":"p01nmqg2","urn":"urn:bbc:radio:brand:p01nmqg2","title":"Jason Rosam","synopses":{"short":"Join Jason Rosam on BBC Radio London.","medium":null,"long":null},"activities":[]},"download":{"type":"drm","quality_variants":{"low":{"bitrate":96,"file_url":null,"file_size":97000000,"label":"97 MB"},"medium":{"bitrate":128,"file_url":null,"file_size":128000000,"label":"128 MB"},"high":{"bitrate":320,"file_url":null,"file_size":318000000,"label":"318 MB"}}},"availability":{"from":"2023-06-20T06:00:00Z","to":"2023-07-20T06:00:00Z","label":"Available for 26 days"},"release":{"date":"2023-06-20T00:00:00Z","label":"20 Jun 2023"},"guidance":{"competition_warning":false,"warnings":null},"activities":[],"uris":[{"type":"latest","id":null,"label":"Latest","uri":"/v2/programmes/playable?container=p01nmqg2&sort=sequential&type=episode&experience=domestic"}],"play_context":null,"recommendation":null}},{"type":"broadcast_summary","id":"p0frdf1l","urn":"urn:bbc:radio:episode:p0frdf0k","start":"2023-06-20T05:00:00Z","end":"2023-06-20T09:00:00Z","service_id":"bbc_radio_solent","duration":14400,"network":{"id":"bbc_radio_solent","key":"radiosolent","short_title":"Radio Solent","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_radio_solent/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Alun Newman","secondary":"20/06/2023","tertiary":""},"synopses":{"short":"All you need for the day ahead, the South's news and great music at breakfast with Alun.","medium":null,"long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0fdrpnt.jpg","container":{"type":"brand","id":"p01n3ch8","urn":"urn:bbc:radio:brand:p01n3ch8","title":"Alun Newman","activities":[]},"playable_item":{"type":"playable_item","id":"p0frdf14","urn":"urn:bbc:radio:episode:p0frdf0k","network":{"id":"bbc_radio_solent","key":"radiosolent","short_title":"Radio Solent","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_radio_solent/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Alun Newman","secondary":"20/06/2023","tertiary":null},"synopses":{"short":"All you need for the day ahead, the South's news and great music at breakfast with Alun.","medium":null,"long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0fdrpnt.jpg","duration":{"value":14400,"label":"240 mins"},"progress":null,"container":{"type":"brand","id":"p01n3ch8","urn":"urn:bbc:radio:brand:p01n3ch8","title":"Alun Newman","synopses":{"short":"Join Alun Newman on BBC Radio Solent.","medium":"Join Alun Newman on BBC Radio Solent.","long":null},"activities":[]},"download":{"type":"drm","quality_variants":{"low":{"bitrate":96,"file_url":null,"file_size":192000000,"label":"192 MB"},"medium":{"bitrate":128,"file_url":null,"file_size":255000000,"label":"255 MB"},"high":{"bitrate":320,"file_url":null,"file_size":635000000,"label":"635 MB"}}},"availability":{"from":"2023-06-20T09:00:00Z","to":"2023-07-20T09:00:00Z","label":"Available for 26 days"},"release":{"date":"2023-06-20T00:00:00Z","label":"20 Jun 2023"},"guidance":{"competition_warning":false,"warnings":null},"activities":[],"uris":[{"type":"latest","id":null,"label":"Latest","uri":"/v2/programmes/playable?container=p01n3ch8&sort=sequential&type=episode&experience=domestic"}],"play_context":null,"recommendation":null}},{"type":"broadcast_summary","id":"p0frdf34","urn":"urn:bbc:radio:episode:p0frdf23","start":"2023-06-20T09:00:00Z","end":"2023-06-20T13:00:00Z","service_id":"bbc_radio_solent","duration":14400,"network":{"id":"bbc_radio_solent","key":"radiosolent","short_title":"Radio Solent","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_radio_solent/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Pat Sissons","secondary":"Gregg Wallace, Russell Watson and World Refugee Day","tertiary":""},"synopses":{"short":"Stories of the day, celebrity guests, your calls and fantastic music in the Solent Years.","medium":null,"long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0fdrq1w.jpg","container":{"type":"brand","id":"p00m9d46","urn":"urn:bbc:radio:brand:p00m9d46","title":"Pat Sissons","activities":[]},"playable_item":{"type":"playable_item","id":"p0frdf2p","urn":"urn:bbc:radio:episode:p0frdf23","network":{"id":"bbc_radio_solent","key":"radiosolent","short_title":"Radio Solent","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_radio_solent/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Pat Sissons","secondary":"Gregg Wallace, Russell Watson and World Refugee Day","tertiary":null},"synopses":{"short":"Stories of the day, celebrity guests, your calls and fantastic music in the Solent Years.","medium":null,"long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0fdrq1w.jpg","duration":{"value":14400,"label":"240 mins"},"progress":null,"container":{"type":"brand","id":"p00m9d46","urn":"urn:bbc:radio:brand:p00m9d46","title":"Pat Sissons","synopses":{"short":"Join Pat Sissons on BBC Radio Solent.","medium":null,"long":null},"activities":[]},"download":{"type":"drm","quality_variants":{"low":{"bitrate":96,"file_url":null,"file_size":192000000,"label":"192 MB"},"medium":{"bitrate":128,"file_url":null,"file_size":255000000,"label":"255 MB"},"high":{"bitrate":320,"file_url":null,"file_size":635000000,"label":"635 MB"}}},"availability":{"from":"2023-06-20T13:00:00Z","to":"2023-07-20T13:00:00Z","label":"Available for 26 days"},"release":{"date":"2023-06-20T00:00:00Z","label":"20 Jun 2023"},"guidance":{"competition_warning":false,"warnings":null},"activities":[],"uris":[{"type":"latest","id":null,"label":"Latest","uri":"/v2/programmes/playable?container=p00m9d46&sort=sequential&type=episode&experience=domestic"}],"play_context":null,"recommendation":null}},{"type":"broadcast_summary","id":"p0frdf4f","urn":"urn:bbc:radio:episode:p0frdf3q","start":"2023-06-20T13:00:00Z","end":"2023-06-20T17:00:00Z","service_id":"bbc_radio_solent","duration":14400,"network":{"id":"bbc_radio_solent","key":"radiosolent","short_title":"Radio Solent","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_radio_solent/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Lou Hannan","secondary":"20/06/2023","tertiary":""},"synopses":{"short":"Supercharging your afternoon with news, guests, music, and inspirational local stories.","medium":"Supercharging your afternoon with news, guests and great music. Plus, inspirational local stories.","long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0fdrqj8.jpg","container":{"type":"brand","id":"p0036nbt","urn":"urn:bbc:radio:brand:p0036nbt","title":"Lou Hannan","activities":[]},"playable_item":{"type":"playable_item","id":"p0frdf43","urn":"urn:bbc:radio:episode:p0frdf3q","network":{"id":"bbc_radio_solent","key":"radiosolent","short_title":"Radio Solent","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_radio_solent/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Lou Hannan","secondary":"20/06/2023","tertiary":null},"synopses":{"short":"Supercharging your afternoon with news, guests, music, and inspirational local stories.","medium":"Supercharging your afternoon with news, guests and great music. Plus, inspirational local stories.","long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0fdrqj8.jpg","duration":{"value":14400,"label":"240 mins"},"progress":null,"container":{"type":"brand","id":"p0036nbt","urn":"urn:bbc:radio:brand:p0036nbt","title":"Lou Hannan","synopses":{"short":"Join Lou Hannan on BBC Radio Solent.","medium":null,"long":null},"activities":[]},"download":{"type":"drm","quality_variants":{"low":{"bitrate":96,"file_url":null,"file_size":192000000,"label":"192 MB"},"medium":{"bitrate":128,"file_url":null,"file_size":255000000,"label":"255 MB"},"high":{"bitrate":320,"file_url":null,"file_size":635000000,"label":"635 MB"}}},"availability":{"from":"2023-06-20T17:00:00Z","to":"2023-07-20T17:00:00Z","label":"Available for 27 days"},"release":{"date":"2023-06-20T00:00:00Z","label":"20 Jun 2023"},"guidance":{"competition_warning":false,"warnings":null},"activities":[],"uris":[{"type":"latest","id":null,"label":"Latest","uri":"/v2/programmes/playable?container=p0036nbt&sort=sequential&type=episode&experience=domestic"}],"play_context":null,"recommendation":null}},{"type":"broadcast_summary","id":"p0frdf5b","urn":"urn:bbc:radio:episode:p0frdf4m","start":"2023-06-20T17:00:00Z","end":"2023-06-20T21:00:00Z","service_id":"bbc_radio_solent","duration":14400,"network":{"id":"bbc_radio_solent","key":"radiosolent","short_title":"Radio Solent","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_radio_solent/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Summer Sport","secondary":"Gloucestershire v Hampshire Hawks (20/06/2023)","tertiary":""},"synopses":{"short":"T20 Blast cricket coverage of Gloucestershire v Hampshire Hawks (weather permitting).","medium":null,"long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0cjdyns.jpg","container":{"type":"brand","id":"p02rfdhr","urn":"urn:bbc:radio:brand:p02rfdhr","title":"Summer Sport","activities":[]},"playable_item":{"type":"playable_item","id":"p0frdf4w","urn":"urn:bbc:radio:episode:p0frdf4m","network":{"id":"bbc_radio_solent","key":"radiosolent","short_title":"Radio Solent","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_radio_solent/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Summer Sport","secondary":"Gloucestershire v Hampshire Hawks (20/06/2023)","tertiary":null},"synopses":{"short":"T20 Blast cricket coverage of Gloucestershire v Hampshire Hawks (weather permitting).","medium":null,"long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0cjdyns.jpg","duration":{"value":14400,"label":"240 mins"},"progress":null,"container":{"type":"brand","id":"p02rfdhr","urn":"urn:bbc:radio:brand:p02rfdhr","title":"Summer Sport","synopses":{"short":"Local summer sports coverage and updates (weather permitting), and a selection of music.","medium":null,"long":null},"activities":[]},"download":{"type":"drm","quality_variants":{"low":{"bitrate":96,"file_url":null,"file_size":192000000,"label":"192 MB"},"medium":{"bitrate":128,"file_url":null,"file_size":255000000,"label":"255 MB"},"high":{"bitrate":320,"file_url":null,"file_size":635000000,"label":"635 MB"}}},"availability":{"from":"2023-06-20T21:00:00Z","to":"2023-07-20T21:00:00Z","label":"Available for 27 days"},"release":{"date":"2023-06-20T00:00:00Z","label":"20 Jun 2023"},"guidance":{"competition_warning":false,"warnings":null},"activities":[],"uris":[{"type":"latest","id":null,"label":"Latest","uri":"/v2/programmes/playable?container=p02rfdhr&sort=sequential&type=episode&experience=domestic"}],"play_context":null,"recommendation":null}},{"type":"broadcast_summary","id":"p0frdf61","urn":"urn:bbc:radio:episode:p0fr6xf5","start":"2023-06-20T21:00:00Z","end":"2023-06-21T00:00:00Z","service_id":"bbc_radio_solent","duration":10800,"network":{"id":"bbc_radio_solent","key":"radiosolent","short_title":"Radio Solent","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_radio_solent/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Paul Miller","secondary":"Holiday friends","tertiary":""},"synopses":{"short":"Lively late-night chat and stories, all fuelled by your calls, texts, emails and tweets.","medium":null,"long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0fdrskt.jpg","container":{"type":"brand","id":"p001d7lw","urn":"urn:bbc:radio:brand:p001d7lw","title":"Paul Miller","activities":[]},"playable_item":{"type":"playable_item","id":"p0fr6xfp","urn":"urn:bbc:radio:episode:p0fr6xf5","network":{"id":"bbc_local_radio","key":null,"short_title":"Local Radio","logo_url":"https://sounds.files.bbci.co.uk/3.5.0/networks/bbc_local_radio/{type}_{size}.{format}","network_type":"master_brand"},"titles":{"primary":"Paul Miller","secondary":"Holiday friends","tertiary":null},"synopses":{"short":"Lively late-night chat and stories, all fuelled by your calls, texts, emails and tweets.","medium":null,"long":null},"image_url":"https://ichef.bbci.co.uk/images/ic/{recipe}/p0fdrskt.jpg","duration":{"value":10800,"label":"180 mins"},"progress":null,"container":{"type":"brand","id":"p001d7lw","urn":"urn:bbc:radio:brand:p001d7lw","title":"Paul Miller","synopses":{"short":"Lively late-night chat and stories, all fuelled by your calls, texts, emails and tweets.","medium":null,"long":null},"activities":[]},"download":{"type":"drm","quality_variants":{"low":{"bitrate":96,"file_url":null,"file_size":144000000,"label":"144 MB"},"medium":{"bitrate":128,"file_url":null,"file_size":192000000,"label":"192 MB"},"high":{"bitrate":320,"file_url":null,"file_size":477000000,"label":"477 MB"}}},"availability":{"from":"2023-06-21T00:00:00Z","to":"2023-07-21T00:00:00Z","label":"Available for 27 days"},"release":{"date":"2023-06-20T00:00:00Z","label":"20 Jun 2023"},"guidance":{"competition_warning":false,"warnings":null},"activities":[],"uris":[{"type":"latest","id":null,"label":"Latest","uri":"/v2/programmes/playable?container=p001d7lw&sort=sequential&type=episode&experience=domestic"}],"play_context":null,"recommendation":null}}],"image_url":null}]}`)
