# go-local-smtp-app
A local SMTP server app

## Setup

Create an empty directory named `attachments` in the main app directory:

```
mkdir attachments
```

Rename the `authfile.sample` file to `authfile`:

```
cp authfile.sample authfile
```

Fill the `authfile` with the username and MD5-encoded password to
login via an SMTP session:

```
account@local:81dc9bdb52d04dc20036dbd8313ed055
```

The `:` character is the separator between username and password.
