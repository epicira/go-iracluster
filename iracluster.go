package iracluster

import (
	"unsafe"
)

/*
#cgo CPPFLAGS: -I${SRCDIR}/include
#cgo darwin LDFLAGS: -lstdc++ -L/opt/homebrew/lib -lnats.3.8.2 -lbrotlicommon.1.1.0 -lbrotlidec.1.1.0 -lbrotlienc.1.1.0 -lboost_system-mt -lboost_coroutine-mt -lboost_stacktrace_basic-mt -lboost_thread-mt -lboost_timer-mt -lboost_date_time-mt -lboost_filesystem-mt -lssl.3 -lcrypto.3 -lsqlite3 -lcurl -lz -lfmt -lspdlog -L/usr/local/lib -liracommon -liracluster
#cgo linux LDFLAGS: -L/usr/lib/x86_64-linux-gnu/ -lm -lstdc++ -lnats -lbrotlicommon -lbrotlidec -lbrotlienc -lssl3 -lcrypto -lsqlite3 -lcurl -lz -lfmt -lspdlog -L/usr/local/lib -lboost_system-mt-x64 -lboost_coroutine-mt-x64 -lboost_stacktrace_backtrace-mt-x64 -lboost_thread-mt-x64 -lboost_timer-mt-x64 -lboost_date_time-mt-x64 -lboost_filesystem-mt-x64 -liracommon -liracluster
#include "include/ira_cluster_callback.h"
#include "include/ira_cluster_c.h"
*/
import "C"

type IraCluster struct {
	AppName   string
	ClusterID string
	cPtr      unsafe.Pointer
}

func New(appName string, clusterId string, cb func(string)) *IraCluster {
	arena := &Arena{}
	defer arena.free()

	iraCluster := &IraCluster{
		AppName:   appName,
		ClusterID: clusterId,
		cPtr:      C.NewIraCluster(arena.CString(appName), arena.CString(clusterId), arena.CString("1.0.0")),
	}
	Callbacks = append(Callbacks, cb)
	return iraCluster
}

func (ic *IraCluster) JoinCluster() bool {
	return bool(C.joinCluster(ic.cPtr))
}

func (ic *IraCluster) IsSenior() bool {
	return bool(C.isSenior(ic.cPtr))
}

func (ic *IraCluster) PeerCount() int {
	return int(C.peerCount(ic.cPtr))
}

type PrivacyLevel int

const (
	None PrivacyLevel = iota
	Local
	Shared
	Private
)

func (p *PrivacyLevel) String() string {
	return [...]string{"none", "local", "shared", "private"}[*p]
}

type Privacy struct {
	level      PrivacyLevel
	publicKey  string
	privateKey string
}

type TDB struct {
	ClusterID          string
	Privacy            *Privacy
	DBName             string
	InitSQL            string
	IndexSQL           string
	PublishChanges     bool
	iraClusterInstance *IraCluster
}

func NewTDB(ic *IraCluster, dbName, initSql, indexSql string, publish bool, privacy *Privacy) *TDB {
	if privacy == nil {
		privacy = &Privacy{
			level:      None,
			publicKey:  "",
			privateKey: "",
		}
	}
	return &TDB{
		ClusterID:          ic.ClusterID,
		Privacy:            privacy,
		DBName:             dbName,
		InitSQL:            initSql,
		IndexSQL:           indexSql,
		PublishChanges:     publish,
		iraClusterInstance: ic,
	}
}

func (tdb *TDB) Open() string {
	arena := &Arena{}
	defer arena.free()

	cResponse := C.C_TdbOpen(
		arena.CString(tdb.ClusterID),
		arena.CString(tdb.DBName),
		arena.CString(tdb.Privacy.level.String()),
		arena.CString(tdb.InitSQL),
		arena.CString(tdb.IndexSQL),
		arena.CString(tdb.Privacy.publicKey),
		arena.CString(tdb.Privacy.privateKey),
	)
	arena.Add(unsafe.Pointer(cResponse))
	return C.GoString(cResponse)
}

func (tdb *TDB) Select(sql string) string {
	arena := &Arena{}
	defer arena.free()

	cResponse := C.C_TdbSelect(arena.CString(tdb.ClusterID), arena.CString(tdb.DBName), arena.CString(sql))
	arena.Add(unsafe.Pointer(cResponse))
	return C.GoString(cResponse)
}

func (tdb *TDB) Count(tableName, where string) int {
	arena := &Arena{}
	defer arena.free()

	return int(
		C.C_TdbCount(
			arena.CString(tdb.ClusterID),
			arena.CString(tdb.DBName),
			arena.CString(tableName),
			arena.CString(where),
		),
	)
}

func (tdb *TDB) Execute(sql string, publish ...bool) string {
	arena := &Arena{}
	defer arena.free()

	publishChanges := tdb.PublishChanges
	if len(publish) > 0 {
		publishChanges = publish[0]
	}

	cResponse := C.C_TdbExec(
		arena.CString(tdb.ClusterID),
		arena.CString(tdb.DBName),
		arena.CString(sql),
		C.bool(publishChanges),
	)
	arena.Add(unsafe.Pointer(cResponse))
	return C.GoString(cResponse)
}

func (tdb *TDB) ExecuteAsync(sql string, delay int, publish ...bool) {
	arena := &Arena{}
	defer arena.free()

	publishChanges := tdb.PublishChanges
	if len(publish) > 0 {
		publishChanges = publish[0]
	}
	C.C_TdbExecAsync(
		tdb.iraClusterInstance.cPtr,
		C.CString(tdb.ClusterID),
		C.CString(tdb.DBName),
		C.CString(sql),
		C.bool(publishChanges),
		C.int(delay),
	)
}

func (tdb *TDB) Close() bool {
	arena := &Arena{}
	defer arena.free()

	return bool(C.C_TdbClose(arena.CString(tdb.ClusterID), arena.CString(tdb.DBName)))
}

type callback func(string)

var Callbacks = []callback{}

//export IraClusterCallback
func IraClusterCallback(message *C.char) {
	goMsg := C.GoString(message)
	C.free(unsafe.Pointer(message))
	for _, cb := range Callbacks {
		cb(goMsg)
	}
}
