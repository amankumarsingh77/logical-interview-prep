package main

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {

}

type LogEntry struct {
	Timestamp int
	Service   string
}

type UserState struct {
	RecentLogs    *list.List
	UniqueService map[string]bool
	LastSeen      int64
}

var userActivity = make(map[string]*UserState)

const windowLimit = 3600
const serviceLimit = 5
const cleanThreshold = 7200

func processLog(log string) {
	userId, entry := formatLog(log)
	addLog(entry, userId)

}

func cleanClutteredLogs() {
	for userId, state := range userActivity {
		if time.Now().UnixNano()-state.LastSeen > cleanThreshold {
			delete(userActivity, userId)
		}
	}
}

func addLog(log LogEntry, userId string) {
	for e := userActivity[userId].RecentLogs.Front(); e != nil; e = e.Next() {
		if log.Timestamp-e.Value.(LogEntry).Timestamp > windowLimit {
			userActivity[userId].RecentLogs.Remove(e)
		}
	}
	userActivity[userId].RecentLogs.PushBack(log)
}

func formatLog(log string) (string, LogEntry) {
	data := strings.Split(log, ",")
	if len(data) < 3 {
		fmt.Println("invalid log : Skipping")
	}
	timeStamp, _ := strconv.Atoi(data[0])
	userId := data[1]
	service := data[2]
	return userId, LogEntry{
		Timestamp: timeStamp,
		Service:   service,
	}
}
