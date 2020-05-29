package config

import (
	"fmt"
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
	SetupTestContext("testContext_1", "testkubecli")
}

func (s *SingleTestSuite) TearDownSuite() {
	DeleteTestContext("testContext_1", "testkubecli")
}
func (s *SingleTestSuite) TestCase() {
	TestCreateRole(s.T())
}
func (s *SanitySuite) SetupTest() {
	SetupTestContext("testContext_1", "testkubecli")
}

func (s *SanitySuite) TearDownSuite() {
	DeleteTestContext("testContext_1", "testkubecli")
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
func (s *SanitySuite) TestCase_TestDeleteRole() {
	TestDeleteRole(s.T())
}
func (s *SanitySuite) TestCase_TestCreateRoleBinding() {
	TestCreateRoleBinding(s.T())
}
func TestSuite(t *testing.T) {
	suite.Run(t, new(SanitySuite))
}
func TestOnlySuite(t *testing.T) {
	suite.Run(t, new(SingleTestSuite))
}
func TestOnly(t *testing.T) {

	SetupTestContext("testContext_1", "testkubecli")
	TestCreateRoleBinding(t)
	defer DeleteTestContext("testContext_1", "testkubecli")
}
func TestAll(t *testing.T) {

	SetupTestContext("testContext_1", "testkubecli")

	defer DeleteTestContext("testContext_1", "testkubecli")
	fmt.Println("TestLoadWithRules")
	TestLoadWithRules(t)
	fmt.Println("TestLocalCache")
	TestLocalCache(t)
	fmt.Println("TestCluster")
	TestCluster(t)
	fmt.Println("TestRoleOpts")
	TestRoleOpts(t)
	fmt.Println("TestCreateServiceAccount")
	TestCreateServiceAccount(t)
	fmt.Println("TestConnection")
	TestConnection(t)
	fmt.Println("TestCreateContext")
	TestCreateContext(t)
	fmt.Println("TestClusterRoleCreate")
	TestClusterRoleCreate(t)
	fmt.Println("TestCreateRole")
	TestCreateRole(t)
	fmt.Println("TestDeleteRole")
	TestDeleteRole(t)
	fmt.Println("TestCreateRoleBinding")
	TestCreateRoleBinding(t)

	fmt.Println("all tests were executed")

}

func TestDelete(t *testing.T) {
	SetupTestContext("testContext_1", "testkubecli")

	defer DeleteTestContext("testContext_1", "testkubecli")
	err := DeleteTestContext("testContext_1", "kubeclitesting")
	if err != nil {
		t.Fatal(err)
	}
}
