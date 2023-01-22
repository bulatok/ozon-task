package usecase

import (
	"errors"
	"go.uber.org/zap"
	"net/http"

	"github.com/bulatok/ozon-task/internal/ozon-task/models"
	"github.com/bulatok/ozon-task/internal/ozon-task/store"
)

var (
	ErrLinksProhibited = errors.New("the same short link condition")
)

type Links struct {
	repo   store.LinksRepo
	logger *zap.Logger
}

func ProvideLinks(repo store.LinksRepo, l *zap.Logger) *Links {
	childL := l.With(zap.String("logger", "links_usecase"))
	return &Links{
		repo:   repo,
		logger: childL,
	}
}

// New saves the short and original link and returns the short link
//
// baseServiceLink + "/" + hash(originalLink)
func (l *Links) New(originalLink, baseServiceLink string) (string, error) {
	link := &models.Link{
		Original: originalLink,
	}

	if err := link.SetShortLink(baseServiceLink); err != nil {
		l.logger.Error("could not set short link",
			zap.String("error", err.Error()))
		return "", models.ErrInternalServer
	}

	// checking that the same short link does not exist
	existingLink, err := l.repo.Get(link.Short)

	switch err {
	case models.ErrNotFound:
		if err := l.repo.Save(link); err != nil {
			l.logger.Error("could not save the link",
				zap.String("error", err.Error()))
			return "", models.ErrInternalServer
		}
		return link.Short, nil
	case nil:
		if originalLink == existingLink.Original {
			return link.Short, nil
		}
		l.logger.Info("the same short link condition")
		return "", models.NewCommonErr(ErrLinksProhibited.Error(), http.StatusInternalServerError)
	default:
		l.logger.Error("repo error while getting from it", zap.String("error", err.Error()))
		return "", models.ErrInternalServer
	}
}

// Get returns original link from shortLink
func (l *Links) Get(shortLink string) (string, error) {
	link, err := l.repo.Get(shortLink)
	if err != nil {
		if err == models.ErrNotFound {
			return "", models.ErrNotFound
		}
		l.logger.Error("could not get the link",
			zap.String("error", err.Error()))
		return "", models.ErrInternalServer
	}

	return link.Original, nil
}
