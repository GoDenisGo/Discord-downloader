# Discord-downloader

Small script to download and extract the latest Discord archive.
Tested for Ubuntu 20.04, your mileage may vary.

Discord is a real-time messenger client that hosts shared spaces
for hundreds of users to interact. It is a form of social media,
somewhat like a WhatsApp group but with a few additional features.

## Installation:

Clone the repository:
`git clone github.com/GoDenisGo/Discord-downloader`

Now open a terminal and enter the directory where the repo was cloned.
You will need this to build the package:
`go build Discord-downloader.go`

Run the executable:
`sudo ./Discord-downloader`

## How does Discord-downloader work?

1. The script visits the Discord website to obtain a copy of the gzipped 
Discord tar.
2. The script will automatically unzip and untar the software, saving
raw file contents to the `/opt/Discord/` directory.
3. In the future, the script will run additional functions to make the
installed files executable. Currently, this script works for me as I
have done some additional steps prior to writing this script. At this
stage, if you have never installed the Discord archive for yourself,
then this script will not be enough for you to run Discord immediately.
This is my next priority and it will be fixed when I find the time.