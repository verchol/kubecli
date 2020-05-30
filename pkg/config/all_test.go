package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SanitySuite struct {
	suite.Suite
	suite.SetupTestSuite
	suite.TearDownAllSuite
}
type SingleTestSuite struct {
	suite.Suite
	suite.SetupTestSuite
	suite.TearDownAllSuite
}

func (s *SingleTestSuite) SetupTest() {
	SetupTestContext("testContext1", "testkubecli")
}

func (s *SingleTestSuite) TearDownSuite() {
	DeleteTestContext("testContext1", "testkubecli")
}
func (s *SingleTestSuite) TestCase() {
	TestCreateRole(s.T())
}
func (s *SanitySuite) SetupTest() {
	SetupTestContext("testContext1", "testkubecli")
}

func (s *SanitySuite) TearDownSuite() {
	DeleteTestContext("testContext1", "testkubecli")
}

func (s *SanitySuite) TestCase_TestLoadWithRules() {
	TestLoadWithRules(s.T())

}
func (s *SanitySuite) TestCase_TestLocaCache() {
	TestLocalCache(s.T())
}
func (s *SanitySuite) TestCase_TestCluster() {
	TestCluster(s.T())
}
func (s *SanitySuite) TestCase_TestRoleOpts() {
	TestRoleOpts(s.T())
}

func (s *SanitySuite) TestCase_TestCreateServiceAccount() {
	TestCreateServiceAccount(s.T())
}
func (s *SanitySuite) TestCase_TestConnection() {
	TestConnection(s.T())
}
func (s *SanitySuite) TestCase_TestCreateContext() {
	TestCreateContext(s.T())
}
func (s *SanitySuite) TestCase_TestClusterRoleCreate() {
	TestClusterRoleCreate(s.T())
}
func (s *SanitySuite) TestCase_TestCreateRole() {
	TestCreateRole(s.T())
}

func (s *SanitySuite) TestCase_TestCreateRoleBinding() {
	TestCreateRoleBinding(s.T())
}
func (s *SanitySuite) TestCase_TestCreateAdminContext() {
	TestCreateAdminContext(s.T())
}
func TestSuite(t *testing.T) {
	suite.Run(t, new(SanitySuite))
}
func TestOnlySuite(t *testing.T) {
	suite.Run(t, new(SingleTestSuite))
}
func TestOnly(t *testing.T) {

	SetupTestContext("testContext1", "testkubecli")
	TestCreateRoleBinding(t)
	defer DeleteTestContext("testContext1", "testkubecli")
}

func TestDelete(t *testing.T) {
	SetupTestContext("testContext1", "testkubecli")

	defer DeleteTestContext("testContext1", "testkubecli")
	err := DeleteTestContext("testContext1", "kubeclitesting")
	if err != nil {
		t.Fatal(err)
	}
}
