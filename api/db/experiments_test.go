// +build !integration

package db

import (
	"api/db/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExperimentInteractor(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()

	t.Run("CannotCreateExperimentWithEmptyId", func (t *testing.T) {
		experiment := &models.Experiment{}
		err := testDb.CreateExperiment(experiment)

		assert.Equal(t, "attempt to create Experiment without ID", err.Error())
	})

	t.Run("CreateAndQueryExperimentWithoutTrial", func(t *testing.T) {
		testDb.CreateExperiment(&models.Experiment{ID: "experiment1"})
		experiment, err := testDb.GetExperimentById("experiment1")

		assert.Nil(t, err)
		assert.Equal(t, "experiment1", experiment.ID)
		assert.Equal(t, 0, len(experiment.Trials))
	})

	t.Run("CreateAndQueryExperimentWithTrials", func(t *testing.T) {
		createModelNamed("model1")
		testDb.CreateTrial(&models.Trial{ID: "trial1", ExperimentID: "experiment1", ModelID: 1})
		testDb.CreateTrial(&models.Trial{ID: "trial2", ExperimentID: "experiment1", ModelID: 1})
		experiment, _ := testDb.GetExperimentById("experiment1")

		assert.Equal(t, 2, len(experiment.Trials))
		assert.Equal(t, "trial1", experiment.Trials[0].ID)
		assert.Equal(t, "trial2", experiment.Trials[1].ID)
	})
}
