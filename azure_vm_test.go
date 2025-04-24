package ec2instancesinfo

import (
	"testing"
)

func TestAzureData(t *testing.T) {
	data, err := AzureData()
	if err != nil {
		t.Fatalf("Error getting Azure VM data: %v", err)
	}
	
	if data == nil {
		t.Fatal("Azure VM data is nil")
	}
	
	if len(*data) == 0 {
		t.Fatal("Azure VM data is empty")
	}
	
	// Test a few fields from the first instance
	firstInstance := (*data)[0]
	if firstInstance.InstanceType != "a0" {
		t.Errorf("Expected instance type 'a0', got '%s'", firstInstance.InstanceType)
	}
	
	if firstInstance.VCPU != 1 {
		t.Errorf("Expected 1 VCPU, got %d", firstInstance.VCPU)
	}
	
	if firstInstance.Memory != 0.75 {
		t.Errorf("Expected memory 0.75 GB, got %.2f", firstInstance.Memory)
	}
	
	// Test pricing data
	if pricing, found := firstInstance.Pricing["us-east"]; found {
		if pricing.Linux.OnDemand != 0.02 {
			t.Errorf("Expected Linux OnDemand price 0.02, got %.4f", pricing.Linux.OnDemand)
		}
		
		if pricing.Windows.OnDemand != 0.02 {
			t.Errorf("Expected Windows OnDemand price 0.02, got %.4f", pricing.Windows.OnDemand)
		}
	} else {
		t.Error("Expected pricing for 'us-east' region not found")
	}
}