package assumerole

type (
	AssumeRole struct {
		AccessKeyID     string
		SecretAccessKey string
		Region          string
		RoleARN         string
		ExternalID      string
	}
)
