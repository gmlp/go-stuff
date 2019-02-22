package creational

import "testing"

//Acceptance criteria

// When no counter has been created before, a new one is created with the value 0
// if a Counter has already been created, return this instance that holds the actual count
// if we call the method AddOne, the count must be incremented by 1
func TestGetInstance(t *testing.T) {
	counter1 := GetInstance()

	if counter1 == nil {
		//Test of acceptanxe  criteria 1 failed
		t.Error("expected pointer to Singleton after calling GetInstance(), not nil")
	}

	expectedCounter := counter1
	currentCount := counter1.AddOne()

	if currentCount != 1 {
		t.Errorf("After calling jfor the first time to count, the count must be 1 but is is %d\n", currentCount)
	}

	counter2 := GetInstance()
	if counter2 != expectedCounter {
		// Test 2 failed
		t.Errorf("Expected same instance in counter2 but it got a different instance/n")
	}

	currentCount = counter2.AddOne()
	if currentCount != 2 {
		// Test 3 failded
		t.Errorf("After calling 'AddOne' using the second counter, the currentCount must be 2 but was %d\n", currentCount)
	}

}
