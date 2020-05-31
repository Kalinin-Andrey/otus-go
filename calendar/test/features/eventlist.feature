Feature: receive list of events

  # first step
  Scenario: create event on a next day
    When I create event on a next day
    Then I send request for a list of events on a day
    And I receive list of events with length equal to one

  # second step
  Scenario: create event on a second day
    When I create event on a second day
    Then I send request for a list of events on a week
    And I receive list of events with length equal to two

  # third step
  Scenario: create event on a second week
    When I create event on a second week
    Then I send request for a list of events on a month
    And I receive list of events with length equal to three

