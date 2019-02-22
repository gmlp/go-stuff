package creational

import (
	"strings"
	"testing"
)

func TestCreatePaymentMethodCash(t *testing.T) {
	payment, err := GetPaymentMethod(Cash)
	if err != nil {
		t.Fatal("A payment method of type 'Cash' must exist")
	}

	msg := payment.Pay(10.30)
	if !strings.Contains(msg, "paid using cash") {
		t.Error("The cash payment method message wasn't correct")
	}
	t.Log("Log:", msg)
}

func TestCreatePaymentMethodDebitCard(t *testing.T) {
	payment, err := GetPaymentMethod(DebitCard)
	if err != nil {
		t.Fatal("A payment method of type 'Debit Card' must exist")
	}

	msg := payment.Pay(10.30)
	if !strings.Contains(msg, "paid using debit card") {
		t.Error("The cash payment method message wasn't correct")
	}
	t.Log("Log:", msg)
}

func TestGetPaymentMethodNonExistent(t *testing.T) {
	_, err := GetPaymentMethod(20)

	if err == nil {
		t.Error("A payment method with ID 20 must return error")
	}
	t.Log("Log:", err)
}
