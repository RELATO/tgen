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

The top line is expected to be the column headers and each value

Currently only simple comma separated values are supported. Not even CSV format but a variety of formats
will be allowed later:

- comma separated values
- tab separated valus
- key value pairs: each line `key=value` with empty line between records
- XML markup: <key1>value1</key1>

### ODBC

The plan is to use ODBC as an alternative data source. To do.

### Stdin

Will also be supported to pipe data into the program.

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

### ODBC

In the future we will be able to output to a database and fill a table with records





