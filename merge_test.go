package git

import (
	"testing"
)

func Test_Merge_With_Self(t *testing.T) {

	repo := createTestRepo(t)
	seedTestRepo(t, repo)

	master, err := repo.LookupReference("refs/heads/master")
	checkFatal(t, err)

	mergeHead, err := repo.MergeHeadFromRef(master)
	checkFatal(t, err)

	options := DefaultMergeOptions()
	mergeHeads := make([]*MergeHead, 1)
	mergeHeads[0] = mergeHead
	results, err := repo.Merge(mergeHeads, options)
	checkFatal(t, err)

	if !results.IsUpToDate() {
		t.Fatal("Expected up to date")
	}
}

func Test_Merge_Fast_Forward(t *testing.T) {

	repo := createTestRepo(t)
	commitId, _ := seedTestRepo(t, repo)

	commit, err := repo.LookupCommit(commitId)
	checkFatal(t, err)

	_, err = repo.CreateBranch("branched", commit, false, testSig(t), "")
	checkFatal(t, err)

	updateTestRepo(t, commitId, repo, []byte("this is an update"))

	master, err := repo.LookupReference("refs/heads/master")
	checkFatal(t, err)

	mergeHead, err := repo.MergeHeadFromId(master.Target())
	checkFatal(t, err)

	head, err := repo.LookupReference("HEAD")
	checkFatal(t, err)

	_, err = head.SetSymbolicTarget("refs/heads/branched", testSig(t), "")
	checkFatal(t, err)

	opts := &CheckoutOpts{Strategy: CheckoutForce}
	err = repo.Checkout(opts)
	checkFatal(t, err)

	options := DefaultMergeOptions()
	mergeHeads := make([]*MergeHead, 1)
	mergeHeads[0] = mergeHead
	results, err := repo.Merge(mergeHeads, options)
	checkFatal(t, err)

	if !results.IsFastForward() {
		t.Fatal("Expected fast forward")
	}
}
