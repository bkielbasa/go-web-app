package infra

import (
	"context"

	"github.com/bkielbasa/go-web-app/links/domain"
)

type inMemory struct {
	links map[string]string
}

func NewInMemoryStorage() *inMemory {
	return &inMemory{
		links: map[string]string{},
	}
}

func (im *inMemory) Add(ctx context.Context, targetLink, id string, tags []string) error {
	im.links[id] = targetLink
	return nil
}

func (im *inMemory) Get(ctx context.Context, id string) (string, error) {
	if link, ok := im.links[id]; ok {
		return link, nil
	}

	return "", domain.ErrLinkNotFound
}
