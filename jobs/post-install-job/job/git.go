package job

import (
	"os"

	"encoding/json"
	// "strings"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/opensecurity/jobs/post-install-job/config"
	git2 "github.com/opengovern/opensecurity/jobs/post-install-job/job/git"
	"github.com/opengovern/opensecurity/services/core/client"
	"github.com/opengovern/opensecurity/services/core/db/models"
	"go.uber.org/zap"
)

func GitClone(conf config.MigratorConfig, logger *zap.Logger) (string, error) {
	gitConfig := GitConfig{
		AnalyticsGitURL:         conf.AnalyticsGitURL,
		ControlEnrichmentGitURL: conf.ControlEnrichmentGitURL,
		githubToken:             conf.GithubToken,
	}

	coreClient := client.NewCoreServiceClient(conf.Core.BaseURL)
	value, err := coreClient.GetConfigMetadata(&httpclient.Context{
		UserRole: api.AdminRole,
	}, models.MetadataKeyPlatformConfigurationGitURL)

	if err == nil && len(value.GetValue().(string)) > 0 {
		gitConfig.AnalyticsGitURL = value.GetValue().(string)
	} else if err != nil {
		logger.Error("failed to get analytics git url from metadata", zap.Error(err))
	}

	logger.Info("using git repo", zap.String("url", gitConfig.AnalyticsGitURL))

	refs := make([]string, 0, 2)
	os.RemoveAll(config.ConfigzGitPath)

	res, err := git2.CloneRepository(logger, gitConfig.AnalyticsGitURL, config.ConfigzGitPath)
	if err != nil {
		logger.Error("failed to clone repository", zap.Error(err))
		return "", err
	}

	logger.Info("finished fetching configz data")

	ref, err := res.Head()
	if err != nil {
		logger.Error("failed to get head", zap.Error(err))
		return "", err
	}
	refs = append(refs, ref.Hash().String())

	// logger.Info("using git repo for enrichmentor", zap.String("url", gitConfig.ControlEnrichmentGitURL))

	//os.RemoveAll(config.ControlEnrichmentGitPath)
	//
	//res, err = git2.CloneRepository(logger, gitConfig.ControlEnrichmentGitURL, config.ControlEnrichmentGitPath)
	//if err != nil {
	//	logger.Error("failed to clone repository", zap.Error(err))
	//	return "", err
	//}
	//ref, err = res.Head()
	//if err != nil {
	//	logger.Error("failed to get head", zap.Error(err))
	//	return "", err
	//}
	//refs = append(refs, ref.Hash().String())

	// refsJson, err := json.Marshal(refs)
	// if err != nil {
	// 	logger.Error("failed to marshal refs", zap.Error(err))
	// 	return "", err
	// }
	refsJson, err := json.Marshal(refs)
	if err != nil {
		logger.Error("failed to marshal refs", zap.Error(err))
		return "", err
	}
	return string(refsJson), nil
}
