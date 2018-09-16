package models

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
)

// DefaultSession creates and initializes the underlying default session for
// working with models.
func DefaultSession() *session.Session {
	return session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			},
		),
	)
}

// Init initializes the clients used by the models.
func Init(s *session.Session) {
	EC2Client(ec2.New(s))
	RDSClient(rds.New(s))
}
