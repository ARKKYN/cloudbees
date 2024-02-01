package services

import (
	d "cloudbees/dao"
	e "cloudbees/errors"
	"cloudbees/genproto/posts"
	m "cloudbees/models"
	"context"
	"errors"
	"regexp"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostsService struct {
	posts.UnimplementedBlogServiceServer
	postsDao *d.PostDAO
}

func NewPostsService(dao *d.PostDAO) *PostsService {
	return &PostsService{
		postsDao: dao,
	}
}

func convertToPostResponse(post *m.Post) *posts.PostResponse {
	return &posts.PostResponse{
		PostId:          post.PostId,
		Title:           post.Title,
		Content:         post.Content,
		Author:          post.Author,
		PublicationDate: post.PublicationDate,
		Tags:            post.Tags,
	}
}

func validateDateFormat(date string) error {
	// Define the regular expression pattern for the date format
	pattern := `^\d{2}-(0[1-9]|1[0-2])-\d{4}$`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Check if the date matches the pattern
	if !regex.MatchString(date) {
		return errors.New("invalid date format")
	}

	return nil
}

func ValidateCreatePostRequest(in *posts.CreatePostRequest) error {

	if in.PostId == 0 {
		return e.PostIdMissingError
	}
	if in.Title == "" {
		return e.TitleMissingError
	}
	if in.Content == "" {
		return e.ContentMissingError
	}
	if in.Author == "" {
		return e.AuthorMissingError
	}

	if in.PublicationDate == "" {
		return e.PublicationDateMissingError
	}

	if validateDateFormat(in.PublicationDate) != nil {
		return e.InvalidPublicationDateError
	}

	if len(in.Tags) == 0 {
		return e.TagsMissingError
	}
	return nil
}

// CleanTags cleans and validates tags.
func CleanTags(tags []string) []string {
	cleanedTags := make([]string, 0, len(tags))
	for _, tag := range tags {
		cleanedTag := strings.TrimSpace(tag)
		if cleanedTag != "" {
			cleanedTags = append(cleanedTags, cleanedTag)
		}
	}
	return cleanedTags
}

func (s *PostsService) CreatePost(ctx context.Context, in *posts.CreatePostRequest) (*posts.PostResponse, error) {
	// Validate input fields
	if err := ValidateCreatePostRequest(in); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Clean and validate tags
	cleanedTags := CleanTags(in.Tags)
	if len(cleanedTags) == 0 {
		return nil, status.Error(codes.InvalidArgument, e.TagsMissingError.Error())
	}

	// Create post object
	post := &m.Post{
		PostId:          in.PostId,
		Title:           in.Title,
		Content:         in.Content,
		Author:          in.Author,
		PublicationDate: in.PublicationDate,
		Tags:            cleanedTags,
	}

	// Persist post to database
	err := s.postsDao.Create(post)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert post to response format and return
	return convertToPostResponse(post), nil
}

func (s *PostsService) GetPost(ctx context.Context, in *posts.GetPostRequest) (*posts.PostResponse, error) {
	post, err := s.postsDao.Read(in.PostId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return convertToPostResponse(post), nil
}

func updatePostFields(post *m.Post, in *posts.UpdatePostRequest) {
	if in.Title != "" {
		post.Title = in.Title
	}
	if in.Content != "" {
		post.Content = in.Content
	}
	if in.Author != "" {
		post.Author = in.Author
	}
	if in.PublicationDate != "" {
		post.PublicationDate = in.PublicationDate
	}
	if len(in.Tags) != 0 {
		post.Tags = CleanTags(in.Tags)
	}
}

func (s *PostsService) UpdatePost(ctx context.Context, in *posts.UpdatePostRequest) (*posts.PostResponse, error) {
	post, err := s.postsDao.Read(in.PostId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if validateDateFormat(in.PublicationDate) != nil {
		return nil, status.Error(codes.InvalidArgument, e.InvalidPublicationDateError.Error())
	}
	updatePostFields(post, in)

	err = s.postsDao.Update(post)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertToPostResponse(post), nil
}

func (s *PostsService) DeletePost(ctx context.Context, in *posts.DeletePostRequest) (*posts.DeletePostResponse, error) {
	err := s.postsDao.Delete(in.PostId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &posts.DeletePostResponse{
		Message: "Post deleted successfully",
	}, nil
}
