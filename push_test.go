package git

import (
	"log"
	"testing"
	"time"
)

func Test_Push_ToRemote(t *testing.T) {
	repo := createBareTestRepo(t)
	repo2 := createTestRepo(t)

	remote, err := repo2.CreateRemote("test_push", repo.Path())
	checkFatal(t, err)

	index, err := repo2.Index()
	checkFatal(t, err)

	index.AddByPath("README")

	err = index.Write()
	checkFatal(t, err)

	newTreeId, err := index.WriteTree()
	checkFatal(t, err)

	tree, err := repo2.LookupTree(newTreeId)
	checkFatal(t, err)

	sig := &Signature{Name: "Rand Om Hacker", Email: "random@hacker.com", When: time.Now()}
	// this should cause master branch to be created if it does not already exist
	_, err = repo2.CreateCommit("HEAD", sig, sig, "message", tree)
	checkFatal(t, err)

	push, err := remote.NewPush()
	checkFatal(t, err)

	err = push.AddRefspec("refs/heads/master")
	checkFatal(t, err)

	err = push.Finish()
	checkFatal(t, err)

	err = push.StatusForeach(func(ref string, msg string) int {
		log.Printf("%s -> %s", ref, msg)
		return 0
	})
	checkFatal(t, err)

	if !push.UnpackOk() {
		t.Fatalf("unable to unpack")
	}

	defer remote.Free()
	defer repo.Free()
}
