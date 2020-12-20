// Code generated with autogen.gen DO NOT EDIT.
package autogen

var autogenDef = map[string]string{
	"agenda_item.assign":                       "agenda.can_manage",
	"agenda_item.create":                       "agenda.can_manage",
	"agenda_item.delete":                       "agenda.can_manage",
	"agenda_item.numbering":                    "agenda.can_manage",
	"agenda_item.sort":                         "agenda.can_manage",
	"agenda_item.update":                       "agenda.can_manage",
	"assignment.create":                        "assignment.can_manage",
	"assignment.delete":                        "assignment.can_manage",
	"assignment.update":                        "assignment.can_manage",
	"group.create":                             "user.can_manage",
	"group.delete":                             "user.can_manage",
	"group.set_permission":                     "user.can_manage",
	"group.update":                             "user.can_manage",
	"list_of_speakers.delete":                  "TODO",
	"list_of_speakers.delete_all_speakers":     "agenda.can_manage_list_of_speakers",
	"list_of_speakers.re_add_last":             "agenda.can_manage_list_of_speakers",
	"list_of_speakers.update":                  "agenda.can_manage_list_of_speakers",
	"mediafile.create_directory":               "mediafile.can_manage",
	"mediafile.delete":                         "mediafile.can_manage",
	"mediafile.move":                           "mediafile.can_manage",
	"mediafile.update":                         "mediafile.can_manage",
	"mediafile.upload":                         "mediafile.can_manage",
	"meeting.delete_all_speakers_of_all_lists": "agenda.can_manage_list_of_speakers",
	"meeting.set_font":                         "meeting.can_manage_logos_and_fonts",
	"meeting.set_logo":                         "meeting.can_manage_logos_and_fonts",
	"meeting.unset_font":                       "meeting.can_manage_logos_and_fonts",
	"meeting.unset_logo":                       "meeting.can_manage_logos_and_fonts",
	"motion.follow_recommendation":             "motions.can_manage_metadata",
	"motion.reset_recommendation":              "motions.can_manage_metadata",
	"motion.reset_state":                       "motion.can_manage_metadata",
	"motion.set_recommendation":                "motion.can_manage_metadata",
	"motion.sort":                              "motion.can_manage_metadata",
	"motion.update_metadata":                   "motion.can_manage_metadata",
	"motion_block.create":                      "motion.can_manage",
	"motion_block.delete":                      "motion.can_manage",
	"motion_block.update":                      "motion.can_manage",
	"motion_category.create":                   "motion.can_mange",
	"motion_category.delete":                   "motion.can_mange",
	"motion_category.number_motions":           "motion.can_manage",
	"motion_category.sort":                     "motion.can_manage",
	"motion_category.sort_motions_in_category": "motion.can_manage",
	"motion_category.update":                   "motion.can_mange",
	"motion_change_recommendation.create":      "motion.can_manage",
	"motion_change_recommendation.delete":      "motion.can_manage",
	"motion_change_recommendation.update":      "motion.can_manage",
	"motion_comment_section.create":            "motion.can_manage",
	"motion_comment_section.delete":            "motion.can_manage",
	"motion_comment_section.sort":              "motion.can_manage",
	"motion_comment_section.update":            "motion.can_manage",
	"motion_state.create":                      "motion.can_manage",
	"motion_state.delete":                      "motion.can_manage",
	"motion_state.update":                      "motion.can_manage",
	"motion_statute_paragraph.create":          "motion.can_manage",
	"motion_statute_paragraph.delete":          "motion.can_manage",
	"motion_statute_paragraph.sort":            "motion.can_manage",
	"motion_statute_paragraph.update":          "motion.can_manage",
	"motion_submitter.delete":                  "motion.can_manage",
	"motion_submitter.sort":                    "motion.can_manage",
	"motion_workflow.create":                   "motion.can_manage",
	"motion_workflow.delete":                   "motion.can_manage",
	"motion_workflow.update":                   "motion.can_manage",
	"speaker.end_speech":                       "agenda.can_manage_list_of_speakers",
	"speaker.sort":                             "agenda.can_manage_list_of_speakers",
	"speaker.speak":                            "agenda.can_manage_list_of_speakers",
	"speaker.update":                           "agenda.can_manage_list_of_speakers",
	"tag.create":                               "tag.can_manage",
	"tag.delete":                               "tag.can_manage",
	"tag.update":                               "tag.can_manage",
	"topic.create":                             "agenda.can_manage",
	"topic.delete":                             "agenda.can_manage",
	"topic.update":                             "agenda.can_manage",
	"user.create_temporary":                    "user.can_manage",
	"user.delete_temporary":                    "user.can_manage",
	"user.generate_new_password_temporary":     "user.can_manage",
	"user.reset_password_to_default_temporary": "user.can_manage",
	"user.set_password_temporary":              "user.can_manage",
	"user.update_temporary":                    "user.can_manage",
}