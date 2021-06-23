package voice

import "fmt"

// Webhook _
type Webhook struct {
	ObjectKind string `json:"object_kind"`

	User             user          `json:"user"`
	Project          project       `json:"project"`
	ObjectAttributes *attributes   `json:"object_attributes"`
	MergeRequest     *mergeRequest `json:"merge_request"`
	Issue            *issue        `json:"issue"`
}

type user struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type project struct {
	Path string `json:"path_with_namespace"`
	URL  string `json:"url"`
}

type attributes struct {
	ID           int    `json:"id"`
	Note         string `json:"note"`
	NoteableType string `json:"noteable_type"`
	IID          int    `json:"iid"`
	Title        string `json:"title"`
	State        string `json:"state"`
	URL          string `json:"url"`
	Action       string `json:"action"`
}

type mergeRequest struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	State string `json:"state"`
	IID   int    `json:"iid"`
}

type issue struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	State string `json:"state"`
	IID   int    `json:"iid"`
}

// Notification _
func (wh *Webhook) Notification() string {
	switch wh.ObjectKind {
	case "merge_request":
		return wh.mrNotification()
	case "issue":
		return wh.issueNotification()
	case "note":
		return wh.commentNotification()
	default:
		fmt.Println("webhook", wh.ObjectKind)
	}
	return ""
}

func (wh *Webhook) mrNotification() string {
	switch wh.ObjectAttributes.Action {
	case "open", "merge", "close":
		return fmt.Sprintf("%s\n%s MR [\\!%d](%s) \"%s\" at %s",
			markdownEscape(wh.User.Username),
			wh.ObjectAttributes.Action,
			wh.ObjectAttributes.IID,
			wh.ObjectAttributes.URL,
			markdownEscape(wh.ObjectAttributes.Title),
			markdownEscape(wh.Project.Path),
		)
	default:
		return ""
	}
}

func (wh *Webhook) issueNotification() string {
	switch wh.ObjectAttributes.Action {
	case "open", "merge", "close":
		return fmt.Sprintf("%s\n%s issue [\\#%d](%s) \"%s\" at %s",
			markdownEscape(wh.User.Username),
			wh.ObjectAttributes.Action,
			wh.ObjectAttributes.IID,
			wh.ObjectAttributes.URL,
			markdownEscape(wh.ObjectAttributes.Title),
			markdownEscape(wh.Project.Path),
		)
	default:
		return ""
	}
}

func (wh *Webhook) commentNotification() string {
	switch wh.ObjectAttributes.NoteableType {
	//case "Commit":
	//return ""
	case "MergeRequest":
		return fmt.Sprintf("%s\ncomment [\\!%d](%s) \"%s\" at %s\n%s",
			markdownEscape(wh.User.Username),
			wh.MergeRequest.IID,
			wh.ObjectAttributes.URL,
			markdownEscape(wh.MergeRequest.Title),
			markdownEscape(wh.Project.Path),
			markdownEscape(wh.ObjectAttributes.Note),
		)
	case "Issue":
		return fmt.Sprintf("%s\ncomment [\\#%d](%s) \"%s\" at %s\n%s",
			markdownEscape(wh.User.Username),
			wh.Issue.IID,
			wh.ObjectAttributes.URL,
			markdownEscape(wh.Issue.Title),
			markdownEscape(wh.Project.Path),
			markdownEscape(wh.ObjectAttributes.Note),
		)
	default:
		return ""
	}
}
