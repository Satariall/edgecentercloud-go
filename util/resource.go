package util

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

var (
	errResourceNotDeleted = errors.New("could not delete the resource")
	errGetResourceInfo    = errors.New("error when retrieving resource information")
)

type RetrieveResourceFunc[T any] func(ctx context.Context, id string) (*T, *edgecloud.Response, error)

func ResourceIsDeleted[T any](ctx context.Context, retrieveResourceFunc RetrieveResourceFunc[T], id string) error {
	_, resp, err := retrieveResourceFunc(ctx, id)
	if err == nil {
		return errResourceNotDeleted
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	return errGetResourceInfo
}

type GetResourceFunc[T any] func(ctx context.Context, id string) (*T, *edgecloud.Response, error)

func ResourceIsExist[T any](ctx context.Context, getResourceFunc GetResourceFunc[T], id string) (bool, error) {
	_, resp, _ := getResourceFunc(ctx, id)

	return HandleStatusCode(resp.StatusCode)
}

func HandleStatusCode(statusCode int) (bool, error) {
	switch statusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound, http.StatusForbidden:
		return false, nil
	default:
		return false, fmt.Errorf("%w, status code: %d", errGetResourceInfo, statusCode)
	}
}
