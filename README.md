# Mix networks - CS 588 Project

[Project documentation](https://docs.google.com/document/d/1DW5OnHH5xCbAnnPUICl4LapEKzsGJL_2hQ4sDW-LiB8/edit)

## Generating key pairs

```bash
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -outform PEM -pubout -out public.pem
```

## Running tests

```bash
make test
```
