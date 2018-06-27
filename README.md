<p align="center">
  <img src="logo.png" width="300px">
 </p>
Fastly enumerate through a pre-compiled list of AWS s3 buckets using DNS instead of HTTP with a list of DNS resolvers and multi-threading. Warning: Be aware that this is really shitty golang code. I wrote it without any prior knowledge of Go Lang but it seems to do the job. Feel free to contribute to make the tool better!

# Usage
```
go get github.com/miekg/dns
go get github.com/remeh/sizedwaitgroup
go build
./mass3 -w ./lists/buckets.txt -r ./lists/resolvers.txt -t 100
```
# Arguments

| argument | explanation | 
| --- | --- |
| -w | The wordlist with all the precompiled S3 buckets (bucketname.s3.amazonaws.com) |
| -r | List with all the resolvers |
| -t | The amount of threads to use, 10 is default |

# Questions & Answers
__Q: Why not generate all the "potential" s3 bucket names in the tool?__

A: This tool doesn't know the recon you've already collected, for example, subdomains. When you have a huge list of subdomains you can run alt-dns over it and try to find other S3 buckets that might not have a DNS record configured (yet).

__Q: The tool returns weird non-existing buckets__

A: The tool relies on the `lists/resolvers.txt` file to be acurate without any "bad" resolvers. You can use [fresh.sh](https://github.com/almroot/fresh.sh) to clean up the list of resolvers.

__Q: How many threads should I use?__

A: Depends on your resources, I personally use 500 threads which seems to work fine for me.

# Credits
Credits to [@koenrh](https://github.com/koenrh) who created [s3enum](https://github.com/koenrh/s3enum). I used some parts of his code and the way he identifies if a S3 bucket exists using DNS.
