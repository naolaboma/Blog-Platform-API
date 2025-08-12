package usecase

import (
	"Blog-API/internal/domain"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type blogUseCase struct {
	blogRepo domain.BlogRepository
	userRepo domain.UserRepository
	cache    domain.Cache
}

func NewBlogUseCase(
	blogRepo domain.BlogRepository,
	userRepo domain.UserRepository,
	cache domain.Cache,
) domain.BlogUseCase {
	return &blogUseCase{
		blogRepo: blogRepo,
		userRepo: userRepo,
		cache:    cache,
	}
}

func (uc *blogUseCase) CreateBlog(blog *domain.Blog, authorID primitive.ObjectID) error {
	author, err := uc.userRepo.GetByID(authorID)
	if err != nil {
		return errors.New("author not found")
	}
	//server generated fields
	blog.ID = primitive.NewObjectID()
	blog.AuthorID = authorID
	blog.AuthorUsername = author.Username
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	blog.Comments = []domain.Comment{}
	blog.Likes = []string{}
	blog.Dislikes = []string{}
	blog.ViewCount = 0
	blog.LikeCount = 0
	blog.CommentCount = 0

	if err := uc.blogRepo.Create(blog); err != nil {
		return err
	}

	//invalidate caches that list multiple blogs
	go uc.invalidateBlogListCaches()

	return nil
}

func (uc *blogUseCase) GetBlog(id primitive.ObjectID) (*domain.Blog, error) {
	ctx := context.Background()
	key := fmt.Sprintf("blog:%s", id.Hex())
	var blog domain.Blog

	if err := uc.cache.Get(ctx, key, &blog); err == nil {
		log.Println("CACHE HIT: GetBlog")
		go uc.blogRepo.IncrementViewCount(id)
		return &blog, nil
	}
	log.Println("CACHE MISS: GetBlog")
	dbBlog, err := uc.blogRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("blog not found")
	}
	go uc.cache.Set(ctx, key, dbBlog, 10*time.Minute)
	return dbBlog, nil
}

func (uc *blogUseCase) GetAllBlogs(page, limit int, sort string) ([]*domain.Blog, int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("blogs:all:page=%d:limit=%d:sort=%s", page, limit, sort)
	var cachedResult struct {
		Blogs []*domain.Blog
		Total int64
	}
	if err := uc.cache.Get(ctx, key, &cachedResult); err == nil {
		log.Println("CACHE HIT: GetAllBlogs")
		return cachedResult.Blogs, cachedResult.Total, nil
	}
	log.Println("CACHE MISS: GetAllBlogs")
	blogs, total, err := uc.blogRepo.GetAll(page, limit, sort)
	if err != nil {
		return nil, 0, err
	}
	go uc.cache.Set(ctx, key, struct {
		Blogs []*domain.Blog
		Total int64
	}{
		blogs, total}, 5*time.Minute)
	return blogs, total, nil
}

func (uc *blogUseCase) UpdateBlog(id primitive.ObjectID, blogUpdate *domain.Blog, userID primitive.ObjectID, userRole string) (*domain.Blog, error) {
	originalBlog, err := uc.blogRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("blog not found")
	}
	if originalBlog.AuthorID != userID && userRole != domain.RoleAdmin {
		return nil, errors.New("forbidden: you are not authorized to update this post")
	}

	originalBlog.Title = blogUpdate.Title
	originalBlog.Content = blogUpdate.Content
	originalBlog.Tags = blogUpdate.Tags
	originalBlog.UpdatedAt = time.Now()

	if err := uc.blogRepo.Update(originalBlog); err != nil {
		return nil, err
	}
	//invalidate cache for this specific blog and for all ListenAndServe
	go uc.cache.Delete(context.Background(), fmt.Sprintf("blog:%s", id.Hex()))
	go uc.invalidateBlogListCaches()

	return originalBlog, nil
}

func (uc *blogUseCase) DeleteBlog(id primitive.ObjectID, userID primitive.ObjectID, userRole string) error {
	blog, err := uc.blogRepo.GetByID(id)
	if err != nil {
		return errors.New("blog not found")
	}
	if blog.AuthorID != userID && userRole != domain.RoleAdmin {
		return errors.New("forbidden: you are not authorized to delete this post")
	}
	if err := uc.blogRepo.Delete(id); err != nil {
		return err
	}
	go uc.cache.Delete(context.Background(), fmt.Sprintf("blog:%s", id.Hex()))
	go uc.invalidateBlogListCaches()
	return nil
}

func (uc *blogUseCase) AddComment(blogID primitive.ObjectID, comment *domain.Comment) error {
	author, err := uc.userRepo.GetByID(comment.AuthorID)
	if err != nil {
		return errors.New("comment author not found")
	}
	comment.ID = primitive.NewObjectID()
	comment.AuthorUsername = author.Username
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	if err := uc.blogRepo.AddComment(blogID, comment); err != nil {
		return err
	}

	go uc.cache.Delete(context.Background(), fmt.Sprintf("blog:%s", blogID.Hex()))
	return nil
}

func (uc *blogUseCase) LikeBlog(blogID primitive.ObjectID, userID string) error {
	blog, err := uc.blogRepo.GetByID(blogID)
	if err != nil {
		return errors.New("blog not found")
	}
	isLiked := containsString(blog.Likes, userID)
	isDisliked := containsString(blog.Dislikes, userID)

	var err2 error
	if isLiked {
		err2 = uc.blogRepo.RemoveLike(blogID, userID)
	} else if isDisliked {
		if err := uc.blogRepo.RemoveDislike(blogID, userID); err != nil {
			return err
		}
		err2 = uc.blogRepo.AddLike(blogID, userID)
	} else {
		err2 = uc.blogRepo.AddLike(blogID, userID)
	}

	// Invalidate cache for this specific blog
	if err2 == nil {
		go uc.cache.Delete(context.Background(), fmt.Sprintf("blog:%s", blogID.Hex()))
	}

	return err2
}

func (uc *blogUseCase) DislikeBlog(blogID primitive.ObjectID, userID string) error {
	blog, err := uc.blogRepo.GetByID(blogID)
	if err != nil {
		return errors.New("blog not found")
	}
	isLiked := containsString(blog.Likes, userID)
	isDisliked := containsString(blog.Dislikes, userID)

	var err2 error
	if isDisliked {
		err2 = uc.blogRepo.RemoveDislike(blogID, userID)
	} else if isLiked {
		if err := uc.blogRepo.RemoveLike(blogID, userID); err != nil {
			return err
		}
		err2 = uc.blogRepo.AddDislike(blogID, userID)
	} else {
		err2 = uc.blogRepo.AddDislike(blogID, userID)
	}

	// Invalidate cache for this specific blog
	if err2 == nil {
		go uc.cache.Delete(context.Background(), fmt.Sprintf("blog:%s", blogID.Hex()))
	}

	return err2
}

func (uc *blogUseCase) SearchBlogsByTitle(title string, page, limit int) ([]*domain.Blog, int64, error) {

	ctx := context.Background()
	key := fmt.Sprintf("blogs:search:title=%s:page=%d:limit=%d", title, page, limit)
	var cachedResult struct {
		Blogs []*domain.Blog
		Total int64
	}
	if err := uc.cache.Get(ctx, key, &cachedResult); err == nil {
		log.Println("CACHE HIT: SearchBlogsByTitle")
		return cachedResult.Blogs, cachedResult.Total, nil
	}
	log.Println("CACHE MISS: SearchBlogsByTitle")
	blogs, total, err := uc.blogRepo.SearchByTitle(title, page, limit)
	if err != nil {
		return nil, 0, err
	}
	go uc.cache.Set(ctx, key, struct {
		Blogs []*domain.Blog
		Total int64
	}{
		blogs, total}, 5*time.Minute)
	return blogs, total, nil
}

func (uc *blogUseCase) SearchBlogsByAuthor(author string, page, limit int) ([]*domain.Blog, int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("blogs:search:author=%s:page=%d:limit=%d", author, page, limit)
	var cachedResult struct {
		Blogs []*domain.Blog
		Total int64
	}
	if err := uc.cache.Get(ctx, key, &cachedResult); err == nil {
		log.Printf("CACHE HIT: SearchBlogsByAuthor")
		return cachedResult.Blogs, cachedResult.Total, nil
	}
	log.Println("CACHE MISS: SearchBlogsByAuthor")
	blogs, total, err := uc.blogRepo.SearchByAuthor(author, page, limit)
	if err != nil {
		return nil, 0, err
	}
	go uc.cache.Set(ctx, key, struct {
		Blogs []*domain.Blog
		Total int64
	}{blogs, total}, 5*time.Minute)
	return blogs, total, nil
}

func (uc *blogUseCase) FilterBlogsByTags(tags []string, page, limit int) ([]*domain.Blog, int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("blogs:filter:tags=%s:page=%d:limit=%d", strings.Join(tags, ","), page, limit)
	var cachedResult struct {
		Blogs []*domain.Blog
		Total int64
	}
	if err := uc.cache.Get(ctx, key, &cachedResult); err == nil {
		log.Println("CACHE HIT: FilterBlogsByTags")
		return cachedResult.Blogs, cachedResult.Total, nil
	}
	log.Println("CACHE MISS: FilterBlogsByTags")
	blogs, total, err := uc.blogRepo.FilterByTags(tags, page, limit)
	if err != nil {
		return nil, 0, err
	}
	go uc.cache.Set(ctx, key, struct {
		Blogs []*domain.Blog
		Total int64
	}{blogs, total}, 5*time.Minute)
	return blogs, total, nil
}

func (uc *blogUseCase) FilterBlogsByDate(startDate, endDate time.Time, page, limit int) ([]*domain.Blog, int64, error) {
	return uc.blogRepo.FilterByDate(startDate, endDate, page, limit)
}

func (uc *blogUseCase) GetPopularBlogs(limit int) ([]*domain.Blog, error) {
	ctx := context.Background()
	key := fmt.Sprintf("blogs:popular:limit=%d", limit)
	var blogs []*domain.Blog
	if err := uc.cache.Get(ctx, key, &blogs); err == nil {
		log.Println("CACHE HIT: GetPopularBlogs")
		return blogs, nil
	}
	log.Println("CACHE MISS: GetPopularBlogs")
	dbBlogs, err := uc.blogRepo.GetPopular(limit)
	if err != nil {
		return nil, err
	}
	go uc.cache.Set(ctx, key, dbBlogs, 15*time.Minute)
	return dbBlogs, nil
}

func (uc *blogUseCase) DeleteComment(blogID, commentID primitive.ObjectID, userID primitive.ObjectID) error {
	blog, err := uc.blogRepo.GetByID(blogID)
	if err != nil {
		return errors.New("blog not found")
	}
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	var commentAuthorID primitive.ObjectID
	found := false
	for _, c := range blog.Comments {
		if c.ID == commentID {
			commentAuthorID = c.AuthorID
			found = true
			break
		}
	}

	if !found {
		return errors.New("comment not found")
	}

	isCommentAuthor := commentAuthorID == userID
	isBlogAuthor := blog.AuthorID == userID
	isAdmin := user.Role == domain.RoleAdmin

	if !isCommentAuthor && !isBlogAuthor && !isAdmin {
		return errors.New("forbideen: you are not authorized to delete this comment")
	}

	if err := uc.blogRepo.DeleteComment(blogID, commentID); err != nil {
		return err
	}
	go uc.cache.Delete(context.Background(), fmt.Sprintf("blog:%s", blogID.Hex()))
	return nil
}

func (uc *blogUseCase) UpdateComment(blogID, commentID primitive.ObjectID, content string, userID primitive.ObjectID) error {
	// Only the original comment author can update their comment.
	blog, err := uc.blogRepo.GetByID(blogID)
	if err != nil {
		return errors.New("blog not found")
	}
	var commentAuthorID primitive.ObjectID
	found := false
	for _, c := range blog.Comments {
		if c.ID == commentID {
			commentAuthorID = c.AuthorID
			found = true
			break
		}
	}
	if !found {
		return errors.New("comment not found")
	}
	if commentAuthorID != userID {
		return errors.New("forbidden: you are not the author of this comment")
	}
	if err := uc.blogRepo.UpdateComment(blogID, commentID, content); err != nil {
		return err
	}
	go uc.cache.Delete(context.Background(), fmt.Sprintf("blog:%s", blogID.Hex()))
	return nil
}

// helper functions
func (uc *blogUseCase) invalidateBlogListCaches() {
	ctx := context.Background()
	uc.cache.DeleteByPattern(ctx, "blogs:all:*")
	uc.cache.DeleteByPattern(ctx, "blogs:search:*")
	uc.cache.DeleteByPattern(ctx, "blogs:filter:*")
	uc.cache.DeleteByPattern(ctx, "blogs:popular:*")
	log.Println("CACHE INVALIDATION: Cleared blog list caches")
}
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
