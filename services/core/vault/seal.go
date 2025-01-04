package vault

import (
	"context"
	"fmt"
	"time"

	"github.com/opengovern/og-util/pkg/vault"
	"github.com/opengovern/opencomply/services/core/config"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type SealHandler struct {
	logger           *zap.Logger
	cfg              config.Config
	vaultSealHandler *vault.HashiCorpVaultSealHandler
	kubeClientset    *kubernetes.Clientset
}

func NewSealHandler(ctx context.Context, logger *zap.Logger, cfg config.Config) (*SealHandler, error) {
	hashiCorpVaultSealHandler, err := vault.NewHashiCorpVaultSealHandler(ctx, logger, cfg.Vault.HashiCorp)
	if err != nil {
		logger.Error("new hashicorp vaultClient seal handler", zap.Error(err))
		return nil, fmt.Errorf("new hashicorp vaultClient seal handler: %w", err)
	}

	kuberConfig, err := rest.InClusterConfig()
	if err != nil {
		logger.Error("failed to get kubernetes config", zap.Error(err))
		return nil, fmt.Errorf("failed to get kubernetes config: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(kuberConfig)
	if err != nil {
		logger.Error("failed to create kubernetes clientset", zap.Error(err))
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	return &SealHandler{
		logger:           logger,
		cfg:              cfg,
		vaultSealHandler: hashiCorpVaultSealHandler,
		kubeClientset:    clientset,
	}, nil
}

const (
	secretName = "vault-unseal-keys"
)

func (s *SealHandler) initVault(ctx context.Context) bool {
	initRes, err := s.vaultSealHandler.TryInit(ctx)
	if err != nil {
		s.logger.Fatal("failed to init vault", zap.Error(err))
	}
	if initRes == nil {
		s.logger.Info("vault already initialized")
		return false
	} else {
		keysSecret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: s.cfg.OpengovernanceNamespace,
			},
			StringData: make(map[string]string),
		}

		for i, key := range initRes.Keys {
			keysSecret.StringData[fmt.Sprintf("key-%d", i)] = key
		}
		keysSecret.StringData["root-token"] = initRes.RootToken

		s.logger.Info("vault initialized creating unseal keys secret")
		_, err = s.kubeClientset.CoreV1().Secrets(s.cfg.OpengovernanceNamespace).Create(ctx, &keysSecret, metav1.CreateOptions{})
		if err != nil && !k8serrors.IsAlreadyExists(err) {
			s.logger.Fatal("failed to create vault unseal keys secret", zap.Error(err), zap.Strings("keys", initRes.Keys))
		} else if k8serrors.IsAlreadyExists(err) && len(initRes.Keys) > 0 {
			s.logger.Info("vault unseal keys secret already exists, updating")
			_, err := s.kubeClientset.CoreV1().Secrets(s.cfg.OpengovernanceNamespace).Update(ctx, &keysSecret, metav1.UpdateOptions{})
			if err != nil {
				s.logger.Fatal("failed to update vault unseal keys secret", zap.Error(err), zap.Strings("keys", initRes.Keys))
			}
		}
		return true
	}
}

func (s *SealHandler) unsealChecker(ctx context.Context, initKuber bool, unsealed chan<- struct{}) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("unsealChecker panic", zap.Any("recover", r))
			go s.unsealChecker(ctx, initKuber, unsealed)
		}
	}()

	keysSecret, err := s.kubeClientset.CoreV1().Secrets(s.cfg.OpengovernanceNamespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		s.logger.Error("failed to get vault unseal keys secret", zap.Error(err))
		return
	}

	keys := make([]string, 0, len(keysSecret.Data))
	for k, v := range keysSecret.Data {
		if k == "root-token" {
			continue
		}
		keys = append(keys, string(v))
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	err = s.vaultSealHandler.TryUnseal(ctx, keys)
	if err != nil {
		s.logger.Error("failed to unseal vault", zap.Error(err))
	}
	if unsealed != nil && err == nil {
		rootToken := keysSecret.Data["root-token"]
		err = s.vaultSealHandler.SetupKuberAuth(ctx, string(rootToken))
		if err != nil {
			s.logger.Error("failed to setup kubernetes auth", zap.Error(err))
		}
		unsealed <- struct{}{}
		close(unsealed)
		unsealed = nil
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err = s.vaultSealHandler.TryUnseal(ctx, keys)
			if err != nil {
				s.logger.Error("failed to unseal vault", zap.Error(err))
				continue
			}
			if unsealed != nil {
				rootToken := keysSecret.Data["root-token"]
				err = s.vaultSealHandler.SetupKuberAuth(ctx, string(rootToken))
				if err != nil {
					s.logger.Error("failed to setup kubernetes auth", zap.Error(err))
				}
				unsealed <- struct{}{}
				close(unsealed)
				unsealed = nil
			}
		}
	}
}

func (s *SealHandler) Start(ctx context.Context) {
	isNewInit := s.initVault(ctx)
	unsealChan := make(chan struct{})
	go s.unsealChecker(ctx, isNewInit, unsealChan)
	// block until vault is unsealed
	<-unsealChan
	s.logger.Info("vault unsealed")
}
