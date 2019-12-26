Feature: only accept requests with valid headers
  As an app dev using clubhouse webhooks
  I need to only accept properly signed requests
  So my application only handles valid clubhouse webhooks

  Scenario: Reject requests without header
    Given request does not have clubhouse header
    When request is made
    Then request is rejected with status code 400

  Scenario: Reject requests with invalid header
    Given request has a garbage clubhouse header
    When request is made
    Then request is rejected with status code 400

  Scenario: Reject requests with empty body
    Given request has a valid clubhouse header
    Given request has empty body
    When request is made
    Then request is rejected with status code 401

  Scenario: Reject requests with mismatched signature
    Given request signature does not match request
    When request is made
    Then request is rejected with status code 401

  Scenario: Allow requests with correct signature
    Given request signature does match request
    When request is made
    Then request is accepted with status code 204
