package repository

import(
	"github.com/pkg/errors"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/infrastructure/repository/in_memory"
)

func Get(entity string, repoType string) (repo event.IEventRepository, err error) {

	switch entity {
	case "event":
		repo, err = eventRepo(repoType)
	default:
		err = errors.Errorf("Entity %q not found", entity)
	}
	return repo, err
}

func eventRepo(repoType string) (repo event.IEventRepository, err error){

	switch repoType {
	case "in_memory":
		repo = in_memory.New()
	default:
		err = errors.Errorf("RepoType %q not found", repoType)
	}
	return repo, err
}
