package git

import (
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// CloneRepository clones a Git repository to the specified directory.
// If targetDir is empty, it defaults to the repo's name or "cloned-repo-<timestamp>".
func CloneRepository(logger *zap.Logger, gitRepoURI, targetDir string) (*git.Repository, error) {
	parsedURL, err := url.Parse(gitRepoURI)
	if err != nil {
		logger.Error("Invalid URL")
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	// If targetDir is not provided, derive it from the repo name
	if targetDir == "" {
		repoName := extractRepoName(parsedURL.Path)
		if repoName == "" {
			// Default to "cloned-repo-<date/timestamp>" if the repo name can't be derived
			targetDir = fmt.Sprintf("cloned-repo-%s", time.Now().Format("20060102-150405"))
			logger.Info("Target directory not provided or determined from the URL.", zap.String("default", targetDir))
		} else {
			targetDir = repoName
		}
	}

	// Determine the platform based on the host
	switch parsedURL.Host {
	case "github.com":
		pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
		if len(pathSegments) >= 3 && (pathSegments[2] == "tree" || pathSegments[2] == "releases") {
			return nil, handleGitHubSpecificURL(logger, gitRepoURI, targetDir, pathSegments)
		}
	case "dev.azure.com", "ssh.dev.azure.com":
		return nil, handleAzureDevOpsURL(logger, gitRepoURI, targetDir, parsedURL)
	case "gitlab.com":
		return nil, handleGitLabURL(logger, gitRepoURI, targetDir, parsedURL)
	}

	// Default cloning behavior for standard Git URIs
	cloneOptions := &git.CloneOptions{
		URL:   gitRepoURI,
		Depth: 1,
	}

	if strings.HasPrefix(gitRepoURI, "https://") {
		cloneOptions.Auth = getHTTPAuth()
	} else if strings.HasPrefix(gitRepoURI, "git@") {
		cloneOptions.Auth, err = getSSHAuth(logger)
		if err != nil {
			return nil, err
		}
	} else {
		logger.Error("Unsupported Git repository URI format")
		return nil, fmt.Errorf("unsupported Git repository URI format")
	}

	res, err := git.PlainClone(targetDir, false, cloneOptions)
	if err != nil {
		logger.Error("Failed to clone repository", zap.Error(err))
		return nil, err
	}

	logger.Info("Successfully cloned repository", zap.String("source", gitRepoURI), zap.String("target", targetDir))
	return res, nil
}

// extractRepoName extracts the repository name from the URL path.
func extractRepoName(urlPath string) string {
	segments := strings.Split(strings.Trim(urlPath, "/"), "/")
	if len(segments) >= 2 {
		return strings.TrimSuffix(segments[len(segments)-1], ".git") // Handles URLs with or without the .git suffix
	}
	return ""
}

// handleGitHubSpecificURL handles cloning based on GitHub's tree (branch) or releases (tag) URLs.
func handleGitHubSpecificURL(logger *zap.Logger, gitRepoURI, targetDir string, pathSegments []string) error {
	// Construct the base Git URL (e.g., https://github.com/user/repo.git)
	baseURL := fmt.Sprintf("https://github.com/%s/%s.git", pathSegments[0], pathSegments[1])

	// Clone the repository
	cloneOptions := &git.CloneOptions{
		URL:   baseURL,
		Auth:  getHTTPAuth(),
		Depth: 1,
	}

	var err error
	if strings.HasPrefix(gitRepoURI, "git@") {
		cloneOptions.Auth, err = getSSHAuth(logger)
		if err != nil {
			return err
		}
	}

	repo, err := git.PlainClone(targetDir, false, cloneOptions)
	if err != nil {
		logger.Error("Failed to clone repository", zap.Error(err))
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	// Determine whether we're dealing with a branch or a tag
	if pathSegments[2] == "tree" && len(pathSegments) > 3 {
		// Handle branch
		return checkoutBranch(logger, repo, pathSegments[3])
	} else if pathSegments[2] == "releases" && len(pathSegments) > 4 && pathSegments[3] == "tag" {
		// Handle tag
		return checkoutTag(logger, repo, pathSegments[4])
	}

	logger.Error("Unsupported GitHub URL structure")
	return fmt.Errorf("unsupported GitHub URL structure")
}

// handleAzureDevOpsURL handles cloning based on Azure DevOps URLs.
func handleAzureDevOpsURL(logger *zap.Logger, gitRepoURI, targetDir string, parsedURL *url.URL) error {
	// Example Azure DevOps HTTPS URL: https://dev.azure.com/{organization}/{project}/_git/{repo}?version=GB{branch} or GT{tag}
	// Example Azure DevOps SSH URL: git@ssh.dev.azure.com:v3/{organization}/{project}/{repo}

	// Parse query parameters for branch or tag
	query := parsedURL.Query()
	version := query.Get("version") // e.g., GBbranch or GTtag

	var branchName, tagName string
	if strings.HasPrefix(version, "GB") {
		branchName = strings.TrimPrefix(version, "GB")
	} else if strings.HasPrefix(version, "GT") {
		tagName = strings.TrimPrefix(version, "GT")
	}

	// Extract organization, project, repo from path
	pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	var organization, project, repo string

	if parsedURL.Host == "dev.azure.com" {
		// HTTPS URL: /{organization}/{project}/_git/{repo}
		if len(pathSegments) >= 4 && pathSegments[2] == "_git" {
			organization = pathSegments[0]
			project = pathSegments[1]
			repo = pathSegments[3]
		} else {
			logger.Error("Invalid Azure DevOps HTTPS URL structure")
			return fmt.Errorf("invalid Azure DevOps HTTPS URL structure")
		}
	} else if parsedURL.Host == "ssh.dev.azure.com" {
		// SSH URL: v3/{organization}/{project}/{repo}
		sshPath := strings.TrimPrefix(parsedURL.Path, "/v3/")
		sshSegments := strings.Split(strings.Trim(sshPath, "/"), "/")
		if len(sshSegments) >= 3 {
			organization = sshSegments[0]
			project = sshSegments[1]
			repo = sshSegments[2]
		} else {
			logger.Error("Invalid Azure DevOps SSH URL structure")
			return fmt.Errorf("invalid Azure DevOps SSH URL structure")
		}
	} else {
		logger.Error("Unsupported Azure DevOps host")
		return fmt.Errorf("unsupported Azure DevOps host: %s", parsedURL.Host)
	}

	// Construct the base Git URL
	var baseURL string
	if parsedURL.Host == "dev.azure.com" {
		baseURL = fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", organization, project, repo)
	} else { // ssh.dev.azure.com
		baseURL = fmt.Sprintf("git@ssh.dev.azure.com:v3/%s/%s/%s", organization, project, repo)
	}

	// Clone the repository
	cloneOptions := &git.CloneOptions{
		URL:   baseURL,
		Auth:  getHTTPAuth(),
		Depth: 1,
	}

	var err error
	if strings.HasPrefix(gitRepoURI, "git@") {
		cloneOptions.Auth, err = getSSHAuth(logger)
		if err != nil {
			logger.Error("Failed to clone repository", zap.Error(err))
		}
	}

	repoObj, err := git.PlainClone(targetDir, false, cloneOptions)
	if err != nil {
		logger.Error("Failed to clone Azure DevOps repository", zap.Error(err))
		return fmt.Errorf("failed to clone Azure DevOps repository: %v", err)
	}

	// Checkout the specified branch or tag
	if branchName != "" {
		return checkoutBranch(logger, repoObj, branchName)
	} else if tagName != "" {
		return checkoutTag(logger, repoObj, tagName)
	}

	// If no specific branch or tag, default to default branch
	logger.Info("No specific branch or tag specified. Using default branch.")
	return nil
}

// handleGitLabURL handles cloning based on GitLab URLs.
func handleGitLabURL(logger *zap.Logger, gitRepoURI, targetDir string, parsedURL *url.URL) error {
	// Example GitLab HTTPS URL: https://gitlab.com/{group}/{project}.git?ref={branch or tag}
	// Example GitLab SSH URL: git@gitlab.com:{group}/{project}.git

	// Parse query parameters for branch or tag
	query := parsedURL.Query()
	ref := query.Get("ref") // e.g., branch or tag name

	var branchName, tagName string
	if ref != "" {
		// GitLab doesn't differentiate between branch and tag in the ref parameter
		// We'll attempt to checkout as a branch first, then as a tag
		branchName = ref
		tagName = ref
	}

	// Extract group and project from path
	pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	var group, project string

	if parsedURL.Host == "gitlab.com" {
		if len(pathSegments) >= 2 {
			group = pathSegments[0]
			project = strings.TrimSuffix(pathSegments[1], ".git")
		} else {
			logger.Error("Invalid GitLab HTTPS URL structure")
			return fmt.Errorf("invalid GitLab HTTPS URL structure")
		}
	} else if parsedURL.Host == "ssh.gitlab.com" {
		// SSH URL: {group}/{project}.git
		sshPath := strings.TrimPrefix(parsedURL.Path, "/")
		sshSegments := strings.Split(strings.TrimSuffix(sshPath, ".git"), "/")
		if len(sshSegments) >= 2 {
			group = sshSegments[0]
			project = sshSegments[1]
		} else {
			logger.Error("Invalid GitLab SSH URL structure")
			return fmt.Errorf("invalid GitLab SSH URL structure")
		}
	} else {
		logger.Error("Unsupported GitLab host")
		return fmt.Errorf("unsupported GitLab host: %s", parsedURL.Host)
	}

	// Construct the base Git URL
	var baseURL string
	if parsedURL.Host == "gitlab.com" {
		baseURL = fmt.Sprintf("https://gitlab.com/%s/%s.git", group, project)
	} else { // ssh.gitlab.com
		baseURL = fmt.Sprintf("git@gitlab.com:%s/%s.git", group, project)
	}

	// Clone the repository
	cloneOptions := &git.CloneOptions{
		URL:   baseURL,
		Auth:  getHTTPAuth(),
		Depth: 1,
	}

	var err error
	if strings.HasPrefix(gitRepoURI, "git@") {
		cloneOptions.Auth, err = getSSHAuth(logger)
		if err != nil {
			logger.Error("Failed to clone GitLab repository", zap.Error(err))
			return fmt.Errorf("failed to clone GitLab repository: %v", err)
		}
	}

	repoObj, err := git.PlainClone(targetDir, false, cloneOptions)
	if err != nil {
		logger.Error("Failed to clone GitLab repository", zap.Error(err))
		return fmt.Errorf("failed to clone GitLab repository: %v", err)
	}

	// Checkout the specified branch or tag
	if branchName != "" {
		err := checkoutBranch(logger, repoObj, branchName)
		if err != nil {
			// Attempt to checkout as a tag if branch checkout fails
			return checkoutTag(logger, repoObj, tagName)
		}
		return nil
	}

	// If no specific ref, default to default branch
	logger.Info("No specific branch or tag specified. Using default branch.")
	return nil
}

// checkoutBranch checks out the specified branch in the cloned repository.
func checkoutBranch(logger *zap.Logger, repo *git.Repository, branchName string) error {
	worktree, err := repo.Worktree()
	if err != nil {
		logger.Error("Failed to get worktree", zap.Error(err))
		return fmt.Errorf("failed to get worktree: %v", err)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
	})
	if err != nil {
		logger.Error("Failed to checkout branch", zap.String("branch", branchName), zap.Error(err))
		return fmt.Errorf("failed to checkout branch %s: %v", branchName, err)
	}
	logger.Info("Checked out branch", zap.String("branch", branchName))
	return nil
}

// checkoutTag checks out the specified tag in the cloned repository.
func checkoutTag(logger *zap.Logger, repo *git.Repository, tagName string) error {
	tagRef, err := repo.Tag(tagName)
	if err != nil {
		logger.Error("Failed to find tag", zap.String("tagName", tagName), zap.Error(err))
		return fmt.Errorf("failed to find tag %s: %v", tagName, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		logger.Error("Failed to get worktree", zap.Error(err))
		return fmt.Errorf("failed to get worktree: %v", err)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Hash: tagRef.Hash(),
	})
	if err != nil {
		logger.Error("Failed to checkout tag", zap.String("tagName", tagName), zap.Error(err))
		return fmt.Errorf("failed to checkout tag %s: %v", tagName, err)
	}
	logger.Info("Checked out tag", zap.String("tagName", tagName))
	return nil
}

// getSSHAuth handles SSH authentication for git@ URIs.
func getSSHAuth(logger *zap.Logger) (transport.AuthMethod, error) {
	sshAuth, err := gitssh.NewPublicKeysFromFile("git", filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"), "")
	if err != nil {
		logger.Error("Failed to create SSH auth method")
		return nil, err
	}

	// Load known_hosts file to avoid host key verification failure
	knownHostsFile := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	hostKeyCallback, err := knownhosts.New(knownHostsFile)
	if err != nil {
		logger.Error("Failed to load known_hosts file", zap.Error(err))
		return nil, err
	}

	sshAuth.HostKeyCallback = hostKeyCallback
	return sshAuth, nil
}

// getHTTPAuth handles HTTP authentication (for private HTTPS repos).
func getHTTPAuth() *http.BasicAuth {
	username := os.Getenv("GIT_USERNAME")
	password := os.Getenv("GIT_PASSWORD")

	if username == "" || password == "" {
		return nil
	}

	return &http.BasicAuth{
		Username: username,
		Password: password,
	}
}
