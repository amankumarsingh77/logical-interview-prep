# Scenario: The Resumable File Downloader

Imagine you are tasked with building a command-line utility to download large files, potentially several gigabytes in size. The network connection is unreliable and might drop at any time.

The core requirement is this: If the download is interrupted, running the utility again with the same arguments must resume the download from where it left off, without re-downloading the parts you already have.

## Your Task

Design and implement a function or class that takes two arguments:

- `url`: The URL for the file to download
- `local_filepath`: The path where the file should be saved locally

## Key Questions to Consider

1. **Determining Partial Downloads:**
    - How can you determine if a partial file already exists locally and where you should resume from?

2. **Requesting File Portions:**
    - How do you instruct a web server to send you only a specific portion of a file, rather than the whole thing from the beginning?

3. **Detecting File Changes:**
    - What happens if the file on the server has changed since your last download attempt?
    - How could you detect that?

4. **Handling Servers Without Partial Download Support:**
    - How should you handle cases where the server doesn't support partial downloads?

