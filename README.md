# encrypt-dir

A very simple cli helper that encrypts files in directories using AES-GCM (128bit). Particularly suited for en-/decrypting
sensitive information in git repositories.

## Installation

```
$ go get -u github.com/ory/encrypt-dir
$ encrypt-dir encrypt --key=JephkRhbfqzDHMAYUtHa6qcys4R4D48w some-directory-1 some-other-directory
```

## Usage

### Encrypt

```
$ ls some-directory-1
foo.txt
$ ls some-other-directory
bar.txt

$ encrypt-dir encrypt --key=<some-key> some-directory-1 some-other-directory

$ ls some-directory-1
foo.txt foo.txt.secure
$ ls some-other-directory
bar.txt foo.txt.secure
```

### Decrypt

```
$ ls some-directory-1
foo.txt.secure
$ ls some-other-directory
bar.txt.secure

$ encrypt-dir decrypt --key=<some-key> some-directory-1 some-other-directory

$ ls some-directory-1
foo.txt foo.txt.secure
$ ls some-other-directory
bar.txt foo.txt.secure
```
