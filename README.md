## kstore

KStore is a simple console password manager storing its data encrypted on Yandex.Disk

For security reasons OAuth credentials are not included in the source code, so to build it you have to create a new App at oauth.yandex.ru, give it permissions to read and write its own application folder, acquire your `client_id` and `client_secret` from there and use the following command to build `kstore`:

```bash
CLIENT_ID=<Your client id>
CLIENT_SECRET=<Your client secret>

go build -ldflags "-X 'github.com/viert/kstore/manager.clientID=$CLIENT_ID' -X 'github.com/viert/kstore/manager.clientSecret=$CLIENT_SECRET'" cmd/kstore/main.go
```

At start kstore will ask you to type in your master password. Then it will try to load and decrypt (with AES key based on your password) yandex credentials file. If the program is running for the first time and there's no credentials file, kstore will generate a link to Yandex OAuth following which you give it permissions it needs and get a confirmation code. Once confirmation code is typed into kstore, the program will create an encrypted credentials file for you and use it the next time you start kstore.

Important note: all the data starting from credentials file to the actual passwords database is encrypted with your master password. In case of losing it you will not be able to restore your data. The credentials file is important but it's quite safe to delete it. It will be generated upon next start using the previously mentioned OAuth "handshake".
