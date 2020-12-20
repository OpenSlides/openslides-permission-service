// Code generated with autogen.gen DO NOT EDIT.
package collection

var derivatePerms = map[string][]string{
	"agenda.can_be_speaker":              {"agenda.can_see"},
	"agenda.can_manage":                  {"agenda.can_manage_list_of_speakers", "agenda.can_see_list_of_speakers", "agenda.can_see", "agenda.can_be_speaker", "agenda.can_see", "agenda.can_see_internal_items", "agenda.can_see"},
	"agenda.can_manage_list_of_speakers": {"agenda.can_see_list_of_speakers", "agenda.can_see"},
	"agenda.can_see":                     {},
	"agenda.can_see_internal_items":      {"agenda.can_see"},
	"agenda.can_see_list_of_speakers":    {"agenda.can_see"},
	"assignment.can_manage":              {"assignment.can_nominate_other", "assignment.can_see", "assignment.can_nominate_self", "assignment.can_see"},
	"assignment.can_nominate_other":      {"assignment.can_see"},
	"assignment.can_nominate_self":       {"assignment.can_see"},
	"assignment.can_see":                 {},
	"mediafile.can_manage":               {"mediafile.can_see"},
	"mediafile.can_see":                  {},
	"meeting.can_manage":                 {"meeting.can_manage_logos_and_fonts", "meeting.can_manage_projector", "meeting.can_see_projector", "meeting.can_see_history", "meeting.can_see_autopilot", "meeting.can_see_frontpage"},
	"meeting.can_manage_logos_and_fonts": {},
	"meeting.can_manage_projector":       {"meeting.can_see_projector"},
	"meeting.can_see_autopilot":          {"meeting.can_see_frontpage"},
	"meeting.can_see_frontpage":          {},
	"meeting.can_see_history":            {},
	"meeting.can_see_projector":          {},
	"motion.can_create":                  {"motion.can_see"},
	"motion.can_create_amendments":       {"motion.can_see"},
	"motion.can_manage":                  {"motion.can_manage_metadata", "motion.can_support", "motion.can_see", "motion.can_see_internal", "motion.can_see", "motion.can_create", "motion.can_see", "motion.can_create_amendments", "motion.can_see"},
	"motion.can_manage_metadata":         {"motion.can_support", "motion.can_see"},
	"motion.can_see":                     {},
	"motion.can_see_internal":            {"motion.can_see"},
	"motion.can_support":                 {"motion.can_see"},
	"tag.can_manage":                     {},
	"user.can_change_password":           {},
	"user.can_manage":                    {"user.can_see_extra_data", "user.can_see"},
	"user.can_see":                       {},
	"user.can_see_extra_data":            {"user.can_see"},
}
