<p align="center">
  <img src="logo.png" width="300px">
 </p>
Quickly enumerate through a pre-compiled list of AWS S3 buckets using DNS instead of HTTP with a list of DNS resolvers and multi-threading. Warning: Be aware that this is really shitty golang code. I wrote it without any prior knowledge of Go Lang but it seems to do the job. Feel free to contribute to make the tool better!

# Install
```shell
go get -u github.com/smiegles/mass3
```

# Usage
```shell
mass3 -w ./lists/buckets.txt -r ./lists/resolvers.txt -t 100
```

# Arguments
| argument | explanation | 
| --- | --- |
| -w | The wordlist with all the pre-compiled S3 buckets (bucketname.s3.amazonaws.com) |
| -r | List with all the resolvers |
| -t | The amount of threads to use, 10 is default |
| -o | The file path to save the output (This is optional. By default, it will be saved to out.csv) |

# Building the Docker Image
`docker build -t <name> .`

# Running the Docker Image
`docker run -it <name> -w buckets.txt -r resolvers.txt -t 100 -o /tmp/out.csv`

# Questions & Answers
__Q: Why not generate all the "potential" s3 bucket names in the tool?__

A: This tool doesn't know the recon you've already collected, for example, subdomains. When you have a huge list of subdomains you can run alt-dns over it and try to find other S3 buckets that might not have a DNS record configured (yet).

__Q: The tool returns weird non-existing buckets__

A: The tool relies on the `lists/resolvers.txt` file to be accurate without any "bad" resolvers. You can use [fresh.sh](https://github.com/almroot/fresh.sh) to clean up the list of resolvers.

__Q: How many threads should I use?__

A: Depends on your resources, I personally use 500 threads which seems to work fine for me.

# Credits
Credits to [@koenrh](https://github.com/koenrh) who created [s3enum](https://github.com/koenrh/s3enum). I used some parts of his code and the way he identifies if a S3 bucket exists using DNS.
