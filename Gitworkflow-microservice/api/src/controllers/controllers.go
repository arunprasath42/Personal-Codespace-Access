package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// CloneRepository - Clones a repository from the given URL to the given path
func CloneRepository(c *gin.Context) {
	repoURL := c.PostForm("repo_url")
	clonePath := c.PostForm("clone_path")

	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repository cloned successfully"})
}

// CreateBranch - Creates a new branch in the given repository
func CreateBranch(c *gin.Context) {
	repoPath := c.PostForm("repo_path")
	branchName := c.PostForm("branch_name")

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a new branch
	branchRef := plumbing.NewBranchReferenceName("refs/heads/" + branchName)

	fmt.Println(branchRef)
	err = repo.CreateBranch(&config.Branch{
		Name:   branchName,
		Remote: "origin",
		Merge:  plumbing.ReferenceName("refs/heads/master"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch created successfully"})
}

// CommitChanges - Commits the changes to the given repository
func CommitChanges(c *gin.Context) {
	repoPath := c.PostForm("repo_path")
	filePath := c.PostForm("file_path")
	fileContent := c.PostForm("file_content")
	commitMessage := c.PostForm("commit_message")

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	w, err := repo.Worktree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	file := filepath.Join(repoPath, filePath)
	err = os.WriteFile(file, []byte(fileContent), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = w.Add(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Commit the changes
	_, err = w.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "arunprasath42",
			Email: "arunsundar4298@gmail.com",
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Changes committed successfully"})
}

// ViewHistory - Returns the commit history of the given repository
/*func ViewHistory(c *gin.Context) {
	repoPath := c.PostForm("repo_path")

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ref, err := repo.Head()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var commits []string
	err = commitIter.ForEach(func(commit *object.Commit) error {
		commits = append(commits, commit.Hash.String())
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"commits": commits})
}*/

func ViewHistory1(c *gin.Context) {
	repoPath := c.PostForm("repo_path")
	branchName := c.PostForm("branch_name")

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	branchRef, err := repo.Reference(plumbing.ReferenceName("refs/heads/"+branchName), true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commitIter, err := repo.Log(&git.LogOptions{From: branchRef.Hash()})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var commits []gin.H
	err = commitIter.ForEach(func(commit *object.Commit) error {
		changes, err := commit.Stats()
		if err != nil {
			return err
		}

		changedFiles := make([]string, 0)
		for _, change := range changes {
			changedFiles = append(changedFiles, change.Name)
		}

		commits = append(commits, gin.H{
			"commit_id":     commit.Hash.String(),
			"commit_time":   commit.Author.When,
			"changed_files": changedFiles,
		})
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"commit_history": commits})
}

func ViewHistory(c *gin.Context) {
	repoURL := viper.GetString("github.url")
	branchName := c.PostForm("branch_name")

	memStorage := memory.NewStorage()
	r, err := git.Clone(memStorage, nil, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	branchRef, err := r.Reference(plumbing.ReferenceName("refs/heads/"+branchName), true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commitIter, err := r.Log(&git.LogOptions{From: branchRef.Hash()})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var commits []gin.H
	err = commitIter.ForEach(func(commit *object.Commit) error {
		changes, err := commit.Stats()
		if err != nil {
			return err
		}

		changedFiles := make([]string, 0)
		for _, change := range changes {
			changedFiles = append(changedFiles, change.Name)
		}

		commits = append(commits, gin.H{
			"commit_id":     commit.Hash.String(),
			"commit_time":   commit.Author.When,
			"changed_files": changedFiles,
		})
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"commit_history": commits})
}

// ShowDiff1 - to be deleted
func ShowDiff1(c *gin.Context) {
	repoPath := c.PostForm("repo_path")
	commitHash1 := c.PostForm("commit_hash_1")
	commitHash2 := c.PostForm("commit_hash_2")

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commitObj1, err := repo.CommitObject(plumbing.NewHash(commitHash1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commitObj2, err := repo.CommitObject(plumbing.NewHash(commitHash2))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	patch, err := commitObj2.Patch(commitObj1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	diffStats := patch.Stats()
	diffChanges := make([]string, len(diffStats))

	for i, stats := range diffStats {
		diffChanges[i] = fmt.Sprintf("File: %s\nAdded: %d, Deleted: %d", stats.Name, stats.Addition, stats.Deletion)
	}

	c.JSON(http.StatusOK, gin.H{"diff_changes": diffChanges})
}

// ShowDiff - Show the diff between two commits
func ShowDiff(c *gin.Context) {
	repoURL := c.PostForm("repo_url")
	commitHash1 := c.PostForm("commit_hash_1")
	commitHash2 := c.PostForm("commit_hash_2")

	memStorage := memory.NewStorage()

	r, err := git.Clone(memStorage, nil, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commitObj1, err := r.CommitObject(plumbing.NewHash(commitHash1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commitObj2, err := r.CommitObject(plumbing.NewHash(commitHash2))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	patch, err := commitObj2.Patch(commitObj1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	diffStats := patch.Stats()
	diffChanges := make([]gin.H, len(diffStats))

	for i, stats := range diffStats {
		diffChanges[i] = gin.H{
			"file":     stats.Name,
			"added":    stats.Addition,
			"deleted":  stats.Deletion,
			"modified": stats.String(),
		}
	}

	c.JSON(http.StatusOK, gin.H{"diff_changes": diffChanges})
}
