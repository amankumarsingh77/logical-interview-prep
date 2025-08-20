package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Log struct {
	Timestamp int
	User      string
	Service   string
}

func findSuspiciousUsers(logs []string) []string {
	processedLogs := getLogsByUser(logs)
	suspiciousSet := make(map[string]bool)
	for key, user := range processedLogs {
		j, i := 0, 0
		serviceUsed := make(map[string]int)
		for j < len(user) {
			serviceUsed[user[j].Service]++
			for user[j].Timestamp-user[i].Timestamp > 3600 {
				serviceUsed[user[i].Service]--
				if serviceUsed[user[i].Service] == 0 {
					delete(serviceUsed, user[i].Service)
				}
				i++
			}
			if len(serviceUsed) > 5 {
				suspiciousSet[key] = true
				break
			}
			j++
		}
	}
	suspicious := []string{}
	for sus := range suspiciousSet {
		suspicious = append(suspicious, sus)
	}
	sort.Strings(suspicious)
	return suspicious
}

func getLogsByUser(logs []string) map[string][]*Log {
	processedLogs := make(map[string][]*Log)
	for _, log := range logs {
		data := strings.Split(log, ",")
		if len(data) < 3 {
			fmt.Println("invalid log : Skipping")
			continue
		}
		timeStamp, _ := strconv.Atoi(data[0])
		userId := data[1]
		service := data[2]
		processedLogs[userId] = append(processedLogs[userId], &Log{
			Timestamp: timeStamp,
			User:      userId,
			Service:   service,
		})
	}
	for _, entries := range processedLogs {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp < entries[j].Timestamp
		})
	}
	return processedLogs
}
