package main

import (
	"cloudbees/dao"
	e "cloudbees/errors"
	"cloudbees/genproto/posts"
	"cloudbees/services"
	"context"
	"fmt"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func setupServer() (*grpc.Server, *net.Listener) {
	server := grpc.NewServer()
	postsDao := dao.NewPostDAO()
	postsService := services.NewPostsService(postsDao)
	posts.RegisterBlogServiceServer(server, postsService)
	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(fmt.Errorf("failed to listen: %v", err))
	}
	go func() {
		if err := server.Serve(listen); err != nil {
			panic(fmt.Errorf("failed to serve gRPC server: %v", err))
		}
	}()
	return server, &listen
}

func setupClient(serverAddress string) posts.BlogServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Errorf("failed to dial server: %v", err))
	}
	return posts.NewBlogServiceClient(conn)
}

func TestCreatePostIntegration(t *testing.T) {
	server, listen := setupServer()

	serverAddress := "localhost:8080"
	client := setupClient(serverAddress)

	testCases := []struct {
		name     string
		request  *posts.CreatePostRequest
		expected error
	}{
		{
			name: "Successful creation of a post",
			request: &posts.CreatePostRequest{
				PostId:          1,
				Title:           "Test Post",
				Content:         "Test Content",
				Author:          "Test Author",
				PublicationDate: "01-01-2024",
				Tags:            []string{"test", "integration"},
			},
			expected: nil,
		},
		{
			name: "Validation error scenario: empty title",
			request: &posts.CreatePostRequest{
				PostId:  1,
				Title:   "",
				Content: "Test Content",
				Author:  "Test Author",
				Tags:    []string{"test", "integration"},
			},
			expected: status.Error(codes.InvalidArgument, e.TitleMissingError.Error()),
		},
		{
			name: "Validation error scenario: empty content",
			request: &posts.CreatePostRequest{
				PostId:          1,
				Title:           "Test Post",
				Content:         "",
				Author:          "Test Author",
				PublicationDate: "01-01-2024",
				Tags:            []string{"test", "integration"},
			},
			expected: status.Error(codes.InvalidArgument, e.ContentMissingError.Error()),
		},
		{
			name: "Validation error scenario: empty author",
			request: &posts.CreatePostRequest{
				PostId:          1,
				Title:           "Test Post",
				Content:         "Test Content",
				Author:          "",
				PublicationDate: "01-01-2024",
				Tags:            []string{"test", "integration"},
			},
			expected: status.Error(codes.InvalidArgument, e.AuthorMissingError.Error()),
		},
		{
			name: "Validation error scenario: empty publication date",
			request: &posts.CreatePostRequest{
				PostId:          1,
				Title:           "Test Post",
				Content:         "Test Content",
				Author:          "Test Author",
				PublicationDate: "",
				Tags:            []string{"test", "integration"},
			},
			expected: status.Error(codes.InvalidArgument, e.PublicationDateMissingError.Error()),
		},
		{
			name: "Validation error scenario: Invalid publication date format",
			request: &posts.CreatePostRequest{
				PostId:          1,
				Title:           "Test Post",
				Content:         "Test Content",
				Author:          "Test Author",
				PublicationDate: "2024-01-01",
				Tags:            []string{"test", "integration"},
			},
			expected: status.Error(codes.InvalidArgument, e.InvalidPublicationDateError.Error()),
		},
		{
			name: "Validation error scenario: empty tags",
			request: &posts.CreatePostRequest{
				PostId:          1,
				Title:           "Test Post",
				Content:         "Test Content",
				Author:          "Test Author",
				PublicationDate: "01-01-2024",
				Tags:            []string{},
			},
			expected: status.Error(codes.InvalidArgument, e.TagsMissingError.Error()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			createdPost, err := client.CreatePost(context.Background(), tc.request)

			if fmt.Sprint(err) != fmt.Sprint(tc.expected) {
				t.Fatalf(" expected %v, got  %v", err, tc.expected)
			}
			if tc.expected == nil && createdPost == nil {
				t.Fatal("expected non-nil post response")
			}
		})
	}
	server.Stop()
	(*listen).Close()
}

func TestReadPostIntegration(t *testing.T) {
	server, listen := setupServer()

	serverAddress := "localhost:8080"
	client := setupClient(serverAddress)

	testCases := []struct {
		name     string
		request  *posts.GetPostRequest
		expected error
	}{
		{
			name: "Successful read of a post",
			request: &posts.GetPostRequest{
				PostId: 1,
			},
			expected: nil,
		},
		{
			name: "Error reading  post",
			request: &posts.GetPostRequest{
				PostId: 2,
			},
			expected: status.Error(codes.NotFound, e.EnitityNotFoundError.Error()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			createdPost, err := client.GetPost(context.Background(), tc.request)
			if fmt.Sprint(err) != fmt.Sprint(tc.expected) {
				t.Fatalf(" %v,  %v", err, tc.expected)
			}
			if tc.expected == nil && createdPost == nil {
				t.Fatal("expected non-nil post response")
			}
		})
	}
	server.Stop()
	(*listen).Close()

}

func TestUpdatePostIntegration(t *testing.T) {
	server, listen := setupServer()

	serverAddress := "localhost:8080"
	client := setupClient(serverAddress)

	updateRequest := &posts.UpdatePostRequest{
		PostId:          1,
		Title:           "Updated Test Post",
		Content:         "Updated Test Content",
		Author:          "Updated Test Author",
		PublicationDate: "01-01-2025",
		Tags:            []string{"updated", "integration"},
	}
	updatedPost, err := client.UpdatePost(context.Background(), updateRequest)
	if err != nil {
		t.Fatalf("failed to update post: %v", err)
	}

	expectedPost := &posts.PostResponse{
		PostId:          1,
		Title:           "Updated Test Post",
		Content:         "Updated Test Content",
		Author:          "Updated Test Author",
		PublicationDate: "01-01-2025",
		Tags:            []string{"updated", "integration"},
	}

	if !isPostEqual(updatedPost, expectedPost) {
		t.Fatalf("updated post does not match expected post")
	}

	testCases := []struct {
		name     string
		request  *posts.UpdatePostRequest
		expected error
	}{
		{
			name: "Error scenario: post not found",
			request: &posts.UpdatePostRequest{
				PostId:          2,
				Title:           "Updated Test Post",
				Content:         "Updated Test Content",
				Author:          "Updated Test Author",
				PublicationDate: "01-01-2025",
				Tags:            []string{"updated", "integration"},
			},
			expected: status.Error(codes.NotFound, e.EnitityNotFoundError.Error()),
		},
		{
			name: "Error scenario: invalid publication date format",
			request: &posts.UpdatePostRequest{
				PostId:          1,
				Title:           "Updated Test Post",
				Content:         "Updated Test Content",
				Author:          "Updated Test Author",
				PublicationDate: "2025-01-01",
				Tags:            []string{"updated", "integration"},
			},
			expected: status.Error(codes.InvalidArgument, e.InvalidPublicationDateError.Error()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.UpdatePost(context.Background(), tc.request)
			if err == nil {
				t.Fatalf("expected error: %v", tc.expected)
			}
			if err.Error() != tc.expected.Error() {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}

	server.Stop()
	(*listen).Close()
}

func isPostEqual(post1, post2 *posts.PostResponse) bool {
	if post1.PostId != post2.PostId {
		return false
	}
	if post1.Title != post2.Title {
		return false
	}
	if post1.Content != post2.Content {
		return false
	}
	if post1.Author != post2.Author {
		return false
	}
	if post1.PublicationDate != post2.PublicationDate {
		return false
	}
	if len(post1.Tags) != len(post2.Tags) {
		return false
	}
	for i := range post1.Tags {
		if post1.Tags[i] != post2.Tags[i] {
			return false
		}
	}
	return true
}

func TestDeletePostIntegration(t *testing.T) {
	server, listen := setupServer()

	serverAddress := "localhost:8080"
	client := setupClient(serverAddress)

	deleteRequest := &posts.DeletePostRequest{
		PostId: 1,
	}
	_, err := client.DeletePost(context.Background(), deleteRequest)
	if err != nil {
		t.Fatalf("failed to delete post: %v", err)
	}

	getRequest := &posts.GetPostRequest{
		PostId: 1,
	}
	_, err = client.GetPost(context.Background(), getRequest)
	expectedErr := status.Error(codes.NotFound, e.EnitityNotFoundError.Error())
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected error: %v, got: %v", expectedErr, err)
	}

	server.Stop()
	(*listen).Close()
}
