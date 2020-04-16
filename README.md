# Mix networks - CS 588 Project

[Project documentation](https://docs.google.com/document/d/1DW5OnHH5xCbAnnPUICl4LapEKzsGJL_2hQ4sDW-LiB8/edit)

## Generating key pairs

```bash
openssl genrsa out privkey.pem 4096
openssl rsa -in privkey.pem -out pubkey.pem -outform PEM
```
