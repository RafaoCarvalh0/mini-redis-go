package main

const RedisEmptyArray = "*0\r\n"

const (
	RDB_OPCODE_EOF           = 0xFF
	RDB_OPCODE_SELECTDB      = 0xFE
	RDB_OPCODE_EXPIRETIME    = 0xFD
	RDB_OPCODE_EXPIRETIME_MS = 0xFC
	RDB_OPCODE_RESIZEDB      = 0xFB
	RDB_OPCODE_AUX           = 0xFA
)

const (
	RDB_TYPE_STRING       = 0x00
	RDB_TYPE_LIST         = 0x01
	RDB_TYPE_SET          = 0x02
	RDB_TYPE_ZSET         = 0x03
	RDB_TYPE_HASH         = 0x04
	RDB_TYPE_ZSET2        = 0x05
	RDB_TYPE_MODULE       = 0x06
	RDB_TYPE_MODULE2      = 0x07
	RDB_TYPE_HASH_ZIPMAP  = 0x09
	RDB_TYPE_LIST_ZIPLIST = 0x0A
	RDB_TYPE_SET_INTSET   = 0x0B
	RDB_TYPE_ZSET_ZIPLIST = 0x0C
	RDB_TYPE_HASH_ZIPLIST = 0x0D
)

const (
	LENGTH_6BIT  = 0
	LENGTH_14BIT = 1
	LENGTH_32BIT = 2
)

const (
	LENGTH_6BIT_MASK  = 0x3F // 00111111
	LENGTH_14BIT_MASK = 0x3F // 00111111
	LENGTH_TYPE_MASK  = 0xC0 // 11000000
)

const (
	RDB_HEADER_SIZE        = 5
	RDB_VERSION_SIZE       = 4
	RDB_EXPIRETIME_SIZE    = 4
	RDB_EXPIRETIME_MS_SIZE = 8
)
