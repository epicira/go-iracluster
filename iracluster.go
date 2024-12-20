package iracluster

import (
	"time"
	"unsafe"
)

/*
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
	internalCallback = func(event string) {
		switch event {
		case "iracluster::received_pass":
			licenseReceived = true
		case "iracluster::cluster_joined":
			clusterJoined = true
		}
	}
	Callbacks = append(Callbacks, cb)
	return iraCluster
}

func (ic *IraCluster) JoinCluster() bool {
	ok := bool(C.joinCluster(ic.cPtr))
	if !ok {
		return false
	}
	println("Waiting to join cluster and receive licenses...")
	for !(clusterJoined && licenseReceived) {
		time.Sleep(1 * time.Second)
	}
	println("Joined cluster and received licenses")
	return true
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

func (p PrivacyLevel) String() string {
	return [...]string{"none", "local", "shared", "private"}[p]
}

func ParsePrivacyLevel(level string) (PrivacyLevel, bool) {
	switch level {
	case None.String():
		return None, true
	case Local.String():
		return Local, true
	case Shared.String():
		return Shared, true
	case Private.String():
		return Private, true
	default:
		return None, false
	}
}

type Privacy struct {
	Level      PrivacyLevel
	PublicKey  string
	PrivateKey string
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
			Level:      None,
			PublicKey:  "",
			PrivateKey: "",
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
		arena.CString(tdb.Privacy.Level.String()),
		arena.CString(tdb.InitSQL),
		arena.CString(tdb.IndexSQL),
		arena.CString(tdb.Privacy.PublicKey),
		arena.CString(tdb.Privacy.PrivateKey),
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
var internalCallback callback

var clusterJoined bool
var licenseReceived bool

//export IraClusterCallback
func IraClusterCallback(message *C.char) {
	goMsg := C.GoString(message)
	C.free(unsafe.Pointer(message))
	for _, cb := range Callbacks {
		cb(goMsg)
	}
}

//export InternalCallback
func InternalCallback(message *C.char) {
	goMsg := C.GoString(message)
	// I am not freeing the memory here as the responsibility is taken care by the caller.
	// But, this behavior is not consistent across different callbacks, so, keeping this as a reminder to make the behavior consistent.
	// C.free(unsafe.Pointer(message))
	if internalCallback != nil {
		internalCallback(goMsg)
	}
}
