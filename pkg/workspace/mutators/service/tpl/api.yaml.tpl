openapi: "3.0.0"
info:
  version: 1.0.0
  title: {{.ServiceName}}
  description: Service description
  contact:
    name: Maintainer name
    email: Maintainer email
servers:
  - url: {{.ServiceName}}.example.com
paths: {}
# Example of a handler, uncomment and remove the above 'paths: {}' line.
# Check Petstore OpenAPI example for more possible options:
# https://github.com/OAI/OpenAPI-Specification/blob/main/examples/v3.0/petstore-expanded.yaml
#
# paths:
#   /path/to/api:
#     get:
#       summary: sample handler
#       responses:
#         '200':
#           description: OK
#           content:
#             application/json:
#               schema:
#                 $ref: '#/components/schemas/PathToApiResponse'
# components:
#   schemas:
#     PathToApiResponse:
#       type: object
#       properties:
#         value:
#           type: string
#       required:
#         - value
