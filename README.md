# Paramstore loader

Paramstore loader is a CLI utility to load data, usually secrets, into AWS Param store. This way, app configuration can be stored in version control, allowing all aspects of an application deployment reproduceable at any time.

Storing secrets in a git repository is a solved problem: you can use [SOPS](https://github.com/getsops/sops), or [redact](https://github.com/julian7/redact) to keep your secrets safe in a git repository.

## Usage

```
paramstore-loader -input <inputfile>
```

The application requires AWS authentication, with the appropriate AWS profile set.

## Input file format

The input file is in JSON format:

```json
{
  "key_id": "alias/key",
  "param_root": "/appsecrets",
  "secrets": {
    "username": "admin",
    "password": "$6$rccsyyvVl2SEjDQB$eJ2IW2klbSMWAOED9Q.16OP7sByg2JV8Hxv/1NP.8McoIZuoSJ/5TYh81iY7t6RwDWOreXsVKX5tlD.SK.xlZ/",
    "config": {"file": "files/config.json"}
  }
}
```

Notes:

- `key_id` can store an AWS KMS key ID or a KMS alias name.
- `key_id` is not required. Parameters will be stored as Strings and not as SecureStrings if not given.
- Secrets can be provided inline (as strings), or as an object with a `file` member to read the value from.
- The encrypted `password` is "password" if you're curious.

## Anything's missing?

Open a ticket, perhaps a pull request. We support [GitHub Flow](https://guides.github.com/introduction/flow/). You might want to [fork](https://guides.github.com/activities/forking/) this project first.
