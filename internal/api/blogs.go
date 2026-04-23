package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/itshivams/studex-cli/internal/config"
)

type BlogAuthor struct {
	Username string `json:"username"`
	UserPic  string `json:"userpic"`
	Status   string `json:"status"`
	FullName string `json:"full_name"`
}

type BlogItem struct {
	ID            string     `json:"_id"`
	Title         string     `json:"title"`
	Slug          string     `json:"slug"`
	Excerpt       string     `json:"excerpt"`
	Tags          []string   `json:"tags"`
	CoverImage    string     `json:"coverImage"`
	Visibility    string     `json:"visibility"`
	Author        BlogAuthor `json:"author"`
	LikesCount    int        `json:"likes_count"`
	CommentsCount int        `json:"comments_count"`
	CreatedAt     string     `json:"createdAt"`
	UpdatedAt     string     `json:"updatedAt"`
	Views         int        `json:"views"`
	ReadTime      int        `json:"readTime"`
	Markdown      string     `json:"markdown,omitempty"`
	LikedBy       []string   `json:"likedBy,omitempty"`
}

type BlogViewResponse struct {
	Ok   bool     `json:"ok"`
	Blog BlogItem `json:"blog"`
}

type BlogListResponse struct {
	Count int        `json:"count"`
	Items []BlogItem `json:"items"`
}

func GetBlogs() (*BlogListResponse, error) {
	url := fmt.Sprintf("%s/blog/list", BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	token := config.GetToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch blogs with status: %d", resp.StatusCode)
	}

	var res BlogListResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func GetBlogView(slug string) (*BlogViewResponse, error) {
	url := fmt.Sprintf("%s/blog/view?slug=%s", BaseURL, slug)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	token := config.GetToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch blog view with status: %d", resp.StatusCode)
	}

	var res BlogViewResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
