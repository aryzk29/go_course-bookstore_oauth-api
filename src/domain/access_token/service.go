package access_token

import (
	"github.com/aryzk29/bookstore-utils-go/rest_errors"
	"github.com/aryzk29/go_course-bookstore_oauth-api/src/repository/db"
	"github.com/aryzk29/go_course-bookstore_oauth-api/src/repository/rest"
	"strings"
)

type Service interface {
	GetById(string) (*AccessToken, *rest_errors.RestErr)
	Create(AccessTokenRequest) (*AccessToken, *rest_errors.RestErr)
	UpdateExpirationTime(AccessToken) *rest_errors.RestErr
}

type service struct {
	restUsersRepo rest.RestUserRepository
	dbRepo        db.DbRepository
}

func NewService(usersRepo rest.RestUserRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepo:        dbRepo,
	}
}

func (s *service) GetById(accessTokenId string) (*AccessToken, *rest_errors.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, rest_errors.NewBadRequestError("invalid access token id")
	}

	accessToken, err := s.dbRepo.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(req AccessTokenRequest) (*AccessToken, *rest_errors.RestErr) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	//TODO: Support both grant types: client_credentials and password

	// Authenticate the user against the Users API:
	user, err := s.restUsersRepo.LoginUser(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	// Generate a new access token:
	at := GetNewAccessToken(user.Id)
	at.Generate()

	// Save the new access token in Cassandra:
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at AccessToken) *rest_errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}
