package git

import (
	"os"

	"github.com/abr-ooo/VulEQ/models"
	"github.com/abr-ooo/VulEQ/log"
	"github.com/abr-ooo/VulEQ/configs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func CloneCode(p models.Project, commitHash string) error {

	path := configs.Cfg.Git.ClonePath + p.Name

	CloneOptions := &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "oauth2",
			Password: p.GitToken,
		},
		URL:           p.GitURL,
		ReferenceName: plumbing.NewBranchReferenceName(p.GitBranch),
		SingleBranch:  true,
		Progress:      os.Stdout,
	}

	repo, err := git.PlainClone(path, false, CloneOptions)

	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"error":         err.Error(),
			"path":          path,
			"is_bare":       "false",
			"clone_options": CloneOptions,
		}).Debug("Git. Clone code failed!")
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		log.Log.Debug("Git. Get work tree from repo failed with this error : \n", err)
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commitHash),
	})
	if err != nil {
		log.Log.Debug("Git. Checkout on commit hash failed with this error : \n", err)
		return err
	}

	log.Log.Debug("Git. clone code successful :)")
	return nil
}
