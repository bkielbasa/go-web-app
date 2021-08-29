package app

import (
	"context"
	"fmt"

	"github.com/bkielbasa/go-web-app/links/domain"
)

type Link struct {
	storage Storage
	serv    domain.LinksService
}

func NewLinks(serv domain.LinksService, storage Storage) Link {
	return Link{serv: serv, storage: storage}
}

type Storage interface {
	Add(ctx context.Context, targetLink, shortLink string, tags []string) error
	// Returns the target for the given ID
	Get(ctx context.Context, id string) (string, error)
}

func (l Link) Add(ctx context.Context, targetLink string, tags []string) (string, error) {
	link := l.serv.Create(targetLink, tags)
	if err := l.storage.Add(ctx, link.Target(), link.ID(), link.Tags()); err != nil {
		return "", fmt.Errorf("cannot add the link to the storage: %w", err)
	}

	shortURL := l.serv.ConstructShortURL(link)

	return shortURL, nil
}

func (l Link) ClickLink(ctx context.Context, id string) (string, error) {
	targetLink, err := l.storage.Get(ctx, id)
	if err != nil {
		return "", fmt.Errorf("cannot get the link: %w", err)
	}

	// add counting statistics

	return targetLink, nil
}
