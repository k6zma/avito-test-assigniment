package mappers

import (
	dtos "github.com/k6zma/avito-test-assigniment/api/generated"
	"github.com/k6zma/avito-test-assigniment/internal/domain/models"
)

func ToDTOPRShort(pr *models.PullRequest) dtos.PullRequestShort {
	return dtos.PullRequestShort{
		PullRequestId:   pr.ID.String(),
		PullRequestName: pr.Name,
		AuthorId:        pr.AuthorID.String(),
		Status:          dtos.PullRequestShortStatus(pr.Status),
	}
}
