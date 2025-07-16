package models

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PlantsTestSuite defines our test suite
type PlantsTestSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

// SetupSuite runs once before all tests
func (suite *PlantsTestSuite) SetupSuite() {
	mockDb, mock, err := sqlmock.New()
	require.NoError(suite.T(), err)

	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), mock)
	require.NotNil(suite.T(), db)
	suite.db = db
	suite.mock = mock
}

// TearDownTest runs after each test to clean up
func (suite *PlantsTestSuite) TearDownTest() {
	// Clean up the table after each test
	// suite.db.Exec("DELETE FROM plants")

	// TODO - make this do something
}

// Helper function to create valid datetime string
func validDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (suite *PlantsTestSuite) CreatePlant(successExpected bool, plants []*Plants) []*gorm.DB {

	suite.mock.ExpectBegin()
	results := make([]*gorm.DB, len(plants))
	for index, plant := range plants {
		suite.mock.ExpectQuery(`INSERT INTO "plants"`).WithArgs(plant.Name, plant.OutdoorSowDate, plant.DaysToGermination, plant.DaysToHarvest,
			plant.CanStartIndoors, plant.IndoorSowDate, plant.TransplantDate).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(index + 1))
		suite.mock.ExpectCommit()
		result := suite.db.Create(&plant)

		if successExpected {
			assert.NoError(suite.T(), result.Error)
			assert.Equal(suite.T(), int64(1), result.RowsAffected)
			assert.Equal(suite.T(), index+1, plant.ID)

			err := suite.mock.ExpectationsWereMet()
			assert.NoError(suite.T(), err)

		} else {
			assert.Error(suite.T(), result.Error)
			assert.Equal(suite.T(), int64(0), result.RowsAffected)
		}

		results[index] = result
	}
	return results
}

// Test creating a valid plant entry
func (suite *PlantsTestSuite) TestCreateValidPlant() {
	plant := Plants{
		Name:              "Tomato",
		OutdoorSowDate:    "05-01",
		DaysToGermination: 7,
		DaysToHarvest:     85,
		CanStartIndoors:   true,
		IndoorSowDate:     "04-01",
		TransplantDate:    "05-15",
	}
	plants := []*Plants{&plant}
	suite.CreatePlant(true, plants)
}

// Test creating plant without optional fields
func (suite *PlantsTestSuite) TestCreatePlantWithoutOptionalFields() {
	plant := Plants{
		Name:              "Lettuce",
		OutdoorSowDate:    validDateTime(),
		DaysToGermination: 5,
		DaysToHarvest:     60,
		CanStartIndoors:   false,
		// IndoorSowDate and TransplantDate are nil (optional)
	}

	plants := []*Plants{&plant}
	suite.CreatePlant(true, plants)
	assert.Equal(suite.T(), "", plant.IndoorSowDate)
	assert.Equal(suite.T(), "", plant.TransplantDate)

}

func (suite *PlantsTestSuite) TestUniqueNameConstraint() {
	plant1 := Plants{
		Name:              "Carrot",
		OutdoorSowDate:    validDateTime(),
		DaysToGermination: 10,
		DaysToHarvest:     70,
		CanStartIndoors:   false,
	}

	// First plant should succeed
	suite.CreatePlant(true, plant1)

	// Second plant with same name should fail
	plant2 := Plants{
		Name:              "Carrot",
		OutdoorSowDate:    validDateTime(),
		DaysToGermination: 12,
		DaysToHarvest:     75,
		CanStartIndoors:   false,
	}
	fail := suite.CreatePlant(false, plant2)
	assert.NotNil(suite.T(), fail.Error)
	// TODO - check fail error contains some description that indicates unique constraint violation
}

// Test name length constraint (max 50 characters)
func (suite *PlantsTestSuite) TestNameSizeConstraint() {
	longName := string(make([]byte, 51)) // 51 characters
	for i := range longName {
		longName = longName[:i] + "a" + longName[i+1:]
	}

	plant := Plants{
		Name:              longName,
		OutdoorSowDate:    validDateTime(),
		DaysToGermination: 7,
		DaysToHarvest:     85,
		CanStartIndoors:   true,
	}

	result := suite.db.Create(&plant)
	assert.Error(suite.T(), result.Error)
}

// Test required fields validation
func (suite *PlantsTestSuite) TestRequiredFieldsValidation() {
	testCases := []struct {
		name  string
		plant Plants
	}{
		{
			name: "missing name",
			plant: Plants{
				// Name is empty
				OutdoorSowDate:    validDateTime(),
				DaysToGermination: 7,
				DaysToHarvest:     85,
				CanStartIndoors:   true,
			},
		},
		{
			name: "missing outdoor sow date",
			plant: Plants{
				Name: "Test Plant",
				// OutdoorSowDate is empty
				DaysToGermination: 7,
				DaysToHarvest:     85,
				CanStartIndoors:   true,
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			result := suite.db.Create(&tc.plant)
			assert.Error(t, result.Error)
		})
	}
}
