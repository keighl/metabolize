package metabolize

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_Metabolize(t *testing.T) {
	doc := `
	<html prefix="og: http://ogp.me/ns#">
	<head>
	<title>The Rock (1996)</title>
	<meta property="title" content="The Rock" />
	<meta property="og:title" content="The Rock" />
	</head>
	</html>
	`

	obj := struct {
		Title string `meta:"title"`
	}{}
	err := Metabolize(strings.NewReader(doc), &obj)
	expect(t, err, nil)
	expect(t, obj.Title, "The Rock")
}

func Test_ParseDocument(t *testing.T) {
	doc := `
<html prefix="og: http://ogp.me/ns#">
<head>
<title>The Rock (1996)</title>
<meta property="" content="your mom" />
<meta property="title" content="The Rock" />
<meta property="og:title" content="The Rock" />
<meta property="og:TYPE" content="video.movie" />
<meta property="og:url" content="http://www.imdb.com/title/tt0117500/" />
<meta property="og:image" content="http://ia.media-imdb.com/images/rock.jpg" />
<meta property="og:image" content="http://ia.media-imdb.com/images/cheese.jpg" />
<meta property="og:cows    " content="mooo" />
</head>
</html>
	`
	expectedData := MetaData{
		"title":    "The Rock",
		"og:title": "The Rock",
		"og:type":  "video.movie",
		"og:url":   "http://www.imdb.com/title/tt0117500/",
		"og:image": "http://ia.media-imdb.com/images/cheese.jpg",
		"og:cows":  "mooo",
	}

	data, err := ParseDocument(strings.NewReader(doc))
	expect(t, err, nil)
	expect(t, len(data), len(expectedData))
	expect(t, reflect.DeepEqual(data, expectedData), true)
}

func Test_ParseDocument_Whacky(t *testing.T) {
	doc := `____`
	_, err := ParseDocument(strings.NewReader(doc))
	expect(t, err, nil)
}

func Test_Decode_NoTag(t *testing.T) {
	data := MetaData{}
	obj := struct {
		Title string
	}{}
	err := Decode(data, &obj)
	expect(t, err, nil)
}

func Test_Decode_NotAStruct(t *testing.T) {
	data := MetaData{}
	obj := "CHEESE"
	err := Decode(data, &obj)
	expect(t, err, NotStructError)
}

func Test_Decode_String(t *testing.T) {
	data := MetaData{
		"title":    "The Rock",
		"og:type":  "video.movie",
		"og:image": "http://ia.media-imdb.com/images/cheese.jpg",
	}

	obj := struct {
		Title string `meta:"title"`
		Type  string `meta:"og:type"`
		Image string `meta:"og:image"`
	}{}

	err := Decode(data, &obj)
	expect(t, err, nil)
	expect(t, obj.Title, "The Rock")
	expect(t, obj.Type, data["og:type"])
	expect(t, obj.Image, data["og:image"])
}

func Test_Decode_URL_Invalid(t *testing.T) {
	data := MetaData{
		"og:image": "CHEESE",
	}

	obj := struct {
		Image url.URL `meta:"og:image"`
	}{}

	err := Decode(data, &obj)
	expect(t, err, nil)
	u, _ := url.Parse(data["og:image"])
	expect(t, obj.Image, *u)
}

func Test_Decode_URL(t *testing.T) {
	data := MetaData{
		"og:image": "http://ia.media-imdb.com/images/cheese.jpg",
	}

	obj := struct {
		Image url.URL `meta:"og:image"`
	}{}

	err := Decode(data, &obj)
	expect(t, err, nil)
	u, _ := url.Parse(data["og:image"])
	expect(t, obj.Image, *u)
}

func Test_Decode_ISO8601(t *testing.T) {
	data := MetaData{
		"og:published_at": "2015-09-14T17:51:31+00:00",
	}

	obj := struct {
		PublishedAt time.Time `meta:"og:published_at"`
	}{}

	err := Decode(data, &obj)
	expect(t, err, nil)
	val, _ := time.Parse(time.RFC3339, data["og:published_at"])
	expect(t, obj.PublishedAt.String(), val.String())
}

func Test_Decode_ISO8601_Invalid(t *testing.T) {
	data := MetaData{
		"og:published_at": "CHEESE",
	}

	obj := struct {
		PublishedAt time.Time `meta:"og:published_at"`
	}{}

	err := Decode(data, &obj)
	expect(t, err, nil)
	val, _ := time.Parse(time.RFC3339, data["og:published_at"])
	expect(t, obj.PublishedAt.String(), val.String())
}

func Test_Decode_Int(t *testing.T) {
	data := MetaData{
		"widgets": "-20",
	}

	obj := struct {
		Widgets int64 `meta:"widgets"`
	}{}

	err := Decode(data, &obj)
	expect(t, err, nil)
	expect(t, obj.Widgets, int64(-20))
}

func Test_Decode_Int_Invalid(t *testing.T) {
	data := MetaData{
		"widgets": "CHEESE",
	}

	obj := struct {
		Widgets int64 `meta:"widgets"`
	}{}

	err := Decode(data, &obj)
	expect(t, err, nil)
	expect(t, obj.Widgets, int64(0))
}

func Test_Decode_Float(t *testing.T) {
	data := MetaData{
		"widgets": "-2.0",
	}

	obj := struct {
		Widgets float64 `meta:"widgets"`
	}{}

	err := Decode(data, &obj)
	expect(t, err, nil)
	expect(t, obj.Widgets, float64(-2.0))
}

func Test_Decode_Float_Invalid(t *testing.T) {
	data := MetaData{
		"widgets": "CHEESE",
	}

	obj := struct {
		Widgets float64 `meta:"widgets"`
	}{}

	err := Decode(data, &obj)
	expect(t, err, nil)
	expect(t, obj.Widgets, float64(0))
}

func Test_Decode_Bool(t *testing.T) {
	data := MetaData{
		"published_1": "true",
		"published_2": "false",
		"published_3": "1",
		"published_4": "0",
	}

	obj := struct {
		Published1 bool `meta:"published_1"`
		Published2 bool `meta:"published_2"`
		Published3 bool `meta:"published_3"`
		Published4 bool `meta:"published_4"`
	}{}

	err := Decode(data, &obj)
	expect(t, err, nil)
	expect(t, obj.Published1, true)
	expect(t, obj.Published2, false)
	expect(t, obj.Published3, true)
	expect(t, obj.Published4, false)
}

func Test_Decode_Boo_Invalid(t *testing.T) {
	data := MetaData{
		"published": "CHEES",
	}

	obj := struct {
		Published bool `meta:"published"`
	}{}

	err := Decode(data, &obj)
	expect(t, err, nil)
	expect(t, obj.Published, false)
}

func Test_Decode_Fallbacks(t *testing.T) {
	data := MetaData{
		"title": "Your mom",
	}

	obj := struct {
		Title string `meta:"og:title,title"`
	}{}

	err := Decode(data, &obj)
	expect(t, err, nil)
	expect(t, obj.Title, data["title"])
}
