package git

/*
#include <git2.h>
#include <git2/errors.h>

extern void _go_git_set_strarray_n(git_strarray *array, char *str, size_t n);
extern char *_go_git_get_strarray_n(git_strarray *array, size_t n);
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type Remote struct {
	ptr *C.git_remote
}

func newRemoteFromC(c *C.git_remote) *Remote {
	remote := &Remote{ptr: c}
	runtime.SetFinalizer(remote, (*Remote).Free)
	return remote
}

func (r *Remote) Free() {
	C.git_remote_free(r.ptr)
	runtime.SetFinalizer(r, nil)
}

func (r *Repository) ListRemoteNames() ([]string, error) {
	var cnames C.git_strarray
	ret := C.git_remote_list(&cnames, r.ptr)
	if ret < 0 {
		return nil, LastError()
	}
	defer C.git_strarray_free(&cnames)

	names := make([]string, cnames.count)
	for i := 0; i < int(cnames.count); i++ {
		names[i] = C.GoString(C._go_git_get_strarray_n(&cnames, C.size_t(i)))
	}

	return names, nil
}

func (r *Repository) CreateRemote(name string, url string) (*Remote, error) {
	var c *C.git_remote

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))

	runtime.LockOSThread()
	runtime.UnlockOSThread()

	ret := C.git_remote_create(&c, r.ptr, cname, curl)
	if ret < 0 {
		return nil, LastError()
	}

	return newRemoteFromC(c), nil

}

func (r *Repository) CreateRemoteInMemory(name string, url string) (*Remote, error) {
	var c *C.git_remote

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))

	runtime.LockOSThread()
	runtime.UnlockOSThread()

	ret := C.git_remote_create_inmemory(&c, r.ptr, cname, curl)
	if ret < 0 {
		return nil, LastError()
	}

	return newRemoteFromC(c), nil

}

func (r *Repository) LoadRemote(name string) (*Remote, error) {
	var c *C.git_remote

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if C.git_remote_load(&c, r.ptr, cname) < 0 {
		return nil, LastError()
	}

	return newRemoteFromC(c), nil
}

func (r *Remote) Name() string {
	return C.GoString(C.git_remote_name(r.ptr))
}

func (r *Remote) URL() string {
	return C.GoString(C.git_remote_url(r.ptr))
}

type GitDirection C.git_direction

const (
	GitDirectionFetch GitDirection = C.GIT_DIRECTION_FETCH
	GitDirectionPush               = C.GIT_DIRECTION_PUSH
)

func (r *Remote) Connect(direction GitDirection) error {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_remote_connect(r.ptr, C.git_direction(direction))
	if ret < 0 {
		return LastError()
	}
	return nil
}

func (r *Remote) Disconnect() {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	C.git_remote_disconnect(r.ptr)
}

func (r *Remote) UpdateTips(sig *Signature, msg string) error {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_remote_update_tips(r.ptr, nil, nil)
	if ret < 0 {
		return LastError()
	}
	return nil
}

func (r *Remote) Download() error {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_remote_download(r.ptr)

	if ret < 0 {
		return LastError()
	}
	return nil

}
