# 📁 Scenario: Concurrent File Deduplicator

## Background  
Imagine you're tasked with building a **command-line utility** to help users clean up their storage.  
A key feature of this tool is finding all **duplicate files** within a directory, no matter what they're named.  

👉 Two files are considered duplicates **only if their contents are identical**.

---

## The Task  
Your goal is to write a Go function:  

```go
FindDuplicates(rootDir string)
```  

This function should recursively scan a given directory and all of its subdirectories.  
It should identify all files that have duplicate content and return a data structure that groups their paths together.

---

## Input  
- A single string representing the **path to the root directory** you need to scan.

## Output  
- The function should return a:

```go
map[string][]string
```  

Where:  
- **Key** → a unique identifier for the file's content (e.g., a hash).  
- **Value** → a slice of strings, where each string is the full path to a file that has that content.  

⚠️ The final map should only include entries where the slice of file paths has a **length ≥ 2** (i.e., actual duplicates).

---

## Key Requirements & Discussion Points  

### 1. File System Traversal  
- Walk the entire directory tree starting from `rootDir`.  

### 2. Content-Based Comparison  
- Files must be verified by **content**, not just size or name.  
- Use a hashing algorithm like **SHA-256** to compute file fingerprints.  

### 3. Efficiency  
- Reading & hashing large files (e.g., videos) is expensive.  
- Consider optimizations:  
  - ✅ Check file **size** before hashing (quick elimination).  
  - ✅ Use partial hashing (e.g., hash first N KB) as a filter before full hashing.  

### 4. Concurrency  
- This task is **highly parallelizable**.  
- After building a single-threaded solution, extend it to use **multiple CPU cores**.  
- Use goroutines and worker pools to process files concurrently.  

---

🚀 **Goal:** A robust, concurrent Go-based utility that efficiently finds and groups duplicate files by content.
