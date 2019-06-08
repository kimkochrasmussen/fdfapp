package actions

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/tags"
	"github.com/gofrs/uuid"
	"github.com/mogensen/fdfapp/models"
)

var r *render.Engine
var assetsBox = packr.New("app:assets", "../public")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.New("app:templates", "../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			"checkboxChecked": func(id uuid.UUID, slice models.Participants) string {
				for _, c := range slice {
					if id == c.ID {
						return "checked"
					}
				}
				return ""
			},

			"isActive": func(name string, help plush.HelperContext) string {
				if cp, ok := help.Value("current_route").(buffalo.RouteInfo); ok {
					if strings.HasPrefix(cp.PathName, name) {
						return "active"
					}
				}
				return "inactive"
			},
			"getParticipant": func(cms models.ClassMembership, participants *models.Participants) models.Participant {
				for _, p := range *participants {
					if p.ID == cms.ParticipantID {
						return p
					}
				}
				return models.Participant{}
			},
			"getClass": func(cms models.ClassMembership, classes *models.Classes) models.Class {
				for _, c := range *classes {
					if c.ID == cms.ClassID {
						return c
					}
				}
				return models.Class{}
			},

			"buttonGroupButton": func(text, icon, link string, help plush.HelperContext) (template.HTML, error) {
				a := tags.New("a", tags.Options{"class": "btn btn-secondary btn-sm", "href": link})
				i := tags.New("i", tags.Options{"class": fmt.Sprintf("fas fa-%s", icon)})
				a.Append(i)
				a.Append(tags.New("br", tags.Options{}))
				a.Append(text)
				return a.HTML(), nil
			},
			"buttonGroup": func(help plush.HelperContext) (template.HTML, error) {
				group := tags.New("div", tags.Options{"class": "btn-group float-right", "role": "group"})
				if help.HasBlock() {
					bc, err := help.Block()
					if err != nil {
						return "", err
					}
					group.Append(bc)
				}
				return group.HTML(), nil
			},

			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,
		},
	})
}
