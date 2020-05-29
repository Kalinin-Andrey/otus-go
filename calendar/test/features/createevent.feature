Feature: create event

  # first step
  Scenario: Create valid event
    When I create event with UserID=1, Title="Событие 01", Description="Описание события 01"
    Then I receive status is OK
    And I receive event with ID and UserID=1, Title="Событие 01", Description="Описание события 01"

  # second step
  Scenario: Create invalid event
    When I create event with UserID=1, Title=""
    Then I receive status is not OK


