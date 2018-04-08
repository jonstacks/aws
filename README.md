# aws

AWS golang pkg, binaries, utils, etc.

<!-- TOC depthFrom:2 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Install](#install)
- [Reservation Audits](#reservation-audits)
	- [reserved-instance-audit](#reserved-instance-audit)
	- [reserved-rds-audit](#reserved-rds-audit)

<!-- /TOC -->

## Install

You can install most simply by using `go get`:

```
go get -u github.com/jonstacks/aws/cmd/...
```

## Reservation Audits

Occasionally, you might find yourself wanting to quickly audit reservations
that you've purchased, making sure you are using all of your reserved instances,
or maybe finding out how many more you need to reserve. You can use the
following commands found in the `cmd` folder for just that purpose.

### reserved-instance-audit

Audits your reserved EC2 instances against your currently running ones.

![reserved-instance-audit](doc/screenshots/reserved-instance-audit.png)

Options:

* `--only-unmatched`: Show only instance types that are not perfectly reserved.

### reserved-rds-audit

Audits your reserved RDS instances against your currently running ones.

![reserved-rds-audit](doc/screenshots/reserved-rds-audit.png)

## Getting a download URL for your RDS logs

It appears that both the `awscli` and sdk libraries are broken for downloading
RDS logs. I have experienced this grief as well. So, I have made a program
based on the git issues below that will give you a download link to download
the complete logs:

* https://github.com/aws/aws-cli/issues/2268
* https://github.com/aws/aws-cli/issues/3079

You can install it with:

```
go get -u https://github.com/jonstacks/aws/cmd/rds-logs-download-url
```

Along with a program like `curl` you can have a complete solution to
downloading your RDS logs like so:

```sh
#!/bin/bash

set -e

# This program requires that you supply the AWS_DEFAULT_REGION as an env var.
# Can possibly pull this out of an AWS config in the future.
export AWS_DEFAULT_REGION="us-west-2"

LOG_NAME="error/postgresql.log.2018-04-08-15"

LOG_URL=$(rds-logs-download-url fc16fu3t5aah9e9 $LOG_NAME)
# -s for silent, -f for fail so we can retry on failure
curl -sf -o $(basename $LOG_NAME) $LOG_URL
```

*Note*: I have only really tested this against the `us-west-2` region. It might
        need some additional changes to support others use cases.