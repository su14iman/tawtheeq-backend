# tawtheeq-backend
a simple backend server that can be used to sign files (images/pdfs) and verify the signature of files. The backend is written in golang and uses the [exiftool](https://exiftool.org/) to sign the files. The server uses RSA keys to sign and verify the files.

## how to run
first you have install golang on your system, you can follow the instructions [here](https://golang.org/doc/install/source).
make sure you have the requirments app installed on your system and you have the keys generated.
then you clone the repo and run the server using the following commands:
```bash
# run the server
go run main.go
```

## keys generation
to generate keys, you have to run the following commands:
### private key
```bash
openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
```

### public key
```bash
openssl rsa -pubout -in private.pem -out public.pem
```

## requirements
you have to install the following app on your system:
### 1. exiftool
```bash
# for MacOS
brew install exiftool
# for Ubuntu
sudo apt-get install libimage-exiftool-perl
```
