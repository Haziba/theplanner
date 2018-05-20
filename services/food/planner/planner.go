package planner

import "github.com/haziba/theplanner/models/food"

type PlannerService interface {
	CreatePlanner(models.Planner) (models.Planner, error)
	GetPlanner(id string) (*models.Planner, error)
	GetAllPlanners() ([]models.Planner, error)
}
