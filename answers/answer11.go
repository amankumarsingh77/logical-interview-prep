package main

//
//import (
//	"fmt"
//	"sort"
//	"strconv"
//	"strings"
//	"time"
//)
//
//type Log struct {
//	Timestamp int
//	User      string
//	Service   string
//}
//
//func findSuspiciousUsers(logs []string) []string {
//	processedLogs := getLogsByUser(logs)
//	suspiciousSet := make(map[string]bool)
//	for key, user := range processedLogs {
//		j, i := 0, 0
//		serviceUsed := make(map[string]int)
//		for j < len(user) {
//			serviceUsed[user[j].Service]++
//			for user[j].Timestamp-user[i].Timestamp > 3600 {
//				serviceUsed[user[i].Service]--
//				if serviceUsed[user[i].Service] == 0 {
//					delete(serviceUsed, user[i].Service)
//				}
//				i++
//			}
//			if len(serviceUsed) > 5 {
//				suspiciousSet[key] = true
//				break
//			}
//			j++
//		}
//	}
//	var suspicious []string
//	for sus := range suspiciousSet {
//		suspicious = append(suspicious, sus)
//	}
//	return suspicious
//}
//
//func getLogsByUser(logs []string) map[string][]*Log {
//	processedLogs := make(map[string][]*Log)
//	for _, log := range logs {
//		data := strings.Split(log, ",")
//		if len(data) < 3 {
//			fmt.Println("invalid log : Skipping")
//		}
//		timeStamp, _ := strconv.Atoi(data[0])
//		userId := data[1]
//		service := data[2]
//		processedLogs[userId] = append(processedLogs[userId], &Log{
//			Timestamp: timeStamp,
//			User:      userId,
//			Service:   service,
//		})
//	}
//	for _, entries := range processedLogs {
//		sort.Slice(entries, func(i, j int) bool {
//			return entries[i].Timestamp < entries[j].Timestamp
//		})
//	}
//	return processedLogs
//}
//
//func main() {
//	fmt.Println("Running Test Case 1: Basic Scenario")
//	logs1 := []string{
//		"1672531200,user-a,service-auth",    // 10:00:00
//		"1672531260,user-b,service-auth",    // 10:01:00
//		"1672531320,user-a,service-storage", // 10:02:00
//		"1672531380,user-a,service-compute", // 10:03:00
//		"1672531440,user-a,service-db",      // 10:04:00
//		"1672534740,user-a,service-network", // 10:59:00 (5th service)
//		"1672534800,user-a,service-logging", // 11:00:00 (6th service in 60 mins) -> user-a is suspicious
//		"1672534860,user-b,service-cache",   // 11:01:00
//	}
//	expected1 := []string{"user-a"}
//	runTest(logs1, expected1)
//
//	fmt.Println("\nRunning Test Case 2: Multiple Users and Noisy Logs")
//	logs2 := []string{
//		"1672621200,user-c,service-A", // Day 2 10:00
//		"1672531200,user-a,service-A", // Day 1 10:00
//		"1672531260,user-b,service-B", // Day 1 10:01
//		"1672531320,user-a,service-C",
//		"1672621500,user-c,service-B",
//		"1672534740,user-a,service-D",
//		"1672621800,user-c,service-C",
//		"1672534800,user-a,service-E", // user-a has 4 services in 60 mins, not suspicious
//		"1672622100,user-c,service-D",
//		"1672622400,user-c,service-E",
//		"1672622700,user-c,service-F", // user-c, 6th service -> suspicious
//		"1672538400,user-a,service-F", // Day 1 11:00 - this is exactly 60 mins after first log, so window includes it. 5th service.
//		"1672538401,user-a,service-G", // Day 1 11:00:01 - this makes user-a suspicious
//	}
//	expected2 := []string{"user-a", "user-c"}
//	runTest(logs2, expected2)
//
//	fmt.Println("\nRunning Test Case 3: No Suspicious Activity")
//	logs3 := []string{
//		"1672531200,user-d,service-A",
//		"1672534801,user-d,service-B", // access is just outside the 60min window (3601 seconds)
//		"1672538402,user-d,service-C",
//	}
//	expected3 := []string{}
//	runTest(logs3, expected3)
//
//	fmt.Println("\nRunning Test Case 4: Empty Logs")
//	logs4 := []string{}
//	expected4 := []string{}
//	runTest(logs4, expected4)
//}
//
//func runTest(logs, expected []string) {
//	result := findSuspiciousUsers(logs)
//	sort.Strings(result) // Ensure result is sorted for comparison
//
//	fmt.Printf("Input Logs:\n")
//	for _, l := range logs {
//		parts := strings.Split(l, ",")
//		ts, _ := strconv.ParseInt(parts[0], 10, 64)
//		fmt.Printf("  - %s (%s, %s)\n", time.Unix(ts, 0).Format("15:04:05"), parts[1], parts[2])
//	}
//
//	fmt.Printf("Expected: %v\n", expected)
//	fmt.Printf("Got:      %v\n", result)
//
//	if equalSlices(result, expected) {
//		fmt.Println("Result: ✅ PASS")
//	} else {
//		fmt.Println("Result: ❌ FAIL")
//	}
//}
//
//func equalSlices(a, b []string) bool {
//	if len(a) != len(b) {
//		return false
//	}
//	for i := range a {
//		if a[i] != b[i] {
//			return false
//		}
//	}
//	return true
//}
