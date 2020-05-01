package mq

import "Distributed-fileserver/common"

type TransferData struct{
	FileHash string
	CurLocation string
	DestLocation string
	DestStoretype common.StoreType
}