package git

import (
	"testing"
)

func TestCreateRemote(t *testing.T) {
	repo := createTestRepo(t)
	_, err := repo.CreateRemote("test", "http://github.com")

	if err != nil {
		t.Fatal(err)
	}

	names, err := repo.ListRemoteNames()
	if err != nil {
		t.Fatal(err)
	}

	if len(names) != 1 {
		t.Fatal("expected 1 name")
	}

	if names[0] != "test" {
		t.Fatal("expected name[0] == \"test\"")
	}
}

func TestDownloadRemote(t *testing.T) {
	repo := createTestRepo(t)
	remote, err := repo.CreateRemote("test", "http://github.com/jezell/zlibber")

	if err != nil {
		t.Fatal(err)
	}

	err = remote.Connect(GitDirectionFetch)

	if err != nil {
		t.Fatal(err)
	}

	err = remote.Download()

	if err != nil {
		t.Fatal(err)
	}

	remote.Disconnect()
	remote.UpdateTips(nil, "")

}
