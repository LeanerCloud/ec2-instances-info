package ec2instancesinfo

import "testing"

func TestData(t *testing.T) {
	tests := []struct {
		name     string
		instance EC2Instance
		wantErr  bool
	}{
		{
			name: "Parsing t2.nano memory, price, and ebs surcharge",
			instance: EC2Instance{
				InstanceType: "t2.nano",
				Memory:       0.5,
				VCPU:         1,
				Pricing: map[string]RegionPrices{
					"us-east-1": {
						Linux: Pricing{
							OnDemand: 0.0058,
						},
						MSWin: Pricing{
							OnDemand: 0.0081,
						},
						EBSSurcharge: 0.0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Parsing m3.2xlarge memory, price, and ebs surcharge",
			instance: EC2Instance{
				InstanceType: "m3.2xlarge",
				Memory:       30.0,
				VCPU:         8,
				Pricing: map[string]RegionPrices{
					"us-east-1": {
						Linux: Pricing{
							OnDemand: 0.532,
						},
						MSWin: Pricing{
							OnDemand: 1.036,
						},
						EBSSurcharge: 0.0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Parsing p2.16xlarge memory, price, GPUs and EBS surcharge",
			instance: EC2Instance{
				InstanceType: "p2.16xlarge",
				Memory:       732.0,
				VCPU:         64,
				GPU:          16,
				Pricing: map[string]RegionPrices{
					"us-east-1": {
						Linux: Pricing{
							OnDemand: 14.4,
						},
						MSWin: Pricing{
							OnDemand: 17.344,
						},
						EBSSurcharge: 0,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Data()
			if (err != nil) != tt.wantErr {
				t.Errorf("Data() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, i := range *got {
				if i.InstanceType != tt.instance.InstanceType {
					continue
				}

				if i.Memory != tt.instance.Memory {
					t.Errorf("Data(): %v, want memory %v, got %v",
						tt.instance.InstanceType, tt.instance.Memory, i.Memory)
				}

				if i.VCPU != tt.instance.VCPU {
					t.Errorf("Data(): %v, want CPUs %v, got %v",
						tt.instance.InstanceType, tt.instance.VCPU, i.VCPU)
				}

				if i.Pricing["us-east-1"].Linux.OnDemand != tt.instance.Pricing["us-east-1"].Linux.OnDemand {
					t.Errorf("Data(): %v, want price %v, got %v",
						tt.instance.InstanceType,
						tt.instance.Pricing["us-east-1"].Linux.OnDemand,
						i.Pricing["us-east-1"].Linux.OnDemand)
				}

				if i.Pricing["us-east-1"].MSWin.OnDemand != tt.instance.Pricing["us-east-1"].MSWin.OnDemand {
					t.Errorf("Data(): %v, want MSWin price %v, got %v",
						tt.instance.InstanceType,
						tt.instance.Pricing["us-east-1"].MSWin.OnDemand,
						i.Pricing["us-east-1"].MSWin.OnDemand)
				}

				if i.Pricing["us-east-1"].EBSSurcharge != tt.instance.Pricing["us-east-1"].EBSSurcharge {
					t.Errorf("Data(): %v, want ebs cost %v, got %v",
						tt.instance.InstanceType,
						tt.instance.Pricing["us-east-1"].EBSSurcharge,
						i.Pricing["us-east-1"].EBSSurcharge)
				}
			}
		})
	}
}
