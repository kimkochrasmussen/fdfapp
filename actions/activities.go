package actions

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/mogensen/fdfapp/models"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Activity)
// DB Table: Plural (activities)
// Resource: Plural (Activities)
// Path: Plural (/activities)
// View Template Folder: Plural (/templates/activities/)

// ActivitiesResource is the resource for the Activity model
type ActivitiesResource struct {
	buffalo.Resource
}

// List gets all Activities. This function is mapped to the path
// GET /activities
func (v ActivitiesResource) List(c buffalo.Context) error {
	c.Set("TIME_FORMAT", "02 Jan 2006")
	activities := &models.Activities{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := scope(c).PaginateFromParams(c.Params())

	if c.Param("class_id") != "" {
		q = q.Where("class_id = (?)", c.Param("class_id"))

		class := &models.Class{}

		// Retrieve all participants from the DB
		if err := scope(c).Eager().Find(class, c.Param("class_id")); err != nil {
			return err
		}
		c.Set("class", class)
	}

	// Retrieve all Activities from the DB
	if err := q.Eager().All(activities); err != nil {
		return err
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, activities))
}

// Show gets the data for one Activity. This function is mapped to
// the path GET /activities/{activity_id}
func (v ActivitiesResource) Show(c buffalo.Context) error {
	c.Set("TIME_FORMAT", "02 Jan 2006")
	// Allocate an empty Activity
	activity := &models.Activity{}

	// To find the Activity the parameter activity_id is used.
	if err := scope(c).Eager().Find(activity, c.Param("activity_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, activity))
}

// New renders the form for creating a new Activity.
// This function is mapped to the path GET /activities/new/
func (v ActivitiesResource) New(c buffalo.Context) error {
	var err error

	// Se if a date has already been supplied by a calender event
	date := time.Now()
	if c.Param("date") != "" {
		date, err = time.Parse("2006-01-02 15:04:05 -0700 MST", c.Param("date"))
		if err != nil {
			return err
		}

	}

	// Se if a duration has already been supplied by a calender event
	duration := 0.0
	if c.Param("duration") != "" {
		duration, _ = strconv.ParseFloat(c.Param("duration"), 64)
	}

	class := &models.Class{}

	// Retrieve all participants from the DB
	if err := scope(c).Eager("Participants.Image").Find(class, c.Param("class_id")); err != nil {
		return err
	}

	c.Set("participants", class.Participants)
	c.Set("class", class)
	return c.Render(200, r.Auto(c, &models.Activity{
		Date:     date,
		ClassID:  class.ID,
		Title:    c.Param("title"),
		Duration: duration,
	}))
}

// Create adds a Activity to the DB. This function is mapped to the
// path POST /activities
func (v ActivitiesResource) Create(c buffalo.Context) error {
	// Allocate an empty Activity
	activity := &models.Activity{}

	// Bind activity to the html form elements
	if err := c.Bind(activity); err != nil {
		return err
	}
	activity.UserID = currentUser(c).ID

	participants := models.Participants{}
	for _, participantID := range c.Request().Form["participants"] {

		cID, err := uuid.FromString(participantID)
		if err != nil {
			return err
		}
		participants = append(participants, models.Participant{ID: cID})
		fmt.Println(cID)
		fmt.Println(participants)
	}
	activity.Participants = participants

	fmt.Println(activity)
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.Eager().ValidateAndCreate(activity)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, activity))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", T.Translate(c, "activity.created.success"))
	// and redirect to the activities index page
	return c.Render(201, r.Auto(c, activity))
}

// Edit renders a edit form for a Activity. This function is
// mapped to the path GET /activities/{activity_id}/edit
func (v ActivitiesResource) Edit(c buffalo.Context) error {

	// Allocate an empty Activity
	activity := &models.Activity{}

	if err := scope(c).Eager("Participants.Image").Find(activity, c.Param("activity_id")); err != nil {
		return c.Error(404, err)
	}

	class := &models.Class{}

	// Retrieve all participants from the DB
	if err := scope(c).Eager().Find(class, activity.ClassID); err != nil {
		return err
	}

	for i, v := range class.Participants {
		part := &models.Participant{}
		if err := scope(c).Eager().Find(part, v.ID); err != nil {
			return err
		}
		class.Participants[i] = *part
	}

	c.Set("participants", class.Participants)
	c.Set("class", class)

	return c.Render(200, r.Auto(c, activity))
}

// Update changes a Activity in the DB. This function is mapped to
// the path PUT /activities/{activity_id}
func (v ActivitiesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Allocate an empty Activity
	activity := &models.Activity{}

	if err := scope(c).Find(activity, c.Param("activity_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Activity to the html form elements
	if err := c.Bind(activity); err != nil {
		return err
	}
	activity.UserID = currentUser(c).ID

	if err := removeAllActivityParticipants(*activity, tx); err != nil {
		return fmt.Errorf("Failed removing activityParticipants: %v", err)
	}

	for _, participantID := range c.Request().Form["participants"] {

		cID, err := uuid.FromString(participantID)
		if err != nil {
			return err
		}
		ap := &models.ActivityParticipant{
			ActivityID:    activity.ID,
			ParticipantID: cID,
		}
		if err := tx.Create(ap); err != nil {
			return err
		}
	}

	verrs, err := tx.ValidateAndUpdate(activity)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, activity))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", T.Translate(c, "activity.updated.success"))
	// and redirect to the activities index page
	return c.Render(200, r.Auto(c, activity))
}

func removeAllActivityParticipants(activity models.Activity, tx *pop.Connection) error {

	activityParticipants := &models.ActivityParticipants{}

	err := tx.Where("activity_id = ?", activity.ID).All(activityParticipants)
	if err != nil {
		return err
	}
	log.Println(activityParticipants)

	if err := tx.Destroy(activityParticipants); err != nil {
		return fmt.Errorf("Failed destroying activityParticipant: %v", err)
	}
	return nil
}

// Destroy deletes a Activity from the DB. This function is mapped
// to the path DELETE /activities/{activity_id}
func (v ActivitiesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Allocate an empty Activity
	activity := &models.Activity{}

	// To find the Activity the parameter activity_id is used.
	if err := scope(c).Find(activity, c.Param("activity_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(activity); err != nil {
		return err
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", T.Translate(c, "activity.destroyed.success"))
	// Redirect to the activities index page
	return c.Render(200, r.Auto(c, activity))
}
