package domain

import "errors"

var ErrLinkNotFound = errors.New("the link does not exist")

type Link struct {
	id     string
	target string
	tags   []string
}

func NewLink(id, target string, tags []string) Link {
	return Link{
		id:     id,
		target: target,
		tags:   tags,
	}
}

func (l Link) ID() string {
	return l.id
}

func (l Link) Target() string {
	return l.target
}

func (l Link) Tags() []string {
	return l.tags
}

type LinksService struct {
	baseURL string
}

func NewLinksService(baseURL string) LinksService {
	return LinksService{
		baseURL: baseURL,
	}
}

func (ls LinksService) Create(targetLink string, tags []string) Link {
	id := randString(8)
	return NewLink(id, targetLink, tags)
}

func (ls LinksService) ConstructShortURL(link Link) string {
	return ls.baseURL + "/c/" + link.ID()
}
