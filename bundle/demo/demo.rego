package demo

import rego.v1

# METADATA
# description: Allow only until a certain amount
# entrypoint: true
# custom:
#   schema:
#     type: object
#     properties:
#       default_max_amount:
#         type: integer
#         note: Maximal Amount
#       mode:
#         type: string
#         note: The Mode
#         enum: ["info", "warn", "enforce"]
#       valid_from:
#         type: dateTime
#         note: Only Valid From
#       valid_to:
#         type: dateTime
#         note: Only Valid To
#       log_decisions:
#         type: boolean
#         note: To Log Decisions
default amount_allowed := false

amount_allowed if {
	input.amount <= data.parameters.default_max_amount
}
