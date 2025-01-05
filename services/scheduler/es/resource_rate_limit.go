package es

var ResourceRateLimit = map[string]int{
	"Microsoft.Management/groups":                 5,
	"Microsoft.CostManagement/CostByResourceType": 3,
	"Microsoft.Storage/tables":                    5,
	"AWS::Organizations::Account":                 1,
	"AWS::Organizations::Root":                    1,
	"AWS::Organizations::Organization":            1,
	"AWS::Organizations::OrganizationalUnit":      1,
	"AWS::Organizations::PolicyTarget":            1,
	"AWS::Organizations::Policy":                  1,
	"AWS::Shield::ProtectionGroup":                5,
	"AWS::IAM::Policy":                            5,
	"AWS::IAM::Role":                              5,
	"AWS::SES::ConfigurationSet":                  5,
	"AWS::IAM::CredentialReport":                  5,
}
