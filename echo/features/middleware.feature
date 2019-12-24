Feature: reject requests without correct headers
  As a security conscious application developer
  I need to reject requests not signed by Clubhouse

  Scenario: Reject requests without header
    Given request does not have clubhouse header
    When request is made
    Then request is rejected with status code 400