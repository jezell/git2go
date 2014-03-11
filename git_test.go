package git

import (
	"io/ioutil"
	"testing"
	"time"
)

func createTestRepo(t *testing.T) *Repository {
	// figure out where we can create the test repo
	path, err := ioutil.TempDir("", "git2go")
	checkFatal(t, err)
	repo, err := InitRepository(path, false)
	checkFatal(t, err)

	tmpfile := "README"
	err = ioutil.WriteFile(path+"/"+tmpfile, []byte("foo\n"), 0644)
	checkFatal(t, err)

	return repo
}

func createBareTestRepo(t *testing.T) *Repository {
	// figure out where we can create the test repo
	path, err := ioutil.TempDir("", "git2go")
	checkFatal(t, err)
	repo, err := InitRepository(path, true)
	checkFatal(t, err)

	return repo
}

func testSig(t *testing.T) *Signature {

	loc, err := time.LoadLocation("Europe/Berlin")
	checkFatal(t, err)

	return &Signature{
		Name:  "Rand Om Hacker",
		Email: "random@hacker.com",
		When:  time.Date(2013, 03, 06, 14, 30, 0, 0, loc),
	}
}

func seedTestRepo(t *testing.T, repo *Repository) (*Oid, *Oid) {
	loc, err := time.LoadLocation("Europe/Berlin")
	checkFatal(t, err)
	sig := &Signature{
		Name:  "Rand Om Hacker",
		Email: "random@hacker.com",
		When:  time.Date(2013, 03, 06, 14, 30, 0, 0, loc),
	}

	idx, err := repo.Index()
	checkFatal(t, err)
	err = idx.AddByPath("README")
	checkFatal(t, err)
	treeId, err := idx.WriteTree()
	checkFatal(t, err)

	message := "This is a commit\n"
	tree, err := repo.LookupTree(treeId)
	checkFatal(t, err)
	commitId, err := repo.CreateCommit("HEAD", sig, sig, message, tree)
	checkFatal(t, err)

	return commitId, treeId
}

func updateTestRepo(t *testing.T, lastCommitId *Oid, repo *Repository, data []byte) (*Oid, *Oid) {

	tmpfile := "README"
	err := ioutil.WriteFile(repo.Path()+"/../"+tmpfile, data, 0644)
	checkFatal(t, err)

	loc, err := time.LoadLocation("Europe/Berlin")
	checkFatal(t, err)
	sig := &Signature{
		Name:  "Rand Om Hacker",
		Email: "random@hacker.com",
		When:  time.Date(2013, 03, 06, 14, 30, 0, 0, loc),
	}

	idx, err := repo.Index()
	checkFatal(t, err)
	err = idx.AddByPath("README")
	checkFatal(t, err)
	treeId, err := idx.WriteTree()
	checkFatal(t, err)

	message := "This is a commit\n"
	tree, err := repo.LookupTree(treeId)
	checkFatal(t, err)

	lastCommit, err := repo.LookupCommit(lastCommitId)
	checkFatal(t, err)

	commitId, err := repo.CreateCommit("HEAD", sig, sig, message, tree, lastCommit)
	checkFatal(t, err)

	return commitId, treeId
}
