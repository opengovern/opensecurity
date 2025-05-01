package git

import (
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// CommitDetails holds commit metadata.
type CommitDetails struct {
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
}

// ComparisonResultGrouped holds the JSON output structure with files grouped under parent folders.
type ComparisonResultGrouped struct {
	CommitDetails   [2]CommitDetails `json:"commit_details"`
	ModifiedFiles   map[string]bool  `json:"modified_files"`
	CreatedFiles    map[string]bool  `json:"created_files"`
	DeletedFiles    map[string]bool  `json:"deleted_files"`
	UnmodifiedFiles map[string]bool  `json:"unmodified_files"`
	ChangedFolders  []string         `json:"changed_folders"`
}

// CloneAndShowDiff clones and returns the diff comparison
func CloneAndShowDiff(logger *zap.Logger, repoURL, targetDir, firstCommitSHA, secondCommitSHA string) (*ComparisonResultGrouped, error) {
	var repo *git.Repository
	var err error

	if isRemoteURI(repoURL) {
		logger.Info("Cloning repository", zap.String("source", repoURL), zap.String("destination", targetDir))
		_, err = CloneRepository(logger, repoURL, targetDir)
		if err != nil {
			logger.Error("Failed to clone repository using fetch package", zap.Error(err))
			os.Exit(1)
		}

		// Open the cloned repository
		repo, err = git.PlainOpen(targetDir)
		if err != nil {
			logger.Error("Failed to open cloned repository", zap.Error(err))
			return nil, err
		}
	} else {
		// Open the local repository
		repo, err = git.PlainOpen(repoURL)
		if err != nil {
			logger.Error("Failed to open local repository", zap.Error(err))
			return nil, err
		}
	}

	firstCommit, err := getCommit(repo, firstCommitSHA)
	if err != nil {
		logger.Error("Failed to retrieve first_commit", zap.Error(err))
		return nil, err
	}

	if secondCommitSHA == "" {
		secondCommitSHA, err = getLastCommitSHA(repo)
		if err != nil {
			logger.Error("Failed to get the latest commit", zap.Error(err))
			return nil, err
		}
	}

	secondCommit, err := getCommit(repo, secondCommitSHA)
	if err != nil {
		logger.Error("Failed to retrieve second_commit", zap.Error(err))
		return nil, err
	}

	// Ensure first_commit is older
	if firstCommit.Committer.When.After(secondCommit.Committer.When) {
		firstCommit, secondCommit = secondCommit, firstCommit
	}

	result := compareCommits(logger, firstCommit, secondCommit)
	result.CommitDetails[0] = CommitDetails{Hash: firstCommit.Hash.String(), Timestamp: firstCommit.Committer.When}
	result.CommitDetails[1] = CommitDetails{Hash: secondCommit.Hash.String(), Timestamp: secondCommit.Committer.When}

	logger.Info("Diff operation completed successfully")
	return &result, nil
}

// isRemoteURI checks if the provided path is a remote URI.
func isRemoteURI(uri string) bool {
	return strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "git@")
}

// getCommit retrieves the commit object by its SHA-1 hash.
func getCommit(repo *git.Repository, sha string) (*object.Commit, error) {
	commitHash := plumbing.NewHash(sha)
	return repo.CommitObject(commitHash)
}

// getLastCommitSHA gets the SHA of the latest commit on the default branch.
func getLastCommitSHA(repo *git.Repository) (string, error) {
	ref, err := repo.Head()
	if err != nil {
		return "", err
	}
	return ref.Hash().String(), nil
}

// compareCommits analyzes the differences between two commits.
func compareCommits(logger *zap.Logger, firstCommit, secondCommit *object.Commit) ComparisonResultGrouped {
	tree1, err := firstCommit.Tree()
	if err != nil {
		logger.Error("Failed to get tree for first_commit")
		return ComparisonResultGrouped{}
	}
	tree2, err := secondCommit.Tree()
	if err != nil {
		logger.Error("Failed to get tree for second_commit", zap.Error(err))
		return ComparisonResultGrouped{}
	}

	patch, err := tree1.Patch(tree2)
	if err != nil {
		logger.Error("Failed to create patch between commits", zap.Error(err))
		return ComparisonResultGrouped{}
	}

	modifiedFiles := make(map[string]bool)
	createdFiles := make(map[string]bool)
	deletedFiles := make(map[string]bool)
	unmodifiedFiles := make(map[string]bool)
	changedFoldersSet := make(map[string]bool)

	changedFilesSet := make(map[string]bool)
	filePatches := patch.FilePatches()
	for _, filePatch := range filePatches {
		from, to := filePatch.Files()
		if from == nil && to != nil {
			// File was created
			createdPath := to.Path()
			parentFolder := getParentFolder(createdPath)
			createdFiles[parentFolder+"/"+filepath.Base(createdPath)] = true
			changedFilesSet[createdPath] = true
			changedFoldersSet[parentFolder] = true
		} else if from != nil && to == nil {
			// File was deleted
			deletedPath := from.Path()
			parentFolder := getParentFolder(deletedPath)
			deletedFiles[parentFolder+"/"+filepath.Base(deletedPath)] = true
			changedFilesSet[deletedPath] = true
			changedFoldersSet[parentFolder] = true
		} else if from != nil && to != nil {
			// File was modified
			modifiedPath := from.Path()
			parentFolder := getParentFolder(modifiedPath)
			modifiedFiles[parentFolder+"/"+filepath.Base(modifiedPath)] = true
			changedFilesSet[modifiedPath] = true
			changedFoldersSet[parentFolder] = true
		}
	}

	// Collect all files from both trees
	allFilesSet := make(map[string]bool)
	err = tree1.Files().ForEach(func(f *object.File) error {
		allFilesSet[f.Name] = true
		return nil
	})
	if err != nil {
		logger.Error("Failed to iterate over files in first_commit tree", zap.Error(err))
		return ComparisonResultGrouped{}
	}
	err = tree2.Files().ForEach(func(f *object.File) error {
		allFilesSet[f.Name] = true
		return nil
	})
	if err != nil {
		logger.Error("Failed to iterate over files in second_commit tree", zap.Error(err))
		return ComparisonResultGrouped{}
	}

	// Identify unmodified files
	for file := range allFilesSet {
		if !changedFilesSet[file] {
			parentFolder := getParentFolder(file)
			unmodifiedFiles[parentFolder+"/"+filepath.Base(file)] = true
		}
	}

	// Convert changedFoldersSet to a slice
	changedFolders := make([]string, 0, len(changedFoldersSet))
	for folder := range changedFoldersSet {
		changedFolders = append(changedFolders, folder)
	}

	return ComparisonResultGrouped{
		ModifiedFiles:   modifiedFiles,
		CreatedFiles:    createdFiles,
		DeletedFiles:    deletedFiles,
		UnmodifiedFiles: unmodifiedFiles,
		ChangedFolders:  changedFolders,
	}
}

// getParentFolder returns the parent folder of a given file path.
// If the file is at the root, it returns "root".
func getParentFolder(filePath string) string {
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "/" {
		return "root"
	}
	return filepath.ToSlash(dir)
}
