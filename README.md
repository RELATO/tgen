# TGEN

Template based text generator.

This program takes one or more _Go_ template files and a data file and generates text by executing a 
template for each record in the data file. 

Currently only supports output to _stdout_ and sending emails, so this can be uses as a mail merge program.
Please don't use to spam people.

## Templates

The system expects at least one _Go_ template.

`-tg` is the file glob for the templates, e.g. _templates/*.tmpl_

`-tn` is the name of the template to apply to the data file, e.g. _root.tmpl_

## Input methods

### Data file

The path of the data file is given by:

`-d` followed by the path

The top line is expected to be the column headers and each value is a key that can be used in the template to reference the column value, e.g.

```
ID,Name
123,John Smith
456,Jane Jones
```
template:

```
User {{.ID}} name is {{.Name}}
```

Currently only csv, tsv and custom separators are supported. A variety of formats
will be allowed later:

- proper comma separated values, using _Go_ csv parser
- tab separated valus
- key value pairs: each line `key=value` with empty line between records
- fixed column widths
- XML markup: <key1>value1</key1>

### ODBC

The plan is to use ODBC as an alternative data source. To do.

### Stdin

Input is read from Stdin if `-d` is omitted.

## Output methods

`-out` followed by the output methods

### Console

By default prints each executed template to _stdout_.

### HTTP

Not supported yet. Will support execution of HTTP commands.

### SMTP

The output will be sent to an SMTP relay if the output is specified as one of:

- `-out my.from.address@example.com`

- `-out "my.from.address@example.com The subject line"`

The `SMTP` environment variable has to be set. Port must be specified or use `:smtp` as alternative for port 25.

The data input must have a `_rcpt` value for each record which contains the recipient address. In fact it can contain multiple
addresses separated by a comma. There is currently no `_cc` or `_bcc` option.

Attachments are not supported.

Example:

data file `balances.csv`

```
_rcpt,Name,Balance
john@example.com,John Smith,939.88
jane@example.org,Jane Jones,1090.7
```

template `templates/email1.t`

```
Hi {{.Name}},
Your account balance is {{.Balance | printf "%.2f"}}
```

Send emails:

```bash
./tgen -tg templates/email*.t -tn email1.t -d balances.csv -out 'admin@example.com Your balance'
```

The email function will add a `From:`, `To:` and `Subject:` line at the start of the body.

### ODBC

In the future we will be able to output to a database and fill a table with records, maybe.






