# Description

PGP Encryption using Golang & PostgreSQL

# HOW TO

### Create GPG Key

1. Create GPG Key 

`$ gpg --full-generate-key`

2. Export Public Key

`$ gpg -a --export {email} > db.rsa`

3. Export Private Key

`$ gpg -a --export-secret-key {email} > db.rsa.pub`

### Save GPG Key

save db.rsa & db.rsa.pub to folder `config/rsa`

### Update Env

copy file `.env.sample` to `.env` and update to your PostgreSQL connection

### Run Application

`$ go build . && ./data-encryption`