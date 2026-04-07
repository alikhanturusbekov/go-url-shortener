package grpcserver

import (
	"context"
	"github.com/alikhanturusbekov/go-url-shortener/api/proto"
	"net/http"

	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/authorization"
	appError "github.com/alikhanturusbekov/go-url-shortener/pkg/error"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ShortenerServer implements the gRPC ShortenerService
type ShortenerServer struct {
	shortenerpb.UnimplementedShortenerServiceServer
	service *service.URLService
}

// NewShortenerServer creates a new gRPC server instance
func NewShortenerServer(service *service.URLService) *ShortenerServer {
	return &ShortenerServer{
		service: service,
	}
}

// ShortenURL returns a shortened URL for the provided URL
func (s *ShortenerServer) ShortenURL(
	ctx context.Context,
	req *shortenerpb.URLShortenRequest,
) (*shortenerpb.URLShortenResponse, error) {
	userID, _ := authorization.UserIDFromContext(ctx)

	shortURL, appErr := s.service.ShortenURL(req.GetUrl(), userID)
	if appErr != nil && shortURL == "" {
		return nil, toGRPCError(appErr)
	}

	resp := &shortenerpb.URLShortenResponse{}
	resp.SetResult(shortURL)

	if appErr != nil {
		return resp, toGRPCError(appErr)
	}

	return resp, nil
}

// ExpandURL returns original URL for the provided short URL
func (s *ShortenerServer) ExpandURL(
	_ context.Context,
	req *shortenerpb.URLExpandRequest,
) (*shortenerpb.URLExpandResponse, error) {
	originalURL, appErr := s.service.ResolveShortURL(req.GetId())
	if appErr != nil {
		return nil, toGRPCError(appErr)
	}

	resp := &shortenerpb.URLExpandResponse{}
	resp.SetResult(originalURL)
	return resp, nil
}

// ListUserURLs returns the list of user's URLs
func (s *ShortenerServer) ListUserURLs(
	ctx context.Context,
	_ *emptypb.Empty,
) (*shortenerpb.UserURLsResponse, error) {
	userID, ok := authorization.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authorization required")
	}

	userURLs, appErr := s.service.GetUserURLs(userID)
	if appErr != nil {
		return nil, toGRPCError(appErr)
	}

	items := make([]*shortenerpb.URLData, 0, len(userURLs))
	for _, item := range userURLs {
		pbItem := &shortenerpb.URLData{}
		pbItem.SetShortUrl(item.ShortURL)
		pbItem.SetOriginalUrl(item.OriginalURL)
		items = append(items, pbItem)
	}

	resp := &shortenerpb.UserURLsResponse{}
	resp.SetUrl(items)

	return resp, nil
}

// toGRPCError converts internal HTTP errors to gRPC status errors
func toGRPCError(err *appError.HTTPError) error {
	switch err.Code {
	case http.StatusBadRequest:
		return status.Error(codes.InvalidArgument, err.GetFullMessage())
	case http.StatusUnauthorized:
		return status.Error(codes.Unauthenticated, err.GetFullMessage())
	case http.StatusForbidden:
		return status.Error(codes.PermissionDenied, err.GetFullMessage())
	case http.StatusNotFound:
		return status.Error(codes.NotFound, err.GetFullMessage())
	case http.StatusConflict:
		return status.Error(codes.AlreadyExists, err.GetFullMessage())
	case http.StatusGone:
		return status.Error(codes.FailedPrecondition, err.GetFullMessage())
	default:
		return status.Error(codes.Internal, err.GetFullMessage())
	}
}
