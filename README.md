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

### reserved-rds-audit

Audits your reserved RDS instances against your currently running ones.

![reserved-rds-audit](doc/screenshots/reserved-rds-audit.png)
