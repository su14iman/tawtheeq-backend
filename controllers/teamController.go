// TeamController
package controllers

import (
	"fmt"
	"strconv"
	"time"

	"tawtheeq-backend/config"
	"tawtheeq-backend/models"
	"tawtheeq-backend/repositories"
	"tawtheeq-backend/utils"

	"github.com/gofiber/fiber/v2"
)

// CreateTeam godoc
// @Summary Create team
// @Description Create a new team
// @Tags teams
// @Accept json
// @Produce json
// @Param input body models.CreateTeamInput true "Create Team Input"
// @Success 201 {object} models.Team
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /teams [post]
// @Security Bearer
func CreateTeam(c *fiber.Ctx) error {
	repo := repositories.NewTeamRepository(config.DB)

	input := new(models.CreateTeamInput)
	if err := c.BodyParser(input); err != nil {
		utils.HandleError(err, "Failed to parse request body", utils.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "created_at": time.Now()})
	}

	team := &models.Team{
		Name:     input.Name,
		LeaderID: input.LeaderID,
	}

	// change Leader to TeamLeaderRole
	leaderRepo := repositories.NewUserRepository(config.DB)
	leader, err := leaderRepo.FindByID(input.LeaderID)
	if err != nil {
		utils.HandleError(err, "Failed to find leader", utils.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Leader not found", "created_at": time.Now()})
	}
	if leader.Role != models.TeamLeaderRole {
		leader.Role = models.TeamLeaderRole
	}

	if err := leaderRepo.Update(leader); err != nil {
		utils.HandleError(err, "Failed to update leader role", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update leader role", "created_at": time.Now()})
	}

	if err := repo.Create(team); err != nil {
		utils.HandleError(err, "Failed to create team", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create team", "created_at": time.Now()})
	}

	return c.Status(fiber.StatusCreated).JSON(team)
}

// RemoveTeam godoc
// @Summary Remove team
// @Description Remove a team by ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} models.Team
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /teams/{id} [delete]
// @Security Bearer
func RemoveTeam(c *fiber.Ctx) error {
	repo := repositories.NewTeamRepository(config.DB)
	id := c.Params("id")

	if err := repo.Delete(id); err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to delete team %s", id), utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete team"})
	}

	return c.JSON(fiber.Map{"message": "Team deleted", "team_id": id})
}

// UpdateTeamName godoc
// @Summary Update team name
// @Description Update the name of a team by ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Param input body models.ChangeTeamNameInput true "New team name"
// @Success 200 {object} models.Team
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /teams/{id}/name [put]
// @Security Bearer
func UpdateTeamName(c *fiber.Ctx) error {
	repo := repositories.NewTeamRepository(config.DB)
	id := c.Params("id")

	team, err := repo.FindByID(id)
	if err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to find team %s", id), utils.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Team not found"})
	}

	team.Name = c.FormValue("name")
	if err := repo.Update(team); err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to update team %s", id), utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update team name"})
	}

	return c.JSON(fiber.Map{"message": "Team name updated"})
}

// UpdateTeamLeader godoc
// @Summary Update team leader
// @Description Update the leader of a team by ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Param leader_id formData string true "New leader ID"
// @Success 200 {object} models.Team
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /teams/{id}/leader [put]
// @Security Bearer
func UpdateTeamLeader(c *fiber.Ctx) error {
	repo := repositories.NewTeamRepository(config.DB)
	id := c.Params("id")

	team, err := repo.FindByID(id)
	if err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to find team %s", id), utils.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Team not found", "created_at": time.Now()})
	}

	// update old leader role to TeamMemberRole
	oldLeaderRepo := repositories.NewUserRepository(config.DB)
	oldLeader, err := oldLeaderRepo.FindByID(team.LeaderID)
	if err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to find old leader %s", team.LeaderID), utils.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Old leader not found", "created_at": time.Now()})
	}
	if oldLeader.Role != models.TeamLeaderRole {
		utils.HandleError(err, fmt.Sprintf("Old leader %s is not a team leader", team.LeaderID), utils.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Old leader is not a team leader", "created_at": time.Now()})
	}
	oldLeader.Role = models.TeamMemberRole
	if err := oldLeaderRepo.Update(oldLeader); err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to update old leader %s role", team.LeaderID), utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update old leader role", "created_at": time.Now()})
	}

	leaderID := c.FormValue("leader_id")
	leaderRepo := repositories.NewUserRepository(config.DB)
	leader, err := leaderRepo.FindByID(leaderID)
	if err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to find new leader %s", leaderID), utils.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Leader not found", "created_at": time.Now()})
	}
	if leader.Role != models.TeamLeaderRole {
		utils.HandleError(err, fmt.Sprintf("New leader %s is not a team leader", leaderID), utils.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User is not a team leader", "created_at": time.Now()})
	}

	team.LeaderID = leaderID
	if err := repo.Update(team); err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to update team %s leader", id), utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update team leader", "created_at": time.Now()})
	}

	// update new leader role to TeamLeaderRole
	leader.Role = models.TeamLeaderRole
	if err := leaderRepo.Update(leader); err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to update new leader %s role", leaderID), utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update new leader role", "created_at": time.Now()})
	}

	return c.JSON(fiber.Map{"message": "Team leader updated"})
}

// GetMyTeam godoc
// @Summary Get my team
// @Description Get the team of the authenticated user
// @Tags teams
// @Accept json
// @Produce json
// @Success 200 {object} models.Team
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /my/team [get]
// @Security Bearer
func GetMyTeam(c *fiber.Ctx) error {

	teamIdRaw := c.Locals("teamId")
	teamId, ok1 := teamIdRaw.(string)
	if !ok1 {
		utils.HandleError(fmt.Errorf("invalid or missing user context"), "Invalid or missing user context", utils.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid or missing user context",
			"created_at": time.Now(),
		})
	}

	repo := repositories.NewTeamRepository(config.DB)

	team, err := repo.FindByID(teamId)
	if err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to find team %s", teamId), utils.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Team not found"})
	}

	return c.JSON(team)
}

// GetAllTeams godoc
// @Summary Get all teams
// @Description Get all teams with pagination
// @Tags teams
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param page query int false "Page" default(1)
// @Success 200 {array} models.Team
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /teams [get]
// @Security Bearer
func GetAllTeams(c *fiber.Ctx) error {
	repo := repositories.NewTeamRepository(config.DB)

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	total, err := repo.Count()
	if err != nil {
		utils.HandleError(err, "Failed to count teams", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count teams"})
	}

	teams, err := repo.FindAllPaginated(limit, offset)
	if err != nil {
		utils.HandleError(err, "Failed to get teams", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get teams"})
	}

	return c.JSON(fiber.Map{
		"teams": teams,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetAllUsersInTeam godoc
// @Summary Get all users in a team
// @Description Get all users in a team with pagination
// @Tags teams
// @Accept json
// @Produce json
// @Param team_id path string true "Team ID"
// @Param limit query int false "Limit" default(10)
// @Param page query int false "Page" default(1)
// @Success 200 {array} models.UserShortResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /teams/{team_id}/users [get]
// @Security Bearer
func GetAllUsersInTeam(c *fiber.Ctx) error {
	repo := repositories.NewTeamMemberRepository(config.DB)
	teamID := c.Params("team_id")

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	total, err := repo.CountMembers(teamID)
	if err != nil {
		utils.HandleError(err, "Failed to count team members", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count team members", "created_at": time.Now()})
	}

	members, err := repo.GetMembersPaginated(teamID, limit, offset)
	if err != nil {
		utils.HandleError(err, "Failed to get team members", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get team members", "created_at": time.Now()})
	}

	return c.JSON(fiber.Map{
		"members": members,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// AddUserToTeam godoc
// @Summary Add user to team
// @Description Add a user to a team
// @Tags teams
// @Accept json
// @Produce json
// @Param input body models.AddTeamMemberInput true "Add Team Member Input"
// @Success 200 {object} models.TeamMember
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /teams/members [post]
// @Security Bearer
func AddUserToTeam(c *fiber.Ctx) error {
	repo := repositories.NewTeamMemberRepository(config.DB)

	input := new(models.AddTeamMemberInput)
	if err := c.BodyParser(input); err != nil {
		utils.HandleError(err, "Failed to parse request body", utils.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	member := &models.TeamMember{
		TeamID: input.TeamID,
		UserID: input.UserID,
	}

	if err := repo.Add(member); err != nil {
		utils.HandleError(err, "Failed to add user to team", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add user to team"})
	}

	return c.JSON(fiber.Map{"message": "User added to team"})
}

// RemoveUserFromTeam godoc
// @Summary Remove user from team
// @Description Remove a user from a team
// @Tags teams
// @Accept json
// @Produce json
// @Param team_id path string true "Team ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} models.TeamMember
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /teams/members/{team_id}/{user_id} [delete]
// @Security Bearer
func RemoveUserFromTeam(c *fiber.Ctx) error {
	repo := repositories.NewTeamMemberRepository(config.DB)
	teamID := c.Params("team_id")
	userID := c.Params("user_id")

	if err := repo.Remove(teamID, userID); err != nil {
		utils.HandleError(err, "Failed to remove user from team", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove user from team", "created_at": time.Now()})
	}

	return c.JSON(fiber.Map{"message": "User removed from team"})
}

// GetAllUsersInMyTeam godoc
// @Summary Get all users in my team
// @Description Get all users in the authenticated user's team
// @Tags teams
// @Accept json
// @Produce json
// @Success 200 {array} models.UserShortResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /my/team/members [get]
// @Security Bearer
func GetAllUsersInMyTeam(c *fiber.Ctx) error {
	teamIdRaw := c.Locals("teamId")
	teamId, ok1 := teamIdRaw.(string)
	if !ok1 {
		utils.HandleError(fmt.Errorf("invalid or missing user context"), "Invalid or missing user context", utils.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid or missing user context",
			"created_at": time.Now(),
		})
	}

	repo := repositories.NewTeamMemberRepository(config.DB)

	members, err := repo.GetMembers(teamId)
	if err != nil {
		utils.HandleError(err, "Failed to get team members", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get team members", "created_at": time.Now()})
	}

	return c.JSON(members)
}

// AddUserToMyTeam godoc
// @Summary Add user to my team
// @Description Add a user to the authenticated user's team
// @Tags teams
// @Accept json
// @Produce json
// @Param input body models.AddTeamMemberToMyInput true "Add Team Member Input"
// @Success 200 {object} models.TeamMember
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /my/team/members [post]
// @Security Bearer
func AddUserToMyTeam(c *fiber.Ctx) error {
	teamIdRaw := c.Locals("teamId")
	teamId, ok1 := teamIdRaw.(string)
	if !ok1 {
		utils.HandleError(fmt.Errorf("invalid or missing user context"), "Invalid or missing user context", utils.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid or missing user context",
			"created_at": time.Now(),
		})
	}

	repo := repositories.NewTeamMemberRepository(config.DB)

	input := new(models.AddTeamMemberToMyInput)
	if err := c.BodyParser(input); err != nil {
		utils.HandleError(err, "Failed to parse request body", utils.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	member := &models.TeamMember{
		TeamID: teamId,
		UserID: input.UserID,
	}

	if err := repo.Add(member); err != nil {
		utils.HandleError(err, "Failed to add user to team", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add user to team"})
	}

	return c.JSON(fiber.Map{"message": "User added to team"})
}

// RemoveUserFromMyTeam godoc
// @Summary Remove user from my team
// @Description Remove a user from the authenticated user's team
// @Tags teams
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.TeamMember
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /my/team/members/{user_id} [delete]
// @Security Bearer
func RemoveUserFromMyTeam(c *fiber.Ctx) error {
	teamIdRaw := c.Locals("teamId")
	teamId, ok1 := teamIdRaw.(string)
	if !ok1 {
		utils.HandleError(fmt.Errorf("invalid or missing user context"), "Invalid or missing user context", utils.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid or missing user context",
			"created_at": time.Now(),
		})
	}

	repo := repositories.NewTeamMemberRepository(config.DB)
	userID := c.Params("user_id")

	if err := repo.Remove(teamId, userID); err != nil {
		utils.HandleError(err, "Failed to remove user from team", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove user from team", "created_at": time.Now()})
	}

	return c.JSON(fiber.Map{"message": "User removed from team", "created_at": time.Now()})
}
