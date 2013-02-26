package firehose

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

// The firehose event format implemented here is documented in:
//	https://github.com/tumblr/parmesan/wiki/Parmesan-API

type Event struct {
	Activity    Activity
	Private     bool
	PrivateData map[string]interface{}
	Timestamp   time.Time
	Post *Post
	Like *Like
}

type Post struct {
	ID           int64
	BlogID       int64
	BlogName     string
	PostURL      string
	BlogURL      string
	Type         PostType
	Tags         []string
	Title        string
	Body         string
	Caption      string
	SourceURL    string
	SourceTitle  string
	Quote        string
	LinkURL      string	// Link URL, if post is a link
	Description  string
	Photos       []Photo
}

type Photo struct {
	Caption  string
	Alt      []AltPhoto
}

func (ph *Photo) BiggestAlt() *AltPhoto {
	var alt *AltPhoto
	var width int
	for i, _ := range ph.Alt {
		x := &ph.Alt[i]
		if alt == nil || x.Width > width {
			alt, width = x, x.Width
		}
	}
	return alt
}

type AltPhoto struct {
	Width  int
	Height int
	URL    string
}

type Like struct {
	DestPostID    int64
	DestBlogID    int64
	SourceBlogID  int64
	RootPostID    int64
	RootBlogID    int64
	ParentPostID  int64
	ParentBlogID  int64
}

// Events
const (
	CreatePost = iota
	UpdatePost
	DeletePost
	Likes
	Unlikes
	FirehoseCheckpoint
)

// Post types
type PostType byte
const (
	PostUnknown = PostType(iota)
	PostText
	PostQuote
	PostLink
	PostAnswer
	PostVideo
	PostAudio
	PostPhoto
	PostChat
)

func (pt PostType) String() string {
	switch pt {
	case PostText:
		return "text"
	case PostQuote:
		return "quote"
	case PostLink:
		return "link"
	case PostAnswer:
		return "answer"
	case PostVideo:
		return "video"
	case PostAudio:
		return "audio"
	case PostPhoto:
		return "photo"
	case PostChat:
		return "chat"
	}
	return "unknown"
}

func parsePostType(t string) PostType {
	switch t {
	case "text":
		return PostText
	case "quote":
		return PostQuote
	case "link":
		return PostLink
	case "answer":
		return PostAnswer
	case "video":
		return PostVideo
	case "audio":
		return PostAudio
	case "photo":
		return PostPhoto
	case "chat":
		return PostChat
	}
	return PostUnknown
}

// Errors
var (
	ErrParse   = errors.New("unrecognized semantics")
	ErrMissing = errors.New("missing field")
	ErrType    = errors.New("wrong type")
)

type Activity byte

func (a Activity) String() string {
	switch a {
	case CreatePost:
		return "CreatePost"
	case UpdatePost:
		return "UpdatePost"
	case DeletePost:
		return "DeletePost"
	case Likes:
		return "Likes"
	case Unlikes:
		return "Unlikes"
	case FirehoseCheckpoint:
		return "FirehoseCheckpoint"
	}
	return "Unknown"
}

func ParseActivity(activity string) (Activity, error) {
	switch activity {
	case "CreatePost":
		return CreatePost, nil
	case "UpdatePost":
		return UpdatePost, nil
	case "DeletePost":
		return DeletePost, nil
	case "Likes":
		return Likes, nil
	case "Unlikes":
		return Unlikes, nil
	case "FirehoseCheckpoint":
		return FirehoseCheckpoint, nil
	}
	return 0, ErrParse
}

func ParseEvent(m map[string]interface{}) (ev *Event, err error) {
	defer func() {
		if err != nil {
			fmt.Printf("RAW:%#v\n", m)
		}
	}()

	ev = &Event{}
	if ev.Activity, err = ParseActivity(getString(m, "activity")); err != nil {
		return nil, err
	}
	if ev.Activity == FirehoseCheckpoint {
		return ev, nil
	}
	ev.Private = getBool(m, "isPrivate")
	ev.PrivateData = getMap(m, "privateJson")
	if epochms, err := getInt64(m, "timestamp"); err != nil {
		return nil, err
	} else {
		ev.Timestamp = time.Unix(0, epochms*1e6)
	}
	switch ev.Activity {
	case CreatePost, UpdatePost, DeletePost:
		p := &Post{}
		if p.ID, err = getInt64(m, "id"); err != nil {
			return nil, err
		}
		if p.BlogID, err = getInt64(m, "blogId"); err != nil {
			return nil, err
		}

		data := getMap(m, "data")
		if data == nil {
			return nil, ErrMissing
		}

		p.BlogName = getString(data, "blog_name") 
		p.Type = parsePostType(getString(data, "type"))
		p.Title = getString(data, "title") 
		p.Body = getString(data, "body") 
		p.Caption = getString(data, "caption") 
		p.SourceURL = getString(data, "source_url") 
		p.SourceTitle = getString(data, "source_title") 
		p.Quote = getString(data, "text") 
		p.PostURL = getString(data, "post_url")
		p.BlogURL = blogFromPostURL(p.PostURL)
		p.LinkURL = getString(data, "url") 
		p.Description = getString(data, "description") 
		
		switch p.Type {
		case PostPhoto:
			p.Photos = parsePhotos(data)
		}
		
		// Tags
		if tags := getSlice(data, "tags"); tags != nil {
			for _, q := range tags {
				s, ok := q.(string)
				if ok {
					p.Tags = append(p.Tags, s)
				}
			}
		}

		ev.Post = p
		return ev, nil
	case Likes, Unlikes:
		l := &Like{}
		if l.DestPostID, err = getInt64(m, "destPostId"); err != nil {
			return nil, err
		}
		if l.DestBlogID, err = getInt64(m, "destBlogId"); err != nil {
			return nil, err
		}
		if l.SourceBlogID, err = getInt64(m, "sourceBlogId"); err != nil {
			return nil, err
		}
		if l.RootPostID, err = getInt64(m, "rootPostId"); err != nil {
			return nil, err
		}
		if l.RootBlogID, err = getInt64(m, "rootBlogId"); err != nil {
			return nil, err
		}
		if l.ParentPostID, err = getInt64(m, "parentPostId"); err != nil {
			return nil, err
		}
		if l.ParentBlogID, err = getInt64(m, "parentBlogId"); err != nil {
			return nil, err
		}
		ev.Like = l
		return ev, nil
	}
	return nil, ErrParse
}

func parsePhotos(data map[string]interface{}) []Photo {
	rawPhotos := getSlice(data, "photos")
	if rawPhotos == nil {
		return nil
	}
	photos := make([]Photo, len(rawPhotos))
	for i, rawPhoto_ := range rawPhotos {
		rawPhoto, ok := rawPhoto_.(map[string]interface{})
		if !ok {
			continue
		}
		var photo *Photo = &photos[i]
		photo.Caption = getString(rawPhoto, "caption")
		rawAlts := getSlice(rawPhoto, "alt_sizes")
		photo.Alt = make([]AltPhoto, len(rawAlts))
		for j, rawAlt_ := range rawAlts {
			rawAlt, ok := rawAlt_.(map[string]interface{})
			if !ok {
				continue
			}
			photo.Alt[j] = AltPhoto{
				Width:  getInt(rawAlt, "width"),
				Height: getInt(rawAlt, "height"),
				URL:    getString(rawAlt, "url"),
			}
		}

	}
	return photos
}

func blogFromPostURL(postURL string) string {
	u, err := url.ParseRequestURI(postURL)
	if err != nil {
		return ""
	}
	return u.Scheme + "://" + u.Host
}
