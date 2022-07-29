package db

import (
	"github.com/aryzk29/go_course-bookstore_oauth-api/src/clients/cassandra"
	"github.com/aryzk29/go_course-bookstore_oauth-api/src/domain/access_token"
	"github.com/aryzk29/go_course-bookstore_users-api/utils/errors"
	"github.com/gocql/gocql"
)

const (
	queryGetAccessToken = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	queryCreateToken    = "INSERT INTO access_tokens (access_token, user_id, client_id, expires) values (?,?,?,?);"
	queryUpdateExpires  = "UPDATE access_tokens SET expires = ? where access_token = ?;"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessToken) *errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, *errors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&result.AccessToken,
		&result.UserId,
		&result.ClientId,
		&result.Expires); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.NewNotFoundError("no access token found")
		}
		return nil, errors.NewInternalServerError(err.Error())
	}

	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) *errors.RestErr {
	if err := cassandra.GetSession().Query(queryUpdateExpires, at.Expires, at.AccessToken).
		Exec(); err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) *errors.RestErr {
	if err := cassandra.GetSession().Query(queryCreateToken, at.AccessToken, at.UserId, at.ClientId, at.Expires).
		Exec(); err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}
