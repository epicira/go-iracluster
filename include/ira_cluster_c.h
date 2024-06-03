#pragma once
#include <stdbool.h>
#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

char* C_TdbOpen(const char* iraClusterId, const char* dbName, const char* privacyLevel, const char* initSql, 
    const char* indexes, const char* publicKey, const char* privateKey);
char* C_TdbSelect(const char* iraClusterId, const char* dbName, const char* sql);
int C_TdbCount(const char* iraClusterId, const char* dbName, const char* tableName, const char* where);
char* C_TdbExec(const char* iraClusterId, const char* dbName, const char* sql, bool publish);
void C_TdbExecAsync(void* iraCM, const char* iraClusterId, const char* dbName, const char* sql, bool publish, int ms);
bool C_TdbClose(const char* iraClusterId, const char* dbName);

void* NewIraCluster(const char* modName, const char* clusterId, const char* version);
void DestroyIraCluster(void* iraCM);

bool joinCluster(void* iraCM);
bool isSenior(void* iraCM);
int peerCount(void* iraCM);

#ifdef __cplusplus
}
#endif